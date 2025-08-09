# Obsidian MCP Comprehensive Guide

This guide covers all the advanced capabilities of the Obsidian MCP server, with a focus on **markdown structure discovery**, **nested content reading**, and **content patching** - the core features demonstrated in the comprehensive testing suite.

## 🎯 Core Advanced Features

### 1. Markdown Structure Discovery (`obsidian_discover_structure`)

**Purpose**: Comprehensive analysis of markdown files to understand their complete structure and generate patch-friendly targets.

```bash
# Discover complete structure of a markdown file
obsidian_discover_structure --filepath "my-document.md"
```

**Key Capabilities:**
- **Hierarchical Analysis**: Automatically detects headings, paragraphs, lists, code blocks, tables
- **Nested Structure**: Builds a complete tree of content relationships
- **Patch-Friendly Targets**: Generates exact target identifiers for content patching
- **Content Previews**: Shows previews of each section for easy identification

**Sample Output:**
```
📄 Markdown Structure for: comprehensive-test.md

📊 8 heading(s), 15 paragraph(s), 3 list(s), 2 code block(s)

🎯 PATCH-FRIENDLY TARGETS:
Use these exact targets for patch_content operations:

📝 HEADINGS (target_type: heading):
- Comprehensive Test Document
- Quick Summary
- Quick Summary -> Key Points About Testing
- Troubleshooting
- Troubleshooting -> Common Issues
- System Endpoints
- Advanced Features

🔗 BLOCKS (target_type: block):
- 2d9b4a (line 15: Comprehensive Test Document)
- abc123 (line 25: Quick Summary)
- def456 (line 35: Key Points About Testing)

📋 FRONTMATTER (target_type: frontmatter):
- title
- tags
- created
```

### 2. Nested Content Reading (`obsidian_get_nested_content`)

**Purpose**: Precise content retrieval using nested path navigation - the core feature for reading specific sections without knowing exact line numbers.

```bash
# Read content using nested path navigation
obsidian_get_nested_content --filepath "my-document.md" --nested_path "Troubleshooting -> Common Issues"
```

**Path Formats Supported:**
- **Arrow notation**: `"Section -> Subsection -> Sub-subsection"`
- **Colon notation**: `"Section::Subsection::Sub-subsection"`
- **Mixed notation**: `"Section -> Subsection::Detail"`

**Real Examples from Testing:**
```bash
# Test cases from comprehensive-target-test.go
obsidian_get_nested_content --filepath "comprehensive-test.md" --nested_path "Troubleshooting -> Common Issues"
obsidian_get_nested_content --filepath "comprehensive-test.md" --nested_path "Quick Summary -> Key Points About Testing"
obsidian_get_nested_content --filepath "comprehensive-test.md" --nested_path "System Endpoints -> 1. GET / - Server Information"
```

**Use Cases:**
- Reading specific sections without knowing exact line numbers
- Extracting content from deeply nested structures
- Building content summaries and reports
- Automated content analysis

### 3. Content Patching with Precision (`obsidian_patch_content`)

**Purpose**: Surgical content modifications using exact targets - the advanced feature for precise content updates.

```bash
# Append content to a specific heading
obsidian_patch_content --filepath "my-document.md" --operation "append" --target_type "heading" --target "Troubleshooting" --content "## New Section\n\nAdditional content here."

# Prepend content to a block
obsidian_patch_content --filepath "my-document.md" --operation "prepend" --target_type "block" --target "abc123" --content "New content at the beginning."

# Replace frontmatter
obsidian_patch_content --filepath "my-document.md" --operation "replace" --target_type "frontmatter" --target "tags" --content "['updated', 'tags', 'here']"
```

**Target Types:**
- **Headings**: Target by heading name (e.g., "Introduction", "Troubleshooting::Common Issues")
- **Blocks**: Target by block ID (e.g., "2d9b4a", "abc123")
- **Frontmatter**: Target by frontmatter field (e.g., "title", "tags", "created")

**Operations:**
- **append**: Add content at the end of the target
- **prepend**: Add content at the beginning of the target
- **replace**: Replace the entire target content

## 🔍 Advanced Search & Content Operations

### Complex Content Selection (`obsidian_read_content`)

```bash
# Search with exact matching
obsidian_read_content --filepath "my-document.md" --selector_type "heading" --query "Troubleshooting" --exact "true"

# Search with level filtering
obsidian_read_content --filepath "my-document.md" --selector_type "heading" --query "Common" --level "2"
```

### JSON-Based Complex Search (`obsidian_search_json`)

```json
{
  "and": [
    {"in": ["#project", {"var": "tags"}]},
    {"!": {"var": "archived"}},
    {"<": [{"var": "created"}, "2024-01-01"]}
  ]
}
```

## 📊 Real-World Examples from Testing

### Example 1: Document Structure Discovery

```bash
# Discover structure of a complex document
obsidian_discover_structure --filepath "comprehensive-test.md"
```

This reveals the complete structure including:
- All headings with their hierarchy levels
- Block references and IDs
- Frontmatter fields
- Patch-friendly targets for precise operations

### Example 2: Nested Content Reading

```bash
# Read specific nested content
obsidian_get_nested_content --filepath "comprehensive-test.md" --nested_path "Troubleshooting -> Common Issues"
```

**Sample Output:**
```
📖 Content from: Troubleshooting -> Common Issues

- Invalid target format errors
- Missing file errors
- API connection issues
```

### Example 3: Content Patching Examples

```bash
# Append content to a heading
obsidian_patch_content \
  --filepath "my-document.md" \
  --operation "append" \
  --target_type "heading" \
  --target "Troubleshooting" \
  --content "

## Additional Troubleshooting Steps

1. Check the logs for detailed error messages
2. Verify all required dependencies are installed
3. Ensure proper permissions are set
"

# Update document tags
obsidian_patch_content \
  --filepath "my-document.md" \
  --operation "replace" \
  --target_type "frontmatter" \
  --target "tags" \
  --content "['updated', 'documentation', 'guide']"

# Add content at the beginning of a block
obsidian_patch_content \
  --filepath "my-document.md" \
  --operation "prepend" \
  --target_type "block" \
  --target "2d9b4a" \
  --content "**Important Note:** This section has been updated.\n\n"
```

## 🎯 Best Practices for Advanced Usage

### 1. Structure Discovery Workflow

1. **Start with discovery**: Always use `obsidian_discover_structure` first to understand the document layout
2. **Identify targets**: Use the generated patch-friendly targets for precise operations
3. **Validate paths**: Test nested paths with `obsidian_get_nested_content` before patching
4. **Backup first**: Create backups before making bulk changes

### 2. Content Patching Strategy

1. **Use specific targets**: Prefer exact heading names or block IDs over generic selectors
2. **Test operations**: Start with small changes to verify the target is correct
3. **Maintain structure**: Respect the existing document hierarchy when adding content
4. **Handle errors**: Always check for errors and validate the results

### 3. Nested Path Navigation

1. **Use descriptive paths**: Choose meaningful path names that clearly identify the target
2. **Test paths incrementally**: Build complex paths by testing each level
3. **Handle missing content**: Always check if the target exists before attempting operations
4. **Use consistent delimiters**: Stick to one delimiter style (-> or ::) for consistency

## 🚀 Integration Patterns

### Automated Documentation Updates

```bash
# Automated workflow for updating documentation
1. Discover structure of existing docs
2. Read current content using nested paths
3. Generate updated content based on changes
4. Patch content using precise targets
5. Validate the results
```

### Content Migration and Transformation

```bash
# Migrate content between documents
1. Read source content using nested paths
2. Transform content as needed
3. Discover target document structure
4. Patch content to appropriate locations
5. Update references and links
```

## 🔍 Troubleshooting Advanced Features

### Common Issues and Solutions

1. **Target not found**: Use `obsidian_discover_structure` to verify available targets
2. **Nested path errors**: Check path syntax and ensure all levels exist
3. **Patching failures**: Verify target type and operation are correct
4. **Structure parsing issues**: Check markdown syntax and formatting

### Debugging Techniques

1. **Enable verbose logging**: Check server logs for detailed error information
2. **Test with simple examples**: Start with basic operations before complex ones
3. **Validate input parameters**: Ensure all required parameters are provided
4. **Check file permissions**: Verify the MCP server has access to the files

## 📈 Performance Optimization

### Large Document Handling

1. **Use specific selectors**: Avoid reading entire large documents when possible
2. **Batch operations**: Group related operations for efficiency
3. **Cache structure**: Reuse structure discovery results for multiple operations
4. **Incremental updates**: Make small, focused changes rather than bulk replacements

### Search Optimization

1. **Use exact matches**: When possible, use exact matching for faster searches
2. **Limit scope**: Use level filters and specific selectors to narrow search scope
3. **Index large vaults**: Consider indexing strategies for very large vaults
4. **Optimize queries**: Use specific search terms and avoid overly broad queries

## 🎉 Key Takeaways

The Obsidian MCP server excels at:

1. **Precise Content Management**: Using nested paths and exact targets for surgical content operations
2. **Structure Intelligence**: Automatic discovery and analysis of markdown document structure
3. **Advanced Patching**: Three target types (headings, blocks, frontmatter) with three operations (append, prepend, replace)
4. **Comprehensive Testing**: Full test suite demonstrating all capabilities with real examples

This makes it ideal for:
- Automated documentation management
- Content migration and transformation
- Knowledge base organization
- Advanced note-taking workflows
- Integration with other MCP servers and external systems
