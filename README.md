# HTTPie to Postman Collection Converter

A Go command-line tool that converts HTTPie collection files to Postman Collection Format v2.1.0. It aims to unblock users who need to migrate their API collections from HTTPie to Postman or any other tool that supports Postman collections.

## Quick Start

1. Go to the [Releases page](https://github.com/vuon9/httpie-to-postman-converter/releases/latest)
2. Download the binary for your platform
3. Make it executable (Linux/macOS): `chmod +x httpie-to-postman-*`
4. Run: `./httpie-to-postman-* input.json output.postman.json`

## Features

- Converts HTTPie collection JSON files to Postman Collection v2.1.0 format

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest binary for your platform from the [Releases page](https://github.com/vuon9/httpie-to-postman-converter/releases/latest):

**Linux (x64):**
```bash
wget https://github.com/vuon9/httpie-to-postman-converter/releases/latest/download/httpie-to-postman-linux-amd64
chmod +x httpie-to-postman-linux-amd64
# Optional: Move to your PATH
sudo mv httpie-to-postman-linux-amd64 /usr/local/bin/httpie-to-postman
```

**macOS (Intel):**
```bash
wget https://github.com/vuon9/httpie-to-postman-converter/releases/latest/download/httpie-to-postman-darwin-amd64
chmod +x httpie-to-postman-darwin-amd64
# Optional: Move to your PATH
sudo mv httpie-to-postman-darwin-amd64 /usr/local/bin/httpie-to-postman
```

**macOS (Apple Silicon):**
```bash
wget https://github.com/vuon9/httpie-to-postman-converter/releases/latest/download/httpie-to-postman-darwin-arm64
chmod +x httpie-to-postman-darwin-arm64
# Optional: Move to your PATH
sudo mv httpie-to-postman-darwin-arm64 /usr/local/bin/httpie-to-postman
```

**Windows:**
Download `httpie-to-postman-windows-amd64.exe` from the releases page and run it directly.

### Option 2: Build from Source

1. Make sure you have Go installed (version 1.24.3 or later)
2. Clone or download this repository
3. Build the tool:

```bash
cd httpie-to-postman-converter
go build -o httpie-to-postman main.go
```

## Usage

**If you downloaded a pre-built binary:**
```bash
# If you moved it to your PATH
httpie-to-postman <input-httpie-collection> <output-postman-collection>

# If you're running it directly
./httpie-to-postman-* <input-httpie-collection> <output-postman-collection>
```

**If you built from source:**
```bash
./httpie-to-postman <input-httpie-collection> <output-postman-collection>
```

### Example

```bash
# Using binary in PATH
httpie-to-postman collection.json output.postman.json

# Or running directly
./httpie-to-postman-linux-amd64 collection.json output.postman.json
```

### Output

```bash
./httpie-to-postman collection.json output.postman.json
```

### Output

The tool will display:
- Total number of APIs in the input collection
- Number of successfully converted APIs
- Final output filename (with incremental suffix if original exists)

Example output:
```
Migration completed successfully!
Total APIs input: 52
Converted: 52
Output file: output.postman.json
```

If `output.postman.json` already exists, it will create `output_1.postman.json`, `output_2.postman.json`, etc.

## Supported Conversions

### Request Methods
- GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS

### Headers
- All custom headers
- Automatic Content-Type header addition for requests with body
- Proper handling of enabled/disabled headers

### Request Bodies
- JSON bodies with proper formatting
- Form data (URL-encoded and multipart)
- Raw text bodies
- GraphQL queries

### URLs
- Variable substitution (e.g., `{{YOU}}`)
- Query parameter extraction
- Path parameter handling
- Complex URL structures

### Environment Variables
- Automatic extraction from URLs and headers
- Creation of Postman variables section
- Support for nested variable references

## File Format Support

### Input: HTTPie Collection Format
```json
{
  "meta": {
    "format": "httpie",
    "version": "1.0.0"
  },
  "entry": {
    "name": "Collection Name",
    "requests": [...]
  }
}
```

### Output: Postman Collection v2.1.0
```json
{
  "info": {
    "_postman_id": "...",
    "name": "Collection Name",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [...],
  "variable": [...]
}
```

## Error Handling

The tool includes comprehensive error handling for:
- Invalid input file format
- Missing or inaccessible files
- JSON parsing errors
- File system write permissions

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is open source and available under the MIT License.
