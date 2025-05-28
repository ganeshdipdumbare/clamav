![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/ganeshdipdumbare/clamav?utm_source=oss&utm_medium=github&utm_campaign=ganeshdipdumbare%2Fclamav&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)

# ClamAV Virus Scanner

A modern web interface for ClamAV virus scanning with a responsive UI built using SvelteKit and a Go backend.

## Features

- Server status checking (ping and version)
- Text scanning for virus signatures
- File upload and scanning
- Dark/light mode toggle
- Mobile-friendly responsive design
- Docker containerization for easy deployment

## Architecture

The application consists of:

1. **SvelteKit Frontend**: Modern UI with Tailwind CSS
2. **Go REST API**: Backend server that communicates with ClamAV
3. **ClamAV Daemon**: Virus scanning engine

## Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for development)
- Go 1.20+ (for development)

## Quick Start

The easiest way to run the application is using Docker Compose:

```bash
# Production mode
docker-compose up -d

# Development mode
docker-compose -f docker-compose.dev.yml up -d
```

The application will be available at:
- Production: http://localhost:3000
- Development: http://localhost:5173

## Development

### Frontend (SvelteKit)

```bash
cd clamav-sveltekit
npm install
npm run dev
```

### Backend (Go)

```bash
cd clamav-server
go get
go run main.go
```

## Environment Variables

### Frontend

- `API_PROXY_TARGET`: URL of the backend API (default: http://api:8080)
- `ORIGIN`: The origin URL (default: http://localhost:3000 in production)

### Backend

- `CLAMD_ADDRESS`: Address of the ClamAV daemon (default: clamd:3310)
- `PORT`: Port to run the API server (default: 8080)

## License

MIT

## Credits

- [ClamAV](https://www.clamav.net/) - Open-source antivirus engine
- [SvelteKit](https://kit.svelte.dev/) - Web application framework
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [Font Awesome](https://fontawesome.com/) - Icon library 