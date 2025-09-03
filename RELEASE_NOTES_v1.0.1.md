# Release Notes v1.0.1

## 🐛 Bug Fix: Patch Target Format Mismatch

### Issue
The `obsidian_discover_structure` tool was processing heading text through `removeEmojisAndSpecialChars()` which was lowercasing and cleaning the text, but the `obsidian_patch_content` operation needs the **exact original heading text**. This caused "invalid-target" errors when using targets from discover_structure in patch operations.

### Fix
- **Modified `createPatchTarget` function** to return exact original heading text
- **Removed `removeEmojisAndSpecialChars` processing** from target generation
- **Ensures discover_structure targets work with patch_content operations**

### Before
```bash
# Discover showed processed targets
• Step 1: Complete comprehensive inventory of core AWS services and resources: `step 1: complete comprehensive inventory of core aws services and resources`

# Patch failed with "invalid-target" error
obsidian_patch_content --target "step 1: complete comprehensive inventory of core aws services and resources"
```

### After
```bash
# Discover now shows exact targets
• Step 1: Complete comprehensive inventory of core AWS services and resources: `Step 1: Complete comprehensive inventory of core AWS services and resources`

# Patch works correctly
obsidian_patch_content --target "Step 1: Complete comprehensive inventory of core AWS services and resources"
```

### Impact
- ✅ **Fixed workflow**: discover_structure → patch_content now works seamlessly
- ✅ **No more "invalid-target" errors** when using targets from discover_structure
- ✅ **Exact target matching** ensures compatibility with Obsidian API
- ✅ **Backward compatible** - existing functionality remains unchanged

### Files Changed
- `obsidian/handlers/obsidian.go` - Fixed `createPatchTarget` function

### Commit
- **Hash**: 845ca56
- **Message**: "fix: patch target format mismatch in discover_structure"

---

## 🚀 Installation

```bash
# Download the latest release
git clone https://github.com/manishiitg/mcp-obsidian.git
cd mcp-obsidian
git checkout v1.0.1

# Build the binary
go build -o mcp-obsidian .
```

## 📋 Testing

To verify the fix works:

1. **Discover structure** of a file with mixed-case headings
2. **Use exact targets** from discover output in patch operations
3. **Verify patch succeeds** without "invalid-target" errors

```bash
# Test the fix
./mcp-obsidian discover_structure --filepath "your-file.md"
./mcp-obsidian patch_content --filepath "your-file.md" --operation "append" --target_type "heading" --target "Exact Heading Text" --content "test content"
```
