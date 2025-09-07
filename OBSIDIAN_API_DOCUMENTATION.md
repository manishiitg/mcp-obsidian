# Obsidian Local REST API - System Endpoints Documentation

This document provides detailed information about the three system endpoints available in the Obsidian Local REST API and the certificate file.

## Quick Summary

### 🔑 Key Points About the Certificate File

The **`obsidian-local-rest-api.crt`** file is a **self-signed SSL certificate** that serves several critical purposes:

1. **🔒 Secure HTTPS Connections**: Enables encrypted communication between your application and the Obsidian API
2. **🏠 Local Development**: Provides SSL for local development environments without needing a CA-signed certificate
3. **🆔 Server Authentication**: Verifies the identity of the API server to prevent man-in-the-middle attacks
4. **🔐 Data Privacy**: Encrypts all data transmitted between client and server

### 📋 Certificate Details
- **Type**: Self-signed X.509 certificate
- **Validity**: 1 year (365 days) from generation
- **Key Type**: RSA 2048-bit
- **Usage**: HTTPS connections only (port 27124)
- **Auto-generation**: Created automatically when plugin is first installed

### 🎯 Why This Matters for Your Application

- **Development**: Your Go application uses `InsecureSkipVerify: true` for local development
- **Production**: For production use, you should load and validate this certificate
- **Security**: The certificate ensures your API calls are encrypted and secure
- **Compliance**: Proper certificate handling is essential for security compliance

## Overview

The Obsidian Local REST API provides three system-level endpoints that are essential for:
1. **Server Information** - Basic details about the API server
2. **Certificate Management** - SSL certificate used by the API
3. **API Documentation** - OpenAPI specification for all available endpoints

## System Endpoints

### 1. GET `/` - Server Information

**Purpose**: Returns basic details about the Obsidian Local REST API server, including version information, authentication status, and certificate validity.

**Endpoint**: `GET https://127.0.0.1:27124/`

**Headers Required**:
```
Authorization: Bearer <your-api-key>
```

**Response Format**: JSON

**Example Request**:
```bash
curl -k -H "Authorization: Bearer 0488e4465eb68ba7cef5c38e35709e83e2579d67845d21be99edf1936ede1e48" \
  "https://127.0.0.1:27124/"
```

**Example Response**:
```json
{
  "status": "OK",
  "manifest": {
    "id": "obsidian-local-rest-api",
    "name": "Local REST API",
    "version": "3.2.0",
    "minAppVersion": "0.12.0",
    "description": "Get, change or otherwise interact with your notes in Obsidian via a REST API.",
    "author": "Adam Coddington",
    "authorUrl": "https://coddingtonbear.net/",
    "isDesktopOnly": true,
    "dir": ".obsidian/plugins/obsidian-local-rest-api"
  },
  "versions": {
    "obsidian": "1.8.10",
    "self": "3.2.0"
  },
  "service": "Obsidian Local REST API",
  "authenticated": true,
  "certificateInfo": {
    "validityDays": 364.94438386574075,
    "regenerateRecommended": false
  },
  "apiExtensions": []
}
```

**Response Fields**:
- `status`: Server status (usually "OK")
- `manifest`: Plugin manifest information
  - `id`: Plugin identifier
  - `name`: Plugin name
  - `version`: Plugin version
  - `minAppVersion`: Minimum Obsidian version required
  - `description`: Plugin description
  - `author`: Plugin author
  - `authorUrl`: Author's website
  - `isDesktopOnly`: Whether plugin is desktop-only
  - `dir`: Plugin directory path
- `versions`: Version information
  - `obsidian`: Obsidian app version
  - `self`: API plugin version
- `service`: Service name
- `authenticated`: Whether the request was authenticated
- `certificateInfo`: SSL certificate information
  - `validityDays`: Days until certificate expires
  - `regenerateRecommended`: Whether certificate regeneration is recommended
- `apiExtensions`: List of API extensions (usually empty)

**Use Cases**:
- Verify API server is running and accessible
- Check plugin version compatibility
- Monitor certificate expiration
- Confirm authentication is working
- Get server status for health checks

### 2. GET `/obsidian-local-rest-api.crt` - SSL Certificate

**Purpose**: Returns the SSL certificate in use by the Obsidian Local REST API. This is a self-signed certificate used for HTTPS connections.

**Endpoint**: `GET https://127.0.0.1:27124/obsidian-local-rest-api.crt`

**Headers Required**:
```
Authorization: Bearer <your-api-key>
```

**Response Format**: PEM certificate (text)

**Example Request**:
```bash
curl -k -H "Authorization: Bearer 0488e4465eb68ba7cef5c38e35709e83e2579d67845d21be99edf1936ede1e48" \
  "https://127.0.0.1:27124/obsidian-local-rest-api.crt"
```

**Example Response**:
```
-----BEGIN CERTIFICATE-----
MIIDRTCCAi2gAwIBAgIBATANBgkqhkiG9w0BAQsFADAiMSAwHgYDVQQDExdPYnNp
ZGlhbiBMb2NhbCBSRVNUIEFQSTAeFw0yNTA4MDkwNTQyMjNaFw0yNjA4MDkwNTQy
MjNaMCIxIDAeBgNVBAMTF09ic2lkaWFuIExvY2FsIFJFU1QgQVBJMIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAusGhreMueps8DCFqb8VxJqzKSVfdEiCr
iz2wX97v413+ahqVL4gxRWJgZ/3xsaLWhWrdvFqUyvN5C4amoscBcvBGSVsW/EEt
bZ5q5MhzGdL8PvCirAezbJLuW0ZMqXIyP0LidYLfKgMQwgWgjgoG40JqQPQqMPqg
sp7qtxZnN8KCQCP1oY3tKCYzkNKzEcahD+K9QThmJXLfnG0Nth97SaHZTjDaW9iY
RPg6v4roiJvzYX72yY/fKsMP+N/u6uymzwuMlbnQFajcZyoEku79jyekAnCDMZ8x
gyyuyi5C9LxbQCxvl1Um5OdWuPwgoOgchLH7MxCRRvMdjtCun03z2QIDAQABo4GF
MIGCMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgLEMDsGA1UdJQQ0MDIG
CCsGAQUFBwMBBggrBgEFBQcDAgYIKwYBBQUHAwMGCCsGAQUFBwMEBggrBgEFBQcD
CDARBglghkgBhvhCAQEEBAMCAPcwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0B
AQsFAAOCAQEAeZKec2t9WSZleZWFpkqnRXevy88NR2XiNB1EUbFLQ2r11AmcCNgF
PioN0fMgCSGV+XbYwTld+1aGJ5C85EdDaVe6VSA0JFdspw9/NjRTDT1QXXRdLL7y
1qN9BEpj0uQBS4+stzZeSiCz1gqulQBe3mTPQRDCG5i8/1tt2MueAbiSxE+A97Xf
N64Z51aXkWfTOYUrtTSj3W4QrW2HR3xq3ZU6P+WNkn62xD7+ann5cUsn1h9ByHmb
fmbmqXNJ5VfmpdPrjy28jrdcHW7adRuIAjfblE2yIhfQJjPsNzAfySzyuSD1Wi5v
FQtnH5NYzOi7l/wcOK9sTuydcKdGa6gd6w==
-----END CERTIFICATE-----
```

**Certificate Details**:
- **Type**: Self-signed X.509 certificate
- **Issuer**: Obsidian Local REST API
- **Subject**: Obsidian Local REST API
- **Validity**: 1 year (365 days)
- **Key Type**: RSA 2048-bit
- **Signature Algorithm**: SHA256 with RSA

**Use Cases**:
- **SSL/TLS Configuration**: Import certificate into client applications for secure connections
- **Certificate Validation**: Verify certificate authenticity and expiration
- **Development Setup**: Configure development environments to trust the certificate
- **Security Auditing**: Review certificate details for security compliance
- **Troubleshooting**: Diagnose SSL connection issues

**Important Notes**:
- This is a **self-signed certificate** generated by the plugin
- The certificate is valid for **1 year** from creation
- Certificate regeneration is recommended when validity is low
- For production use, consider replacing with a proper CA-signed certificate
- The certificate is used for **HTTPS connections only** (port 27124)

### 3. GET `/openapi.yaml` - API Documentation

**Purpose**: Returns the OpenAPI YAML document describing all capabilities and endpoints of the Obsidian Local REST API.

**Endpoint**: `GET https://127.0.0.1:27124/openapi.yaml`

**Headers Required**:
```
Authorization: Bearer <your-api-key>
```

**Response Format**: YAML (OpenAPI 3.0 specification)

**Example Request**:
```bash
curl -k -H "Authorization: Bearer 0488e4465eb68ba7cef5c38e35709e83e2579d67845d21be99edf1936ede1e48" \
  "https://127.0.0.1:27124/openapi.yaml"
```

**Response Content**: Complete OpenAPI 3.0 specification including:
- All available endpoints
- Request/response schemas
- Authentication methods
- Error codes and messages
- Examples and usage patterns
- Data types and structures

**Use Cases**:
- **API Discovery**: Discover all available endpoints and their capabilities
- **Code Generation**: Generate client libraries and SDKs
- **Documentation**: Create interactive API documentation
- **Testing**: Generate test cases and mock responses
- **Integration**: Understand API structure for third-party integrations
- **Development**: Reference for implementing new features

**Key Sections in OpenAPI Spec**:
1. **Info**: API metadata, version, description
2. **Servers**: Available server configurations (HTTP/HTTPS)
3. **Security**: Authentication requirements
4. **Paths**: All available endpoints with methods
5. **Components**: Reusable schemas and parameters
6. **Tags**: Endpoint categorization (Vault Files, Search, etc.)

## Certificate Management

### Certificate Purpose

The SSL certificate (`obsidian-local-rest-api.crt`) serves several important purposes:

1. **Secure Communication**: Enables HTTPS connections to the API
2. **Local Development**: Provides SSL for local development environments
3. **Authentication**: Verifies the identity of the API server
4. **Privacy**: Encrypts data transmitted between client and server

### Certificate Lifecycle

1. **Generation**: Certificate is automatically generated when the plugin is first installed
2. **Validity**: Certificate is valid for 365 days from generation
3. **Renewal**: Certificate should be regenerated before expiration
4. **Replacement**: Can be replaced with a custom certificate if needed

### Certificate Usage in Applications

#### Go Applications
```go
// Load the certificate for custom TLS configuration
cert, err := os.ReadFile("obsidian-local-rest-api.crt")
if err != nil {
    log.Fatal(err)
}

// Create certificate pool
certPool := x509.NewCertPool()
certPool.AppendCertsFromPEM(cert)

// Configure TLS
tlsConfig := &tls.Config{
    RootCAs: certPool,
}

// Use in HTTP client
client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: tlsConfig,
    },
}
```

#### JavaScript/Node.js Applications
```javascript
const https = require('https');
const fs = require('fs');

// Load certificate
const cert = fs.readFileSync('obsidian-local-rest-api.crt');

// Configure agent
const agent = new https.Agent({
    ca: cert,
    rejectUnauthorized: false // For self-signed certificates
});

// Use in requests
const response = await fetch('https://127.0.0.1:27124/', {
    agent: agent,
    headers: {
        'Authorization': 'Bearer your-api-key'
    }
});
```

#### Python Applications
```python
import requests
import certifi

# For self-signed certificates, disable verification (development only)
response = requests.get(
    'https://127.0.0.1:27124/',
    headers={'Authorization': 'Bearer your-api-key'},
    verify=False  # Only for development!
)

# For production, add certificate to trusted store
response = requests.get(
    'https://127.0.0.1:27124/',
    headers={'Authorization': 'Bearer your-api-key'},
    verify='path/to/obsidian-local-rest-api.crt'
)
```

## Integration Examples

### Health Check Script
```bash
#!/bin/bash

# Check if Obsidian API is running
response=$(curl -s -k -H "Authorization: Bearer $OBSIDIAN_API_KEY" \
  "https://127.0.0.1:27124/")

if [[ $response == *"\"status\":\"OK\""* ]]; then
    echo "✅ Obsidian API is running"
    exit 0
else
    echo "❌ Obsidian API is not responding"
    exit 1
fi
```

### Certificate Monitor
```bash
#!/bin/bash

# Check certificate expiration
response=$(curl -s -k -H "Authorization: Bearer $OBSIDIAN_API_KEY" \
  "https://127.0.0.1:27124/")

# Extract validity days using jq
validity_days=$(echo $response | jq -r '.certificateInfo.validityDays')

if (( $(echo "$validity_days < 30" | bc -l) )); then
    echo "⚠️  Certificate expires in $validity_days days"
    exit 1
else
    echo "✅ Certificate is valid for $validity_days days"
    exit 0
fi
```

### API Documentation Generator
```bash
#!/bin/bash

# Download OpenAPI spec
curl -k -H "Authorization: Bearer $OBSIDIAN_API_KEY" \
  "https://127.0.0.1:27124/openapi.yaml" > obsidian-api.yaml

# Generate documentation using swagger-codegen
swagger-codegen generate -i obsidian-api.yaml -l html2 -o docs/

echo "📚 API documentation generated in docs/"
```

## Troubleshooting

### Common Issues

1. **Certificate Errors**
   - **Problem**: SSL certificate verification fails
   - **Solution**: Import the certificate or disable verification for development
   - **Command**: `curl -k` (disables certificate verification)

2. **Authentication Errors**
   - **Problem**: 401 Unauthorized responses
   - **Solution**: Verify API key is correct and included in Authorization header
   - **Check**: Use the `/` endpoint to verify authentication

3. **Connection Refused**
   - **Problem**: Cannot connect to API server
   - **Solution**: Ensure Obsidian is running and plugin is enabled
   - **Check**: Verify host and port configuration

4. **API Version Mismatch**
   - **Problem**: Features not available or different behavior
   - **Solution**: Check API version using `/` endpoint
   - **Update**: Update Obsidian Local REST API plugin if needed

### Debug Commands

```bash
# Test basic connectivity
curl -k -H "Authorization: Bearer $OBSIDIAN_API_KEY" \
  "https://127.0.0.1:27124/"

# Check certificate
curl -k -H "Authorization: Bearer $OBSIDIAN_API_KEY" \
  "https://127.0.0.1:27124/obsidian-local-rest-api.crt"

# Get API documentation
curl -k -H "Authorization: Bearer $OBSIDIAN_API_KEY" \
  "https://127.0.0.1:27124/openapi.yaml" | head -20
```

## Security Considerations

1. **Self-Signed Certificates**: The default certificate is self-signed and should not be used in production
2. **Local Access Only**: The API is designed for local access only (127.0.0.1)
3. **API Key Security**: Keep API keys secure and don't expose them in client-side code
4. **HTTPS Only**: Always use HTTPS in production environments
5. **Certificate Rotation**: Regularly rotate certificates for security

## Conclusion

These three system endpoints provide essential information for:
- **Server monitoring and health checks**
- **SSL certificate management**
- **API documentation and discovery**
- **Development and integration support**

Understanding and utilizing these endpoints is crucial for building robust applications that integrate with the Obsidian Local REST API.
