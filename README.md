# HTTPie to Postman Collection Converter

A Go command-line tool that converts HTTPie collection files to Postman Collection Format v2.1.0. 

It aims to unblock users who need to migrate their API collections from HTTPie to Postman or any other tool that supports Postman collections.

## Installation

### Option 1: Prebuilt binary
1. Please look at the release page to have your prebuilt binary:

- **Linux (x64)**: `httpie-to-postman-linux-amd64`
- **Linux (ARM64)**: `httpie-to-postman-linux-arm64`
- **macOS (Intel)**: `httpie-to-postman-darwin-amd64`
- **macOS (Apple Silicon)**: `httpie-to-postman-darwin-arm64`
- **Windows (x64)**: `httpie-to-postman-windows-amd64.exe`
- **Windows (ARM64)**: `httpie-to-postman-windows-arm64.exe`

2. Download the binary for your platform
3. Make it executable (Linux/macOS): `chmod +x httpie-to-postman-*`
4. Move to your PATH or run directly

### Option 2: Build from Source

1. Make sure you have Go installed (version 1.24.3 or later)
2. Clone or download this repository
3. Build the tool:

```bash
cd httpie-to-postman-converter
go build -o httpie-to-postman main.go
```

## Usage

```bash
httpie-to-postman collection.json output.postman.json

### Example of output
Migration completed successfully!
Total APIs input: 52
Converted: 52
Output file: output.postman.json
```

If `output.postman.json` already exists, it will create `output_1.postman.json`, `output_2.postman.json`, etc.

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

## License

This project is open source and available under the MIT License.
