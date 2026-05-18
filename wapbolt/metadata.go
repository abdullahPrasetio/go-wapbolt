package wapbolt

// RouteMetadata holds extra information about a specific route
type RouteMetadata struct {
	Description      string
	FieldValidations map[string]interface{}
	Examples         []Response
}

// Metadata Registry (Simple global map for the module)
var registry = make(map[string]RouteMetadata)

// RegisterMetadata associates metadata with a route key (Method + Path)
// Example: RegisterMetadata("POST", "/api/v1/login", RouteMetadata{...})
func RegisterMetadata(method, path string, meta RouteMetadata) {
	key := method + ":" + path
	registry[key] = meta
}

// GetMetadata retrieves metadata for a route
func GetMetadata(method, path string) (RouteMetadata, bool) {
	key := method + ":" + path
	meta, ok := registry[key]
	return meta, ok
}
