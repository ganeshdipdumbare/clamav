import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Function to determine the appropriate API URL
function getApiUrl(): string {
    // In the browser, we need to use relative URLs or the server's public hostname
    if (browser) {
        // For development environment (when using Vite dev server)
        const isDev = window.location.port === '5173';
        
        if (isDev) {
            // In development, we need to specify the API host and port
            return 'http://localhost:8080';
        }
        
        // For production Docker environment, use relative URLs
        // This works because the browser will use the same host it's loaded from
        return '';
    }
    
    // When running on the server (SSR), we need to use the service name from docker-compose
    return 'http://api:8080';
}

// Create the store with the appropriate initial value
export const apiUrl = writable(getApiUrl()); 