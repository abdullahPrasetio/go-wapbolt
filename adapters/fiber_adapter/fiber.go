package fiber_adapter

import (
	"fmt"
	"strings"

	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
	"github.com/gofiber/fiber/v2"
)

// ParseRoutes extracts all routes from a Fiber app and adds them to a Wapbolt collection
func ParseRoutes(app *fiber.App, collection *wapbolt.Collection, baseUrl string) {
	routes := app.GetRoutes()

	// Group routes by first path segment (basic heuristic for folders)
	folders := make(map[string]*wapbolt.Item)

	for _, route := range routes {
		// Skip internal fiber routes or middleware if any
		if route.Method == "USE" || route.Method == "HEAD" || route.Method == "OPTIONS" {
			continue
		}

		path := route.Path
		segments := strings.Split(strings.Trim(path, "/"), "/")
		
		folderName := "General"
		if len(segments) > 0 && segments[0] != "" {
			folderName = strings.Title(segments[0])
		}

		// Get or create folder
		folder, exists := folders[folderName]
		if !exists {
			folder = collection.AddFolder(folderName)
			folders[folderName] = folder
		}

		// Replace fiber params like :id with postman variables {{id}} or just keep as :id
		// Wapbolt usually accepts standard REST params
		normalizedPath := strings.ReplaceAll(path, ":", "") // simplistic param handling for demo

		// Create request item
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

		// Add body placeholder for methods that usually have body
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
