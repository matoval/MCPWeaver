package parser

import "github.com/getkin/kin-openapi/openapi3"

// ParsedAPI represents a parsed OpenAPI specification with extracted information
type ParsedAPI struct {
	Document    *openapi3.T
	Title       string
	Version     string
	Description string
	BaseURL     string
	Servers     []string
	Operations  []Operation
	Schemas     map[string]*openapi3.SchemaRef
}

// Operation represents an API operation (endpoint + method)
type Operation struct {
	ID          string
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[string]*Response
	Security    []map[string][]string
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string
	In          string // path, query, header, cookie
	Description string
	Required    bool
	Schema      *openapi3.SchemaRef
	Example     interface{}
}

// RequestBody represents a request body
type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]*MediaType
}

// MediaType represents a media type in request/response
type MediaType struct {
	Schema   *openapi3.SchemaRef
	Example  interface{}
	Examples map[string]*openapi3.ExampleRef
}

// Response represents an API response
type Response struct {
	Description string
	Headers     map[string]*openapi3.HeaderRef
	Content     map[string]*MediaType
}