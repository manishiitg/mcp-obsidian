# MCP Obsidian Error Handling Demonstration

## Current Error Handling Pattern

The MCP Obsidian server uses a **robust error handling pattern** that **does NOT crash** when API calls fail. Instead, it returns structured error responses.

## How It Works

### 1. **API Level (HTTP Client)**
```go
// In obsidian/client/client.go - makeRequest method
if resp.StatusCode >= 400 {
    defer resp.Body.Close()
    bodyBytes, _ := io.ReadAll(resp.Body)
    
    var obsidianError types.ObsidianError
    if err := json.Unmarshal(bodyBytes, &obsidianError); err == nil {
        return nil, fmt.Errorf("API error %d: %s", obsidianError.ErrorCode, obsidianError.Message)
    }
    
    return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
}
```

### 2. **Handler Level (MCP Tools)**
```go
// Example from ListFilesInVault handler
files, err := obsidianClient.ListFilesInVault()
if err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("failed to list files in vault: %v", err)), nil
}
```

### 3. **Error Response Format**
When a tool fails, it returns:
```json
{
  "error": {
    "message": "failed to list files in vault: API error 404: Vault not found"
  }
}
```

## Error Scenarios

### **Network/Connection Errors**
- Connection refused
- Timeout
- SSL/TLS errors

### **HTTP Status Errors**
- 401: Unauthorized (invalid API key)
- 404: Not found (file/vault doesn't exist)
- 500: Internal server error
- 503: Service unavailable

### **API-Specific Errors**
- Invalid parameters
- Missing required fields
- Permission denied
- Rate limiting

## **Special Case: PutContent to Non-Existent Path**

### **What Happens When You Put Content to a Non-Existent File/Folder:**

1. **MCP Client** calls `obsidian_put_content` with `filepath: "non/existent/file.md"`
2. **Handler** validates parameters and creates Obsidian client
3. **Client** sends PUT request to `/vault/non/existent/file.md`
4. **Obsidian REST API** behavior depends on the API implementation:
   - **If API creates missing folders**: File gets created successfully
   - **If API requires folders to exist**: Returns 404 or 400 error
5. **Error captured** and logged with full context
6. **Tool returns** `mcp.NewToolResultError()` with specific error message
7. **MCP Client** receives structured error response

### **Example Error Response:**
```json
{
  "error": {
    "message": "failed to put content to non/existent/file.md: API error 404: Directory does not exist"
  }
}
```

### **Log Entry:**
```json
{
  "level": "error",
  "message": "Tool execution failed",
  "tool_name": "obsidian_put_content",
  "filepath": "non/existent/file.md",
  "error": "failed to put content to non/existent/file.md: API error 404: Directory does not exist",
  "timestamp": "2025-08-10T21:XX:XX.XXXXXXZ"
}
```

## Benefits of This Approach

✅ **No crashes** - Server continues running  
✅ **Structured errors** - Clear error messages for debugging  
✅ **Graceful degradation** - Other tools continue to work  
✅ **Comprehensive logging** - All errors are logged with context  
✅ **MCP compliant** - Proper error responses for MCP clients  

## Example Error Flow

1. **Client calls tool** → `obsidian_put_content`
2. **Handler processes request** → Creates Obsidian client
3. **API call fails** → 404 error (folder doesn't exist)
4. **Error captured** → Detailed error message created
5. **Tool returns error** → `mcp.NewToolResultError()` with message
6. **Error logged** → Full context written to log file
7. **Server continues** → Ready for next request

## Testing Error Handling

To test error handling, you can:
1. **Disconnect network** - Simulate connection failures
2. **Use invalid API key** - Trigger 401 errors
3. **Request non-existent files** - Generate 404 errors
4. **Put content to non-existent paths** - Test folder creation behavior
5. **Check logs** - See detailed error context

The server will handle all these gracefully and provide useful error messages to MCP clients.
