#!/bin/bash

# Docker Deployment Script for MCP Obsidian Server
# This script helps deploy the MCP Obsidian Server using Docker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Function to check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        print_warning ".env file not found. Creating template..."
        cat > .env << EOF
# Required: Your Obsidian Local REST API key
OBSIDIAN_API_KEY=your-api-key-here

# Optional: Obsidian API configuration
OBSIDIAN_HOST=127.0.0.1
OBSIDIAN_PORT=27124
OBSIDIAN_USE_HTTPS=true
OBSIDIAN_PROTOCOL=https
OBSIDIAN_VAULT_PATH=/path/to/your/vault
EOF
        print_warning "Please edit .env file with your actual API key and configuration before running the server."
        exit 1
    fi
    print_success ".env file found"
}

# Function to build images
build_images() {
    print_status "Building Docker images..."
    
    # Build HTTP image
    print_status "Building HTTP image..."
    docker build -f Dockerfile.http -t mcp-obsidian-http:latest .
    
    # Build SSE image
    print_status "Building SSE image..."
    docker build -f Dockerfile.sse -t mcp-obsidian-sse:latest .
    
    print_success "All images built successfully"
}

# Function to start services
start_services() {
    local service=$1
    
    case $service in
        "http")
            print_status "Starting HTTP service..."
            docker-compose up -d mcp-obsidian-http
            print_success "HTTP service started on port 8080"
            ;;
        "sse")
            print_status "Starting SSE service..."
            docker-compose up -d mcp-obsidian-sse
            print_success "SSE service started on port 8081"
            ;;
        "both")
            print_status "Starting both services..."
            docker-compose up -d
            print_success "Both services started (HTTP: 8080, SSE: 8081)"
            ;;
        *)
            print_error "Invalid service. Use 'http', 'sse', or 'both'"
            exit 1
            ;;
    esac
}

# Function to stop services
stop_services() {
    print_status "Stopping services..."
    docker-compose down
    print_success "Services stopped"
}

# Function to show logs
show_logs() {
    local service=$1
    
    if [ -z "$service" ]; then
        print_status "Showing logs for all services..."
        docker-compose logs -f
    else
        print_status "Showing logs for $service service..."
        docker-compose logs -f $service
    fi
}

# Function to show status
show_status() {
    print_status "Service status:"
    docker-compose ps
}

# Function to show help
show_help() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  build                    Build Docker images"
    echo "  start [http|sse|both]    Start services (default: both)"
    echo "  stop                     Stop all services"
    echo "  logs [service]           Show logs (default: all services)"
    echo "  status                   Show service status"
    echo "  restart [service]        Restart services"
    echo "  clean                    Clean up containers and images"
    echo "  help                     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build                 # Build all images"
    echo "  $0 start http            # Start HTTP service only"
    echo "  $0 start sse             # Start SSE service only"
    echo "  $0 start both            # Start both services"
    echo "  $0 logs mcp-obsidian-http # Show HTTP service logs"
    echo "  $0 stop                  # Stop all services"
}

# Function to restart services
restart_services() {
    local service=$1
    
    print_status "Restarting services..."
    if [ -z "$service" ]; then
        docker-compose restart
        print_success "All services restarted"
    else
        docker-compose restart $service
        print_success "$service service restarted"
    fi
}

# Function to clean up
clean_up() {
    print_warning "This will remove all containers and images. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_status "Cleaning up..."
        docker-compose down -v
        docker rmi mcp-obsidian-http:latest mcp-obsidian-sse:latest 2>/dev/null || true
        print_success "Cleanup completed"
    else
        print_status "Cleanup cancelled"
    fi
}

# Main script logic
main() {
    local command=$1
    local service=$2
    
    # Check if Docker is running
    check_docker
    
    case $command in
        "build")
            check_env_file
            build_images
            ;;
        "start")
            check_env_file
            start_services ${service:-both}
            ;;
        "stop")
            stop_services
            ;;
        "logs")
            show_logs $service
            ;;
        "status")
            show_status
            ;;
        "restart")
            restart_services $service
            ;;
        "clean")
            clean_up
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        "")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
