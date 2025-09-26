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
	"github.com/google/uuid"
)

// HTTPie Workspace Structure (new format)
type HTTPieWorkspace struct {
	Meta         HTTPieMeta          `json:"meta"`
	Entry        HTTPieEntry         `json:"entry"`
	Environments []HTTPieEnvironment `json:"environments,omitempty"`
}

type HTTPieCollection struct {
	Name     string          `json:"name"`
	Icon     HTTPieIcon      `json:"icon"`
	Auth     HTTPieAuth      `json:"auth"`
	Requests []HTTPieRequest `json:"requests"`
}

type HTTPieEnvironment struct {
	Name        string                     `json:"name"`
	Color       string                     `json:"color"`
	IsDefault   bool                       `json:"isDefault"`
	IsLocalOnly bool                       `json:"isLocalOnly"`
	Variables   []HTTPieEnvironmentVariable `json:"variables"`
}

type HTTPieEnvironmentVariable struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
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
	Name        string               `json:"name"`
	Icon        HTTPieIcon           `json:"icon"`
	Auth        HTTPieAuth           `json:"auth"`
	Requests    []HTTPieRequest      `json:"requests,omitempty"`
	Collections []HTTPieCollection  `json:"collections,omitempty"`
}

type HTTPieIcon struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type HTTPieAuth struct {
	Type        string                 `json:"type"`
	Target      string                 `json:"target,omitempty"`
	Credentials HTTPieAuthCredentials  `json:"credentials,omitempty"`
}

type HTTPieAuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	Request *PostmanRequest `json:"request,omitempty"`
	Item    []PostmanItem  `json:"item,omitempty"`
}

type PostmanRequest struct {
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header,omitempty"`
	Body   *PostmanBody    `json:"body,omitempty"`
	Auth   *PostmanAuth    `json:"auth,omitempty"`
	URL    PostmanURL      `json:"url"`
}

type PostmanAuth struct {
	Type   string                 `json:"type"`
	Bearer []PostmanAuthBearer    `json:"bearer,omitempty"`
	Basic  []PostmanAuthBasic     `json:"basic,omitempty"`
	APIKey []PostmanAuthAPIKey    `json:"apikey,omitempty"`
}

type PostmanAuthBearer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PostmanAuthBasic struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PostmanAuthAPIKey struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
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
	ID    string `json:"id,omitempty"` // Optional ID for Postman variables
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "merge":
		handleMergeCommand()
	default:
		handleConvertCommand()
	}
}

func handleMergeCommand() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: httpie-to-postman merge <output-file> <input-file-1> [<input-file-2> ...]")
		fmt.Println("Example: httpie-to-postman merge merged.postman.json collection1.json collection2.json")
		os.Exit(1)
	}

	outputFile := os.Args[2]
	inputFiles := os.Args[3:]

	mergedCollection := PostmanCollection{
		Info: PostmanInfo{
			PostmanID:   generatePostmanID(),
			Name:        "Merged HTTPie Collections",
			Description: "Merged from multiple HTTPie collections",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item:     []PostmanItem{},
		Variable: []PostmanVariable{},
	}

	allVariables := make(map[string]string)

	for _, inputFile := range inputFiles {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			log.Printf("Error reading input file %s: %v. Skipping.", inputFile, err)
			continue
		}

		var httpieWorkspace HTTPieWorkspace
		if err := json.Unmarshal(data, &httpieWorkspace); err != nil {
			log.Printf("Error parsing HTTPie collection from %s: %v. Skipping.", inputFile, err)
			continue
		}

		folderName := httpieWorkspace.Entry.Name
		if folderName == "" {
			folderName = strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		}

		folder := PostmanItem{
			Name: folderName,
			Item: []PostmanItem{},
		}

		// Convert direct requests
		for _, req := range httpieWorkspace.Entry.Requests {
			folder.Item = append(folder.Item, convertRequest(req))
		}

		// Convert requests from collections
		for _, collection := range httpieWorkspace.Entry.Collections {
			for _, req := range collection.Requests {
				folder.Item = append(folder.Item, convertRequest(req))
			}
		}

		if len(folder.Item) > 0 {
			mergedCollection.Item = append(mergedCollection.Item, folder)
		}

		// Extract and merge variables
		vars := extractVariablesFromWorkspace(httpieWorkspace)
		for _, v := range vars {
			if _, exists := allVariables[v.Key]; !exists {
				allVariables[v.Key] = v.Value
			}
		}
	}

	for key, value := range allVariables {
		mergedCollection.Variable = append(mergedCollection.Variable, PostmanVariable{
			ID:    uuid.New().String(),
			Key:   key,
			Value: value,
			Type:  "string",
		})
	}

	// Write merged Postman collection
	outputData, err := json.MarshalIndent(mergedCollection, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling merged Postman collection: %v", err)
	}

	finalOutputPath := generateUniqueFilename(outputFile)
	if err := os.WriteFile(finalOutputPath, outputData, 0644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("Merge completed!")
	fmt.Printf("--> Output file: %s\n", finalOutputPath)
}

func handleConvertCommand() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputPath := os.Args[2]

	// Read HTTPie collection
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Try to parse as workspace format first (newer format)
	var httpieWorkspace HTTPieWorkspace
	if err := json.Unmarshal(data, &httpieWorkspace); err != nil {
		log.Fatalf("Error parsing HTTPie collection: %v", err)
	}

	// Convert to Postman collection
	postmanCollection := convertWorkspaceToPostman(httpieWorkspace)

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
	totalInputAPIs := countTotalRequests(httpieWorkspace)
	convertedAPIs := len(postmanCollection.Item)
	totalVariables := len(postmanCollection.Variable)

	errorStr := "\n"
	if (convertedAPIs != totalInputAPIs) {
		errorStr = fmt.Sprintf("Some requests were not converted correctly\n")
	}

	fmt.Printf("Migration completed!%s", errorStr)
	fmt.Printf("* Total APIs: %d\n", totalInputAPIs)
	fmt.Printf("* Total problematic APIs: %d\n", totalInputAPIs - convertedAPIs)
	fmt.Printf("* Total variables: %d\n", totalVariables)
	fmt.Printf("--> Output file: %s\n", finalOutputPath)
}

func printUsage() {
	fmt.Println("Usage: httpie-to-postman <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  <input-httpie-collection> <output-postman-collection>")
	fmt.Println("    Converts a single HTTPie collection to a Postman collection.")
	fmt.Println("    Example: httpie-to-postman collection.json output.postman.json")
	fmt.Println("\n  merge <output-file> <input-file-1> [<input-file-2> ...]")
	fmt.Println("    Merges multiple HTTPie collections into a single Postman collection.")
	fmt.Println("    Example: httpie-to-postman merge merged.postman.json collection1.json collection2.json")
}

func convertWorkspaceToPostman(httpie HTTPieWorkspace) PostmanCollection {
	postman := PostmanCollection{
		Info: PostmanInfo{
			PostmanID:   generatePostmanID(),
			Name:        httpie.Entry.Name,
			Description: "Converted from HTTPie workspace",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item:     []PostmanItem{},
		Variable: extractVariablesFromWorkspace(httpie),
	}

	// Create a folder for the collection content
	folder := PostmanItem{
		Name: httpie.Entry.Name,
		Item: []PostmanItem{},
	}

	// Convert direct requests (if any)
	for _, req := range httpie.Entry.Requests {
		postmanItem := convertRequest(req)
		folder.Item = append(folder.Item, postmanItem)
	}

	// Convert collections (folders)
	for _, collection := range httpie.Entry.Collections {
		// Add requests to this folder - for now we'll flatten them
		// In a more complex version, we could create nested folders
		for _, req := range collection.Requests {
			postmanItem := convertRequest(req)
			folder.Item = append(folder.Item, postmanItem)
		}
	}

	if len(folder.Item) > 0 {
		postman.Item = append(postman.Item, folder)
	}

	return postman
}


func countTotalRequests(httpie HTTPieWorkspace) int {
	count := len(httpie.Entry.Requests)
	for _, collection := range httpie.Entry.Collections {
		count += len(collection.Requests)
	}
	return count
}

func convertToPostman(httpie HTTPieWorkspace) PostmanCollection {
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
		Auth:   convertAuth(httpieReq.Auth),
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

	// Generate a name if empty
	name := httpieReq.Name
	if name == "" {
		name = httpieReq.Method + " " + httpieReq.URL
	}

	return PostmanItem{
		Name:    name,
		Request: &postmanReq,
	}
}

func convertAuth(httpieAuth HTTPieAuth) *PostmanAuth {
	if httpieAuth.Type == "none" || httpieAuth.Type == "" {
		return nil
	}

	switch httpieAuth.Type {
	case "bearer":
		return &PostmanAuth{
			Type: "bearer",
			Bearer: []PostmanAuthBearer{
				{
					Key:   "token",
					Value: httpieAuth.Credentials.Password,
					Type:  "string",
				},
			},
		}
	case "basic":
		return &PostmanAuth{
			Type: "basic",
			Basic: []PostmanAuthBasic{
				{
					Key:   "username",
					Value: httpieAuth.Credentials.Username,
					Type:  "string",
				},
				{
					Key:   "password",
					Value: httpieAuth.Credentials.Password,
					Type:  "string",
				},
			},
		}
	case "apiKey":
		return &PostmanAuth{
			Type: "apikey",
			APIKey: []PostmanAuthAPIKey{
				{
					Key:   "key",
					Value: httpieAuth.Credentials.Username,
					Type:  "string",
				},
				{
					Key:   "value",
					Value: httpieAuth.Credentials.Password,
					Type:  "string",
				},
			},
		}
	default:
		// For unknown auth types, return nil and let headers handle it
		return nil
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

func extractVariablesFromWorkspace(httpie HTTPieWorkspace) []PostmanVariable {
	variableSet := make(map[string]string) // Use map to store variable names and their default values
	var variables []PostmanVariable

	// Build a map of all environment variables (prioritize default environment)
	envVarMap := make(map[string]string)
	var defaultEnv *HTTPieEnvironment
	for _, env := range httpie.Environments {
		if env.IsDefault {
			defaultEnv = &env
			break
		}
	}
	if defaultEnv == nil && len(httpie.Environments) > 0 {
		defaultEnv = &httpie.Environments[0]
	}
	if defaultEnv != nil {
		for _, envVar := range defaultEnv.Variables {
			envVarMap[envVar.Name] = envVar.Value
		}
	}
	// Add all other environments (do not overwrite default)
	for _, env := range httpie.Environments {
		if defaultEnv != nil && env.Name == defaultEnv.Name {
			continue
		}
		for _, envVar := range env.Variables {
			if _, exists := envVarMap[envVar.Name]; !exists {
				envVarMap[envVar.Name] = envVar.Value
			}
		}
	}

	// Add environment variables to variableSet
	for k, v := range envVarMap {
		variableSet[k] = v
	}

	// Extract variables from all URLs and headers (as before)
	variableRegex := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	// Process direct requests
	for _, req := range httpie.Entry.Requests {
		extractVariablesFromRequest(req, variableRegex, variableSet, envVarMap)
	}

	// Process collections
	for _, collection := range httpie.Entry.Collections {
		for _, req := range collection.Requests {
			extractVariablesFromRequest(req, variableRegex, variableSet, envVarMap)
		}
	}

	// Convert map to slice
	for varName, varValue := range variableSet {
		variables = append(variables, PostmanVariable{
			ID:    uuid.New().String(), // Generate a unique ID for each variable
			Key:   varName,
			Value: varValue,
			Type:  "string",
		})
	}

	return variables
}

// Updated signature to accept envVarMap
func extractVariablesFromRequest(req HTTPieRequest, variableRegex *regexp.Regexp, variableSet map[string]string, envVarMap map[string]string) {
	// Extract from URL
	matches := variableRegex.FindAllStringSubmatch(req.URL, -1)
	for _, match := range matches {
		if len(match) > 1 {
			varName := match[1]
			if val, exists := envVarMap[varName]; exists {
				variableSet[varName] = val
			} else if _, exists := variableSet[varName]; !exists {
				variableSet[varName] = "" // Empty default value
			}
		}
	}

	// Extract from headers
	for _, header := range req.Headers {
		matches := variableRegex.FindAllStringSubmatch(header.Value, -1)
		for _, match := range matches {
			if len(match) > 1 {
				varName := match[1]
				if val, exists := envVarMap[varName]; exists {
					variableSet[varName] = val
				} else if _, exists := variableSet[varName]; !exists {
					variableSet[varName] = ""
				}
			}
		}
	}

	// Extract from body
	if req.Body.Text.Value != "" {
		matches := variableRegex.FindAllStringSubmatch(req.Body.Text.Value, -1)
		for _, match := range matches {
			if len(match) > 1 {
				varName := match[1]
				if val, exists := envVarMap[varName]; exists {
					variableSet[varName] = val
				} else if _, exists := variableSet[varName]; !exists {
					variableSet[varName] = ""
				}
			}
		}
	}
}

func extractVariables(httpie HTTPieWorkspace) []PostmanVariable {
	// This function is kept for backward compatibility but now delegates to the new function
	return extractVariablesFromWorkspace(httpie)
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
