package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// HTTPie Collection Structures
type HTTPieCollection struct {
	Meta  HTTPieMeta  `json:"meta"`
	Entry HTTPieEntry `json:"entry"`
}

type HTTPieMeta struct {
	Format      string `json:"format"`
	Version     string `json:"version"`
	ContentType string `json:"contentType"`
	Schema      string `json:"schema"`
	Docs        string `json:"docs"`
	Source      string `json:"source"`
}

type HTTPieEntry struct {
	Name     string          `json:"name"`
	Icon     HTTPieIcon      `json:"icon"`
	Auth     HTTPieAuth      `json:"auth"`
	Requests []HTTPieRequest `json:"requests"`
}

type HTTPieIcon struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type HTTPieAuth struct {
	Type string `json:"type"`
}

type HTTPieRequest struct {
	Name        string             `json:"name"`
	URL         string             `json:"url"`
	Method      string             `json:"method"`
	Headers     []HTTPieHeader     `json:"headers"`
	QueryParams []HTTPieQueryParam `json:"queryParams"`
	PathParams  []HTTPiePathParam  `json:"pathParams"`
	Auth        HTTPieAuth         `json:"auth"`
	Body        HTTPieBody         `json:"body"`
}

type HTTPieHeader struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type HTTPieQueryParam struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type HTTPiePathParam struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type HTTPieBody struct {
	Type    string        `json:"type"`
	File    HTTPieFile    `json:"file"`
	Text    HTTPieText    `json:"text"`
	Form    HTTPieForm    `json:"form"`
	GraphQL HTTPieGraphQL `json:"graphql"`
}

type HTTPieFile struct {
	Name string `json:"name"`
}

type HTTPieText struct {
	Value  string `json:"value"`
	Format string `json:"format"`
}

type HTTPieForm struct {
	IsMultipart bool              `json:"isMultipart"`
	Fields      []HTTPieFormField `json:"fields"`
}

type HTTPieFormField struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type HTTPieGraphQL struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// Postman Collection v2.1.0 Structures
type PostmanCollection struct {
	Info     PostmanInfo       `json:"info"`
	Item     []PostmanItem     `json:"item"`
	Variable []PostmanVariable `json:"variable,omitempty"`
}

type PostmanInfo struct {
	PostmanID   string `json:"_postman_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

type PostmanItem struct {
	Name    string         `json:"name"`
	Request PostmanRequest `json:"request"`
}

type PostmanRequest struct {
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header,omitempty"`
	Body   *PostmanBody    `json:"body,omitempty"`
	URL    PostmanURL      `json:"url"`
}

type PostmanHeader struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

type PostmanBody struct {
	Mode    string              `json:"mode"`
	Raw     string              `json:"raw,omitempty"`
	Options *PostmanBodyOptions `json:"options,omitempty"`
}

type PostmanBodyOptions struct {
	Raw PostmanBodyRaw `json:"raw"`
}

type PostmanBodyRaw struct {
	Language string `json:"language"`
}

type PostmanURL struct {
	Raw   string              `json:"raw"`
	Host  []string            `json:"host,omitempty"`
	Path  []string            `json:"path,omitempty"`
	Query []PostmanQueryParam `json:"query,omitempty"`
}

type PostmanQueryParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: httpie-to-postman <input-httpie-collection> <output-postman-collection>")
		fmt.Println("Example: httpie-to-postman collection.json output.postman.json")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputPath := os.Args[2]

	// Read HTTPie collection
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var httpieCollection HTTPieCollection
	if err := json.Unmarshal(data, &httpieCollection); err != nil {
		log.Fatalf("Error parsing HTTPie collection: %v", err)
	}

	// Convert to Postman collection
	postmanCollection := convertToPostman(httpieCollection)

	// Generate unique output filename if file exists
	finalOutputPath := generateUniqueFilename(outputPath)

	// Write Postman collection
	outputData, err := json.MarshalIndent(postmanCollection, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling Postman collection: %v", err)
	}

	if err := os.WriteFile(finalOutputPath, outputData, 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	// Print results
	totalInputAPIs := len(httpieCollection.Entry.Requests)
	convertedAPIs := len(postmanCollection.Item)

	fmt.Printf("Migration completed successfully!\n")
	fmt.Printf("Total APIs input: %d\n", totalInputAPIs)
	fmt.Printf("Converted: %d\n", convertedAPIs)
	fmt.Printf("Output file: %s\n", finalOutputPath)
}

func convertToPostman(httpie HTTPieCollection) PostmanCollection {
	postman := PostmanCollection{
		Info: PostmanInfo{
			PostmanID:   generatePostmanID(),
			Name:        httpie.Entry.Name,
			Description: "Converted from HTTPie collection",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item:     []PostmanItem{},
		Variable: extractVariables(httpie),
	}

	for _, req := range httpie.Entry.Requests {
		postmanItem := convertRequest(req)
		postman.Item = append(postman.Item, postmanItem)
	}

	return postman
}

func convertRequest(httpieReq HTTPieRequest) PostmanItem {
	postmanReq := PostmanRequest{
		Method: httpieReq.Method,
		Header: convertHeaders(httpieReq.Headers),
		URL:    convertURL(httpieReq.URL),
	}

	// Convert body if present
	if httpieReq.Body.Type != "none" && httpieReq.Body.Text.Value != "" {
		postmanReq.Body = &PostmanBody{
			Mode: "raw",
			Raw:  httpieReq.Body.Text.Value,
		}

		// Set options for JSON content
		if httpieReq.Body.Text.Format == "application/json" {
			postmanReq.Body.Options = &PostmanBodyOptions{
				Raw: PostmanBodyRaw{
					Language: "json",
				},
			}
		}

		// Add Content-Type header if not present and body has format
		if httpieReq.Body.Text.Format != "" {
			hasContentType := false
			for _, header := range postmanReq.Header {
				if strings.ToLower(header.Key) == "content-type" {
					hasContentType = true
					break
				}
			}
			if !hasContentType {
				postmanReq.Header = append(postmanReq.Header, PostmanHeader{
					Key:   "Content-Type",
					Value: httpieReq.Body.Text.Format,
				})
			}
		}
	}

	return PostmanItem{
		Name:    httpieReq.Name,
		Request: postmanReq,
	}
}

func convertHeaders(httpieHeaders []HTTPieHeader) []PostmanHeader {
	var postmanHeaders []PostmanHeader

	for _, header := range httpieHeaders {
		postmanHeader := PostmanHeader{
			Key:   header.Name,
			Value: header.Value,
			Type:  "text",
		}

		if !header.Enabled {
			postmanHeader.Disabled = true
		}

		postmanHeaders = append(postmanHeaders, postmanHeader)
	}

	return postmanHeaders
}

func convertURL(httpieURL string) PostmanURL {
	// Parse the URL
	parsedURL, err := url.Parse(httpieURL)
	if err != nil {
		// If parsing fails, return raw URL
		return PostmanURL{
			Raw: httpieURL,
		}
	}

	// Extract host
	var host []string
	if parsedURL.Host != "" {
		host = []string{parsedURL.Scheme + "://" + parsedURL.Host}
	} else {
		// Handle cases where URL starts with variable like {{YOU}}
		urlParts := strings.Split(httpieURL, "/")
		if len(urlParts) > 0 {
			host = []string{urlParts[0]}
		}
	}

	// Extract path
	var path []string
	if parsedURL.Path != "" && parsedURL.Path != "/" {
		pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		for _, part := range pathParts {
			if part != "" {
				path = append(path, part)
			}
		}
	}

	// Extract query parameters
	var query []PostmanQueryParam
	if parsedURL.RawQuery != "" {
		queryParams, _ := url.ParseQuery(parsedURL.RawQuery)
		for key, values := range queryParams {
			for _, value := range values {
				query = append(query, PostmanQueryParam{
					Key:   key,
					Value: value,
				})
			}
		}
	}

	return PostmanURL{
		Raw:   httpieURL,
		Host:  host,
		Path:  path,
		Query: query,
	}
}

func extractVariables(httpie HTTPieCollection) []PostmanVariable {
	variableSet := make(map[string]bool)
	var variables []PostmanVariable

	// Extract variables from all URLs and headers
	variableRegex := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	for _, req := range httpie.Entry.Requests {
		// Extract from URL
		matches := variableRegex.FindAllStringSubmatch(req.URL, -1)
		for _, match := range matches {
			if len(match) > 1 {
				varName := match[1]
				if !variableSet[varName] {
					variables = append(variables, PostmanVariable{
						Key:   varName,
						Value: "",
						Type:  "string",
					})
					variableSet[varName] = true
				}
			}
		}

		// Extract from headers
		for _, header := range req.Headers {
			matches := variableRegex.FindAllStringSubmatch(header.Value, -1)
			for _, match := range matches {
				if len(match) > 1 {
					varName := match[1]
					if !variableSet[varName] {
						variables = append(variables, PostmanVariable{
							Key:   varName,
							Value: "",
							Type:  "string",
						})
						variableSet[varName] = true
					}
				}
			}
		}
	}

	return variables
}

func generatePostmanID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateUniqueFilename(basePath string) string {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return basePath
	}

	dir := filepath.Dir(basePath)
	ext := filepath.Ext(basePath)
	name := strings.TrimSuffix(filepath.Base(basePath), ext)

	counter := 1
	for {
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, counter, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}
