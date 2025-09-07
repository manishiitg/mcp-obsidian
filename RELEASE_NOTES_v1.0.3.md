# Release Notes v1.0.3

## 🚀 Critical Bug Fixes & Enhanced Documentation

This release addresses critical issues with `obsidian_patch_content` and `obsidian_list_files_in_dir` tools, along with comprehensive documentation improvements.

## 🐛 Critical Bug Fixes

### 1. Fixed `obsidian_patch_content` Invalid Target Errors

#### Problem
- **Error 40080**: "invalid-target" failures when patching content
- **Root Cause**: Target string mismatch between user input and actual file content
- **Impact**: Tool failed "most of the time" with cryptic error messages

#### Solution
- **Exact heading matching**: Prioritizes exact string matches before case-insensitive fallback
- **Improved target validation**: Better error messages with available headings
- **Enhanced debugging**: Detailed error information for troubleshooting

**Before:**
```
failed to patch content: invalid-target (error code 40080)
```

**After:**
```
target heading 'Completed Steps' not found in file. Available headings: [Completed Steps, Step 1: Inventory, Step 2: Analysis, ...]
```

### 2. Fixed `obsidian_list_files_in_dir` Directory Path Issues

#### Problem
- **Error 40400**: "Not Found" when using trailing slashes in directory paths
- **Root Cause**: Client was forcing trailing slash in endpoint construction
- **Impact**: `Tasks/AWS_Security_Audit/` failed while `Tasks/AWS_Security_Audit` worked

#### Solution
- **Removed forced trailing slash**: Endpoint now respects user's input format
- **Flexible path handling**: Supports both formats seamlessly
- **Improved error parsing**: Better error message handling

**Before:**
```
❌ Tasks/AWS_Security_Audit/ (with trailing slash) → 40400 Not Found
✅ Tasks/AWS_Security_Audit (without trailing slash) → Success
```

**After:**
```
✅ Tasks/AWS_Security_Audit/ (with trailing slash) → Success
✅ Tasks/AWS_Security_Audit (without trailing slash) → Success
```

## 📚 Enhanced Documentation

### 1. Comprehensive Tool Descriptions

#### `obsidian_patch_content`
- **Multiline descriptions** with detailed examples
- **Target type explanations**: Heading, Block, Frontmatter with usage examples
- **Operation guides**: Replace, Append, Prepend with clear use cases
- **Best practices**: Tips for reliable content patching
- **Error handling**: Understanding and resolving common issues

#### `obsidian_discover_structure`
- **Detailed overview** of structure discovery capabilities
- **Output format explanations**: JSON structure and field meanings
- **Best practices**: When and how to use depth control
- **Examples**: Real-world usage scenarios

#### `obsidian_get_nested_content`
- **Path navigation guide**: Understanding nested path formats
- **Content extraction**: How to target specific sections
- **Patch target display**: Exact strings for patch operations
- **Tips**: Optimizing content retrieval

### 2. Updated Configuration Files

- **API key updates**: All test files and documentation updated with new key
- **Environment setup**: Improved test environment configuration
- **Documentation examples**: Updated curl commands and usage examples

## 🧪 Testing Improvements

### 1. New Test Suite
- **Dedicated directory listing test**: `cmd/test-list-files-in-dir.go`
- **Comprehensive target testing**: Enhanced validation for all target types
- **Error scenario testing**: Better coverage of edge cases

### 2. Enhanced Test Environment
- **Updated API keys**: All test files synchronized with new credentials
- **Improved error reporting**: Better debugging information in tests
- **Comprehensive coverage**: Testing both success and failure scenarios

## 🔧 Technical Improvements

### Code Quality
- **Better error handling**: More descriptive error messages
- **Improved validation**: Target existence checks before API calls
- **Enhanced debugging**: Detailed logging for troubleshooting

### API Compatibility
- **Flexible path handling**: Respects user input format
- **Better error parsing**: Improved error message extraction
- **Consistent behavior**: Reliable operation across different scenarios

## 📋 Usage Examples

### Reliable Content Patching
```bash
# Step 1: Discover structure
obsidian_discover_structure --filepath "document.md"

# Step 2: Get exact target
obsidian_get_nested_content --filepath "document.md" --nested_path "Section Name"

# Step 3: Patch with exact target
obsidian_patch_content --filepath "document.md" --operation "append" --target_type "heading" --target "EXACT_TARGET" --content "new content"
```

### Flexible Directory Listing
```bash
# Both formats now work
obsidian_list_files_in_dir --dirpath "Tasks/AWS_Security_Audit"
obsidian_list_files_in_dir --dirpath "Tasks/AWS_Security_Audit/"
```

### Frontmatter Operations
```bash
# Update frontmatter
obsidian_patch_content --filepath "document.md" --operation "replace" --target_type "frontmatter" --target "frontmatter" --content "status: completed\ndate: 2025-01-01"
```

## 🎯 Impact

### For Users
- **Reliable operation**: Tools work consistently without cryptic errors
- **Better error messages**: Clear guidance when things go wrong
- **Flexible usage**: Support for different path formats
- **Comprehensive documentation**: Easy to understand and use

### For Developers
- **Improved debugging**: Better error information for troubleshooting
- **Enhanced testing**: Comprehensive test coverage
- **Better maintainability**: Cleaner code with better error handling
- **Documentation**: Clear examples and best practices

## 🔄 Migration Guide

### No Breaking Changes
- **Backward compatible**: All existing workflows continue to work
- **Enhanced functionality**: Same tools, better reliability
- **Improved error handling**: Better user experience
- **Updated documentation**: More comprehensive guides

### New Capabilities
- **Flexible path handling**: Support for different directory path formats
- **Better error messages**: Clearer guidance for troubleshooting
- **Enhanced documentation**: Comprehensive usage guides
- **Improved testing**: Better coverage and debugging

## 🚀 Installation

```bash
# Download the latest release
git clone https://github.com/manishiitg/mcp-obsidian.git
cd mcp-obsidian
git checkout v1.0.3

# Build the binary
go build -o mcp-obsidian .
```

## 📊 Performance Notes

- **No performance impact**: Fixes maintain existing performance
- **Better error handling**: Faster debugging with improved messages
- **Flexible path handling**: No overhead for different formats
- **Enhanced testing**: Better coverage without performance cost

## 🔮 Future Enhancements

- **Advanced error handling**: More sophisticated error recovery
- **Enhanced validation**: Better input validation and sanitization
- **Performance optimization**: Further improvements to response times
- **Extended testing**: More comprehensive test coverage

---

**This release focuses on reliability and usability, fixing critical bugs that were causing tool failures and providing comprehensive documentation for better user experience.**
