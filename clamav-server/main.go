package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/baruwa-enterprise/clamd"
)

// ScanResult represents the result of a scan
type ScanResult struct {
	Filename string `json:"filename"`
	Status   string `json:"status"`
	Clean    bool   `json:"clean"`
}

// Server represents the REST API server
type Server struct {
	clamdClient *clamd.Client
}

// CORS middleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NewServer creates a new Server instance
func NewServer(clamdAddress string) (*Server, error) {
	client, err := clamd.NewClient("tcp", clamdAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create ClamAV client: %v", err)
	}

	// Set some sensible timeouts
	client.SetConnTimeout(5 * time.Second)
	client.SetCmdTimeout(30 * time.Second)

	return &Server{
		clamdClient: client,
	}, nil
}

// PingHandler handles ping requests to ClamAV
func (s *Server) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pingResult, err := s.clamdClient.Ping(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to ping ClamAV: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": pingResult,
		"message": "ClamAV daemon is responding",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VersionHandler returns the ClamAV version
func (s *Server) VersionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	version, err := s.clamdClient.Version(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get ClamAV version: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"version": version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ScanTextHandler scans text submitted in the request body
func (s *Server) ScanTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	reader := strings.NewReader(string(body))
	scanResults, err := s.clamdClient.ScanReader(ctx, reader)
	if err != nil {
		http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
		return
	}

	results := []ScanResult{}
	for _, result := range scanResults {
		results = append(results, ScanResult{
			Filename: result.Filename,
			Status:   result.Status,
			Clean:    result.Status == "OK",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

// ScanFileHandler handles file upload and scanning
func (s *Server) ScanFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form, 10 MB max
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Received file: %s, size: %d bytes", handler.Filename, handler.Size)

	scanResults, err := s.clamdClient.ScanReader(ctx, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
		return
	}

	results := []ScanResult{}
	for _, result := range scanResults {
		results = append(results, ScanResult{
			Filename: handler.Filename,
			Status:   result.Status,
			Clean:    result.Status == "OK",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"filename": handler.Filename,
		"size":     handler.Size,
		"results":  results,
	})
}

// IndexHandler redirects to the static HTML page
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/static/", http.StatusFound)
}

func main() {
	// Get ClamAV address from environment or use default
	clamdAddress := os.Getenv("CLAMD_ADDRESS")
	if clamdAddress == "" {
		clamdAddress = "clamd:3310" // Using Docker service name
	}

	// Create server
	server, err := NewServer(clamdAddress)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create a new multiplexer with CORS middleware
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/", server.IndexHandler)
	mux.HandleFunc("/ping", server.PingHandler)
	mux.HandleFunc("/version", server.VersionHandler)
	mux.HandleFunc("/scan/text", server.ScanTextHandler)
	mux.HandleFunc("/scan/file", server.ScanFileHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Wrap the mux with CORS middleware
	handler := corsMiddleware(mux)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	serverAddr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting server on %s", serverAddr)
	log.Printf("ClamAV daemon address: %s", clamdAddress)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
