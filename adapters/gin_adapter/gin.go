package gin_adapter

import (
	"fmt"
	"strings"

	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
	"github.com/gin-gonic/gin"
)

// ParseRoutes extracts all routes from a Gin engine and adds them to a Wapbolt collection
func ParseRoutes(engine *gin.Engine, collection *wapbolt.Collection, baseUrl string) {
	routes := engine.Routes()

	folders := make(map[string]*wapbolt.Item)

	for _, route := range routes {
		if route.Method == "HEAD" || route.Method == "OPTIONS" {
			continue
		}

		path := route.Path
		segments := strings.Split(strings.Trim(path, "/"), "/")
		
		folderName := "General"
		if len(segments) > 0 && segments[0] != "" {
			folderName = strings.Title(segments[0])
		}

		folder, exists := folders[folderName]
		if !exists {
			folder = collection.AddFolder(folderName)
			folders[folderName] = folder
		}

		// Normalize gin params (e.g. /users/:id to /users/{{id}})
		// Here we just use a simplistic approach for demonstration
		normalizedPath := strings.ReplaceAll(path, ":", "")

		reqItem := wapbolt.Item{
			Name: fmt.Sprintf("%s", path),
			Request: &wapbolt.Request{
				Method: route.Method,
				URL:    fmt.Sprintf("%s%s", baseUrl, normalizedPath),
				Header: []wapbolt.HeaderField{
					{Key: "Accept", Value: "application/json"},
				},
			},
		}

		// Check for registered metadata
		if meta, ok := wapbolt.GetMetadata(route.Method, path); ok {
			reqItem.Request.Description = meta.Description
			reqItem.Request.FieldValidations = meta.FieldValidations
			reqItem.Request.Responses = meta.Examples
		}

		if route.Method == "POST" || route.Method == "PUT" || route.Method == "PATCH" {
			reqItem.Request.Body = &wapbolt.Body{
				Mode: "raw",
				Raw:  "{\n  \n}",
			}
			reqItem.Request.Header = append(reqItem.Request.Header, wapbolt.HeaderField{
				Key:   "Content-Type",
				Value: "application/json",
			})
		}

		folder.AddRequest(reqItem)
	}
}
