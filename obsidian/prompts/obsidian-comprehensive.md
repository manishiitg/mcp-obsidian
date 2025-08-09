# Obsidian MCP Guide - Advanced Markdown Structure & Content Management

This guide covers the core advanced capabilities of the Obsidian MCP server: **markdown structure discovery**, **nested content reading**, and **content patching**.

## 🎯 Core Features

### 1. Markdown Structure Discovery (`obsidian_discover_structure`)

Comprehensive analysis of markdown files to understand structure and generate patch-friendly targets.

```bash
obsidian_discover_structure --filepath "my-document.md"
```

**Output includes:**
- Structure summary (headings, paragraphs, lists, code blocks)
- Patch-friendly targets for headings, blocks, frontmatter
- Detailed nested structure with proper indentation
- Block references and IDs

### 2. Nested Content Reading (`obsidian_get_nested_content`)

Precise content retrieval using nested path navigation.

```bash
obsidian_get_nested_content --filepath "my-document.md" --nested_path "Troubleshooting -> Common Issues"
```

**Path formats:**
- Arrow notation: `"Section -> Subsection -> Sub-subsection"`
- Colon notation: `"Section::Subsection::Sub-subsection"`

### 3. Content Patching (`obsidian_patch_content`)

Surgical content modifications using exact targets.

```bash
# Append to heading
obsidian_patch_content --filepath "my-document.md" --operation "append" --target_type "heading" --target "Troubleshooting" --content "## New Section\n\nContent here."

# Prepend to block
obsidian_patch_content --filepath "my-document.md" --operation "prepend" --target_type "block" --target "abc123" --content "New content."

# Replace frontmatter
obsidian_patch_content --filepath "my-document.md" --operation "replace" --target_type "frontmatter" --target "tags" --content "['updated', 'tags']"
```

**Target types:** headings, blocks, frontmatter
**Operations:** append, prepend, replace

## 📊 Examples

### Structure Discovery
```bash
obsidian_discover_structure --filepath "comprehensive-test.md"
```
**Sample output:**
```
📄 Markdown Structure for: comprehensive-test.md
📊 8 heading(s), 15 paragraph(s), 3 list(s), 2 code block(s)

🎯 PATCH-FRIENDLY TARGETS:
📝 HEADINGS (target_type: heading):
- Comprehensive Test Document
- Quick Summary -> Key Points About Testing
- Troubleshooting -> Common Issues

🔗 BLOCKS (target_type: block):
- 2d9b4a (line 15: Comprehensive Test Document)
- abc123 (line 25: Quick Summary)

📋 FRONTMATTER (target_type: frontmatter):
- title, tags, created
```

### Nested Content Reading
```bash
obsidian_get_nested_content --filepath "comprehensive-test.md" --nested_path "Troubleshooting -> Common Issues"
```
**Sample output:**
```
📖 Content from: Troubleshooting -> Common Issues
- Invalid target format errors
- Missing file errors
- API connection issues
```

### Content Patching
```bash
# Append to heading
obsidian_patch_content \
  --filepath "my-document.md" \
  --operation "append" \
  --target_type "heading" \
  --target "Troubleshooting" \
  --content "## Additional Steps\n\n1. Check logs\n2. Verify permissions"
```

## 🎯 Best Practices

### Workflow
1. **Discover structure** first with `obsidian_discover_structure`
2. **Use generated targets** for precise operations
3. **Test nested paths** with `obsidian_get_nested_content` before patching
4. **Backup** before bulk changes

### Target Strategy
- Use exact heading names or block IDs
- Test operations with small changes first
- Respect existing document hierarchy
- Handle errors and validate results

### Path Navigation
- Use descriptive paths that clearly identify targets
- Test paths incrementally
- Check if targets exist before operations
- Use consistent delimiters (-> or ::)

## 🔍 Advanced Search

### Content Selection
```bash
obsidian_read_content --filepath "my-document.md" --selector_type "heading" --query "Troubleshooting" --exact "true"
```

### JSON-Based Search
```json
{
  "and": [
    {"in": ["#project", {"var": "tags"}]},
    {"!": {"var": "archived"}}
  ]
}
```

## 🚀 Integration Patterns

### Automated Documentation Updates
1. Discover structure of existing docs
2. Read current content using nested paths
3. Generate updated content
4. Patch content using precise targets
5. Validate results

### Content Migration
1. Read source content using nested paths
2. Transform content as needed
3. Discover target document structure
4. Patch content to appropriate locations

## 🔍 Troubleshooting

### Common Issues
- **Target not found**: Use `obsidian_discover_structure` to verify targets
- **Nested path errors**: Check path syntax and ensure levels exist
- **Patching failures**: Verify target type and operation
- **Structure parsing**: Check markdown syntax

### Performance
- Use specific selectors for large documents
- Batch operations when possible
- Cache structure discovery results
- Make incremental updates

## 🎉 Key Takeaways

The Obsidian MCP server excels at:
1. **Precise Content Management**: Nested paths and exact targets for surgical operations
2. **Structure Intelligence**: Automatic discovery and analysis of markdown structure
3. **Advanced Patching**: Three target types with three operations
4. **Comprehensive Testing**: Full test suite with real examples

**Ideal for:** Automated documentation, content migration, knowledge base organization, advanced note-taking workflows, and integration with other MCP servers.
