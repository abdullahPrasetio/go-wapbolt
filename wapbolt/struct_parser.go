package wapbolt

import (
	"reflect"
	"strconv"
	"strings"
)

// ParseStruct uses reflection to extract validation rules from struct tags
// It supports tags: json (for name), validate (for rules), and description (for info)
func ParseStruct(s interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	v := reflect.TypeOf(s)

	// If pointer, get the underlying element
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		
		// Get field name from JSON tag, fallback to field name
		fieldName := field.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}
		// Clean up json tag options like "email,omitempty"
		fieldName = strings.Split(fieldName, ",")[0]

		rule := ValidationRule{
			Description: field.Tag.Get("description"),
		}

		// Parse validate tag (e.g., "required,email,min=5,max=20")
		validateTag := field.Tag.Get("validate")
		if validateTag != "" {
			tags := strings.Split(validateTag, ",")
			for _, t := range tags {
				t = strings.TrimSpace(t)
				if strings.HasPrefix(t, "min=") {
					val, _ := strconv.ParseFloat(strings.TrimPrefix(t, "min="), 64)
					rule.Min = val
				} else if strings.HasPrefix(t, "max=") {
					val, _ := strconv.ParseFloat(strings.TrimPrefix(t, "max="), 64)
					rule.Max = val
				} else {
					rule.Rules = append(rule.Rules, t)
				}
			}
		}

		// Handle nested structs recursively
		if field.Type.Kind() == reflect.Struct {
			// You could recursively call ParseStruct here if needed for nested objects
		}

		fields[fieldName] = rule
	}

	return fields
}
