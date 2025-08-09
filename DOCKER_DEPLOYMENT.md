# Docker Deployment Guide for MCP Obsidian Server

This guide covers how to deploy the MCP Obsidian Server using Docker with support for both HTTP and SSE transports.

## Prerequisites

1. **Docker**: Install Docker and Docker Compose
2. **Obsidian Local REST API Plugin**: Install and enable the [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) plugin in your Obsidian vault
3. **API Key**: Generate an API key from the plugin settings

## Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd mcp-obsidian
```

### 2. Set Environment Variables
Create a `.env` file in the project root:
```bash
# Required: Your Obsidian Local REST API key
OBSIDIAN_API_KEY=your-api-key-here

# Optional: Obsidian API configuration
OBSIDIAN_HOST=127.0.0.1
OBSIDIAN_PORT=27124
OBSIDIAN_USE_HTTPS=true
OBSIDIAN_PROTOCOL=https
OBSIDIAN_VAULT_PATH=/path/to/your/vault
```

### 3. Deploy with Docker Compose

#### Deploy Both Services (HTTP + SSE)
```bash
# Build and start both services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Deploy Individual Services

**HTTP Transport Only:**
```bash
# Build and start HTTP service
docker-compose up -d mcp-obsidian-http

# View logs
docker-compose logs -f mcp-obsidian-http
```

**SSE Transport Only:**
```bash
# Build and start SSE service
docker-compose up -d mcp-obsidian-sse

# View logs
docker-compose logs -f mcp-obsidian-sse
```

## Manual Docker Build and Run

### HTTP Transport

#### Build the Image
```bash
# Build HTTP Docker image
docker build -f Dockerfile.http -t mcp-obsidian-http .

# Or with a specific tag
docker build -f Dockerfile.http -t mcp-obsidian-http:latest .
```

#### Run the Container
```bash
# Run HTTP container
docker run -d \
  --name mcp-obsidian-http \
  -p 8080:8080 \
  -e OBSIDIAN_API_KEY=your-api-key-here \
  -e OBSIDIAN_HOST=127.0.0.1 \
  -e OBSIDIAN_PORT=27124 \
  -e OBSIDIAN_USE_HTTPS=true \
  mcp-obsidian-http:latest
```

### SSE Transport

#### Build the Image
```bash
# Build SSE Docker image
docker build -f Dockerfile.sse -t mcp-obsidian-sse .

# Or with a specific tag
docker build -f Dockerfile.sse -t mcp-obsidian-sse:latest .
```

#### Run the Container
```bash
# Run SSE container
docker run -d \
  --name mcp-obsidian-sse \
  -p 8081:8081 \
  -e OBSIDIAN_API_KEY=your-api-key-here \
  -e OBSIDIAN_HOST=127.0.0.1 \
  -e OBSIDIAN_PORT=27124 \
  -e OBSIDIAN_USE_HTTPS=true \
  mcp-obsidian-sse:latest
```

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OBSIDIAN_API_KEY` | ✅ Yes | - | Your Obsidian Local REST API key |
| `OBSIDIAN_HOST` | ❌ No | `127.0.0.1` | Obsidian API host |
| `OBSIDIAN_PORT` | ❌ No | `27124` | Obsidian API port |
| `OBSIDIAN_USE_HTTPS` | ❌ No | `true` | Use HTTPS for API calls |
| `OBSIDIAN_PROTOCOL` | ❌ No | `https` | Protocol to use (http/https) |
| `OBSIDIAN_VAULT_PATH` | ❌ No | - | Path to your Obsidian vault |

### Port Configuration

- **HTTP Transport**: Port 8080 (default)
- **SSE Transport**: Port 8081 (default)

You can customize these ports by modifying the `docker-compose.yml` file or using the `-p` flag with `docker run`.

## Health Checks

Both Docker images include health checks that verify the service is running:

- **HTTP Service**: `http://localhost:8080/health`
- **SSE Service**: `http://localhost:8081/health`

## Monitoring and Logs

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f mcp-obsidian-http
docker-compose logs -f mcp-obsidian-sse

# Or with docker run
docker logs -f mcp-obsidian-http
docker logs -f mcp-obsidian-sse
```

### Check Service Status
```bash
# Check running containers
docker-compose ps

# Check container health
docker inspect mcp-obsidian-http | grep -A 10 "Health"
docker inspect mcp-obsidian-sse | grep -A 10 "Health"
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure Obsidian Local REST API plugin is running
   - Verify the API key is correct
   - Check if the host and port are accessible

2. **Permission Denied**
   - The containers run as non-root user (appuser)
   - Ensure proper file permissions if mounting volumes

3. **Port Already in Use**
   - Change the port mapping in `docker-compose.yml`
   - Or stop existing services using the same ports

### Debug Mode

To run in debug mode, you can override the CMD in docker-compose.yml:

```yaml
services:
  mcp-obsidian-http:
    # ... other configuration
    command: ["./mcp-obsidian-http", "--http-port", "8080", "--verbose"]
```

## Production Deployment

### Security Considerations

1. **Use Secrets Management**: Store sensitive data like API keys in Docker secrets or environment variables
2. **Network Security**: Use Docker networks to isolate services
3. **Resource Limits**: Set appropriate CPU and memory limits
4. **Logging**: Configure proper logging and monitoring

### Example Production docker-compose.yml

```yaml
version: '3.8'

services:
  mcp-obsidian-http:
    build:
      context: .
      dockerfile: Dockerfile.http
    container_name: mcp-obsidian-http
    ports:
      - "8080:8080"
    environment:
      - OBSIDIAN_API_KEY=${OBSIDIAN_API_KEY}
    volumes:
      - /path/to/obsidian/vault:/vault:ro
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  default:
    name: mcp-obsidian-network
    driver: bridge
```

## Support

For issues and questions:
1. Check the [main README.md](README.md) for general usage
2. Review the [OBSIDIAN_API_DOCUMENTATION.md](OBSIDIAN_API_DOCUMENTATION.md) for API details
3. Check the [todo.md](todo.md) for current status and known issues
