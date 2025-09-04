# Release Notes v1.0.2

## 🚀 Major Enhancements: JSON Output & Target Validation

This release transforms the MCP Obsidian server into a much more LLM-friendly tool with machine-readable outputs and improved error handling.

## ✨ New Features

### 1. JSON Output for All Structure Tools

#### `obsidian_discover_structure`
- **Converted to JSON format** for better LLM consumption
- **Added `max_depth` parameter** (default: 3, use 0 for unlimited)
- **Cleaner structure** without emoji-heavy human formatting
- **Machine-readable targets** for patch operations

**Before:**
```
📄 Markdown Structure for: document.md
📊 19 heading(s), 2 paragraph(s), 14 list(s)
🎯 PATCH-FRIENDLY TARGETS:
📝 HEADINGS (target_type: heading):
• Step 1: Complete comprehensive inventory: `step 1: complete comprehensive inventory`
```

**After:**
```json
{
  "filepath": "document.md",
  "max_depth": 3,
  "patch_targets": {
    "headings": [
      {
        "title": "Step 1: Complete comprehensive inventory",
        "target": "Step 1: Complete comprehensive inventory",
        "level": 2,
        "line": 15
      }
    ]
  }
}
```

#### `obsidian_list_files_in_dir`
- **Converted to JSON format** for consistency
- **Added `max_depth` parameter** (default: 1, use 0 for unlimited)
- **Clear type indicators** for files vs directories
- **Rich metadata** including size and modification times

**Before:**
```
Files in directory: Tasks
1. AWS_Cost_Analysis
2. AWS_Security_Audit
```

**After:**
```json
{
  "directory": "Tasks",
  "max_depth": 1,
  "total_items": 2,
  "items": [
    {
      "name": "AWS_Cost_Analysis",
      "type": "directory",
      "path": "Tasks/AWS_Cost_Analysis",
      "has_children": true,
      "children_depth": 0
    }
  ]
}
```

### 2. Enhanced Target Validation

#### `obsidian_patch_content`
- **Target existence validation** before API calls
- **Clear error messages** with available headings
- **Case-insensitive matching** for better usability
- **Helpful debugging** when targets are not found

**Before (Generic Error):**
```
failed to patch content: invalid-target
```

**After (Helpful Error):**
```
target heading 'Completed Steps' not found in file. Available headings: [Completed Steps, Step 1: Inventory, Step 2: Analysis, ...]
```

### 3. Improved Emoji Handling

#### `obsidian_get_nested_content`
- **Exact emoji matching** before fallback to cleaned matching
- **Patch target display** showing exact text to use
- **Consistent target format** between get and patch operations

**New Output:**
```
📝 PATCH TARGET: Use this exact target for patch_content operations:
Target: "📋 Current Status"
```

## 🔧 Technical Improvements

### Performance
- **Depth control** reduces processing time for large files
- **Efficient JSON building** with proper Unicode handling
- **Target validation** prevents unnecessary API calls

### LLM Compatibility
- **Machine-readable outputs** for all structure tools
- **Consistent JSON format** across all tools
- **Clear error messages** for debugging
- **Exact target specifications** for patch operations

### Error Handling
- **Target validation** with helpful error messages
- **Available options display** when targets not found
- **Unicode escape handling** for proper character display
- **Graceful fallbacks** for edge cases

## 📋 Usage Examples

### Discover Structure with Depth Control
```bash
# Default depth (3 levels)
obsidian_discover_structure --filepath "document.md"

# Unlimited depth
obsidian_discover_structure --filepath "document.md" --max_depth "0"

# Custom depth
obsidian_discover_structure --filepath "document.md" --max_depth "5"
```

### List Files with Depth Control
```bash
# Default (non-recursive)
obsidian_list_files_in_dir --dirpath "Tasks"

# Recursive listing
obsidian_list_files_in_dir --dirpath "Tasks" --max_depth "3"
```

### Improved Patch Workflow
```bash
# Step 1: Discover structure
obsidian_discover_structure --filepath "document.md"

# Step 2: Get exact target
obsidian_get_nested_content --filepath "document.md" --nested_path "Section Name"

# Step 3: Use exact target in patch
obsidian_patch_content --filepath "document.md" --operation "append" --target_type "heading" --target "EXACT_TARGET_FROM_STEP_2" --content "new content"
```

## 🎯 Impact

### For LLMs
- **Easier parsing** with JSON outputs
- **Better error handling** with clear messages
- **Consistent workflows** across all tools
- **Reduced confusion** with exact target specifications

### For Developers
- **Machine-readable outputs** for automation
- **Clear debugging** with helpful error messages
- **Flexible depth control** for performance
- **Consistent API** across all tools

## 🔄 Migration Guide

### Existing Workflows
- **No breaking changes** - all existing functionality preserved
- **Enhanced outputs** - same tools, better results
- **Optional parameters** - max_depth defaults to sensible values
- **Backward compatible** - old workflows still work

### New Capabilities
- **JSON outputs** for programmatic processing
- **Depth control** for performance optimization
- **Target validation** for better error handling
- **Exact target specifications** for reliable patching

## 🚀 Installation

```bash
# Download the latest release
git clone https://github.com/manishiitg/mcp-obsidian.git
cd mcp-obsidian
git checkout v1.0.2

# Build the binary
go build -o mcp-obsidian .
```

## 📊 Performance Notes

- **Default depth limits** improve processing speed
- **JSON outputs** reduce parsing overhead for LLMs
- **Target validation** prevents failed API calls
- **Efficient Unicode handling** maintains performance

## 🔮 Future Enhancements

- **Recursive directory listing** implementation
- **Advanced search capabilities** with JSON queries
- **Batch operations** for multiple files
- **Real-time file watching** for live updates

---

**This release represents a significant step forward in making the MCP Obsidian server more LLM-friendly and developer-friendly with machine-readable outputs and improved error handling.**
