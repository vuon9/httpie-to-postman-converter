
# Postmanzier

A fast CLI tool to merge and convert HTTPie and Postman collections into Postman Collection v2.1.0 format.

## Installation

- Download the latest binary for your OS from the [releases page](https://github.com/vuon9/postmanzier/releases).
- (Linux/macOS) `chmod +x postmanzier-*`
- Or build from source:
  `go build -o postmanzier main.go`

## Usage

### 1. Convert a Single HTTPie Collection

Convert an HTTPie collection to Postman format.

```bash
postmanzier <input-httpie-collection.json> <output-postman-collection.json>
```

**Input format:**
HTTPie workspace/collection JSON (see below).

**Output:**
A Postman v2.1.0 collection file.

**Example:**
```bash
postmanzier collection.json output.postman.json
```
_Output:_
```
Migration completed!
* Total APIs: 1
* Total problematic APIs: 0
* Total variables: 0
--> Output file: output.postman.json
```

**Supported input format:**
```json
{
  "meta": { "format": "httpie", "version": "1.0.0" },
  "entry": { "name": "Collection Name", "requests": [...] }
}
```

---

### 2. Merge Multiple Collections

Merge multiple HTTPie or Postman collections into a single Postman collection.
The tool auto-detects the format from the first input file.

```bash
postmanzier merge <output-file.json> <input1.json> <input2.json> ...
```

- Each input collection becomes a folder in the output.
- Variables are merged and deduplicated.

**Examples:**

_Merge HTTPie collections:_
```bash
postmanzier merge merged-httpie.postman.json collection1.json collection2.json
```
_Output:_
```
HTTPie collections merge completed!
--> Output file: merged-httpie.postman.json
```

_Merge Postman collections:_
```bash
postmanzier merge merged-postman.postman.json postman_collection1.json postman_collection2.json
```
_Output:_
```
Postman collections merge completed!
--> Output file: merged-postman.postman.json
```

**Supported input formats:**

- HTTPie: see above.
- Postman v2.1.0:
  ```json
  {
    "info": { "_postman_id": "...", "name": "Collection Name", "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json" },
    "item": [...],
    "variable": [...]
  }
  ```

**Output:**
Always Postman Collection v2.1.0.

---

## License

MIT
