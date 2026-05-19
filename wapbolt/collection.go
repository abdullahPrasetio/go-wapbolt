package wapbolt

import (
	"encoding/json"
	"os"
)

// Collection represents the root of a Wapbolt/Postman compatible collection
type Collection struct {
	Info Info   `json:"info"`
	Item []Item `json:"item"`
}

// Info holds metadata about the collection
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema,omitempty"`
}

// Item can be either a folder (containing more items) or a request
type Item struct {
	Name    string   `json:"name"`
	Item    []Item   `json:"item,omitempty"` // For folders
	Request *Request `json:"request,omitempty"` // For actual requests
}

// Request holds the HTTP request details
type Request struct {
	Method           string                 `json:"method"`
	URL              interface{}            `json:"url"`
	Header           []HeaderField          `json:"header,omitempty"`
	Body             *Body                  `json:"body,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Responses        []Response             `json:"response,omitempty"`         // Examples in Postman/Wapbolt
	FieldValidations map[string]interface{} `json:"field_validations,omitempty"` // Wapbolt specific extension
}

// Response represents a request example/response
type Response struct {
	Name   string        `json:"name"`
	Status string        `json:"status,omitempty"`
	Code   int           `json:"code,omitempty"`
	Header []HeaderField `json:"header,omitempty"`
	Body   string        `json:"body,omitempty"`
}

// ValidationRule represents Wapbolt field validation rules
type ValidationRule struct {
	Rules       []string `json:"rules,omitempty"` // e.g., ["required", "email"]
	Min         float64  `json:"min,omitempty"`
	Max         float64  `json:"max,omitempty"`
	Nullable    bool     `json:"nullable"`
	Description string   `json:"description,omitempty"`
}

// HeaderField represents a single HTTP header
type HeaderField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Body represents the request payload
type Body struct {
	Mode string `json:"mode"` // e.g., "raw", "formdata"
	Raw  string `json:"raw,omitempty"`
}

// URL represents a structured URL (optional, often just a string is used)
type URL struct {
	Raw string `json:"raw"`
}

// NewCollection creates a new empty collection
func NewCollection(name, description string) *Collection {
	return &Collection{
		Info: Info{
			Name:        name,
			Description: description,
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: []Item{},
	}
}

// AddFolder adds a new folder to the root of the collection and returns a pointer to it
func (c *Collection) AddFolder(name string) *Item {
	folder := Item{
		Name: name,
		Item: []Item{},
	}
	c.Item = append(c.Item, folder)
	return &c.Item[len(c.Item)-1]
}

// AddRequest adds a request to the root of the collection
func (c *Collection) AddRequest(req Item) {
	c.Item = append(c.Item, req)
}

// AddRequestToFolder adds a request to a specific folder
func (folder *Item) AddRequest(req Item) {
	folder.Item = append(folder.Item, req)
}

// ExportToFile saves the collection to a JSON file
func (c *Collection) ExportToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// ToJSON returns the collection as formatted JSON bytes
func (c *Collection) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
