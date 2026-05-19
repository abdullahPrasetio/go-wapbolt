package wapbolt

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Scanner handles static analysis of Go source code
type Scanner struct {
	fset     *token.FileSet
	structs  map[string]map[string]interface{} // map[structName]FieldValidations
	handlers map[string]ast.Node               // map[handlerName]bodyNode
}

func NewScanner() *Scanner {
	return &Scanner{
		fset:     token.NewFileSet(),
		structs:  make(map[string]map[string]interface{}),
		handlers: make(map[string]ast.Node),
	}
}

func (s *Scanner) ScanDir(dirPath string) (map[string]RouteMetadata, error) {
	metadataMap := make(map[string]RouteMetadata)

	// First pass: Find all definitions
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "/vendor/") || strings.Contains(path, ".git") {
			return nil
		}
		s.collectDefinitions(path)
		return nil
	})

	// Second pass: Scan routes
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "/vendor/") || strings.Contains(path, ".git") {
			return nil
		}
		s.scanRoutes(path, metadataMap)
		return nil
	})

	return metadataMap, nil
}

func (s *Scanner) collectDefinitions(filePath string) {
	f, err := parser.ParseFile(s.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				s.structs[ts.Name.Name] = s.parseStructNode(st)
			}
		}
		if fn, ok := n.(*ast.FuncDecl); ok {
			s.handlers[fn.Name.Name] = fn.Body
		}
		return true
	})
}

func (s *Scanner) parseStructNode(st *ast.StructType) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		fieldName := field.Names[0].Name
		rule := ValidationRule{}

		if field.Tag != nil {
			tagVal := strings.Trim(field.Tag.Value, "`")
			
			// Extract JSON name and check for omitempty
			if jsonTag := getTagValue(tagVal, "json"); jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				fieldName = parts[0]
				for _, p := range parts {
					if p == "omitempty" {
						rule.Nullable = true
					}
				}
			}

			rule.Description = getTagValue(tagVal, "description")
			if valTag := getTagValue(tagVal, "validate"); valTag != "" {
				parts := strings.Split(valTag, ",")
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if strings.HasPrefix(p, "min=") {
						rule.Min, _ = strconv.ParseFloat(strings.TrimPrefix(p, "min="), 64)
					} else if strings.HasPrefix(p, "max=") {
						rule.Max, _ = strconv.ParseFloat(strings.TrimPrefix(p, "max="), 64)
					} else {
						rule.Rules = append(rule.Rules, p)
					}
				}
			}
		}
		fields[fieldName] = rule
	}
	return fields
}

func (s *Scanner) scanRoutes(filePath string, results map[string]RouteMetadata) {
	f, err := parser.ParseFile(s.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return
	}

	groupPrefixes := make(map[string]string)

	ast.Inspect(f, func(n ast.Node) bool {
		if assign, ok := n.(*ast.AssignStmt); ok {
			for i, rhs := range assign.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Group" {
						if len(call.Args) >= 1 {
							if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
								prefix := strings.Trim(lit.Value, "\"")
								if i < len(assign.Lhs) {
									if lhsIdent, ok := assign.Lhs[i].(*ast.Ident); ok {
										parentPrefix := ""
										if callerIdent, ok := sel.X.(*ast.Ident); ok {
											parentPrefix = groupPrefixes[callerIdent.Name]
										}
										groupPrefixes[lhsIdent.Name] = parentPrefix + "/" + strings.Trim(prefix, "/")
									}
								}
							}
						}
					}
				}
			}
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		method := strings.ToUpper(sel.Sel.Name)
		if isHTTPMethod(method) && len(call.Args) >= 1 {
			if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
				path := strings.Trim(lit.Value, "\"")
				prefix := ""
				if callerIdent, ok := sel.X.(*ast.Ident); ok {
					prefix = groupPrefixes[callerIdent.Name]
				}
				
				fullPath := "/" + strings.Trim(prefix, "/") + "/" + strings.Trim(path, "/")
				fullPath = filepath.Clean(fullPath)
				
				key := method + ":" + fullPath
				meta := results[key]

				for _, arg := range call.Args {
					if fields := s.findBodyHint(arg); fields != nil {
						meta.FieldValidations = fields
						break
					}
					var handlerName string
					if ident, ok := arg.(*ast.Ident); ok {
						handlerName = ident.Name
					} else if sel, ok := arg.(*ast.SelectorExpr); ok {
						handlerName = sel.Sel.Name
					}

					if handlerName != "" {
						if body, ok := s.handlers[handlerName]; ok {
							if fields := s.findBodyHint(body); fields != nil {
								meta.FieldValidations = fields
								break
							}
						}
					}
				}
				results[key] = meta
			}
		}
		return true
	})
}

func (s *Scanner) findBodyHint(n ast.Node) map[string]interface{} {
	var foundFields map[string]interface{}
	localVarTypes := make(map[string]string)

	ast.Inspect(n, func(inner ast.Node) bool {
		getTypeName := func(expr ast.Expr) string {
			if ident, ok := expr.(*ast.Ident); ok {
				return ident.Name
			}
			if sel, ok := expr.(*ast.SelectorExpr); ok {
				return sel.Sel.Name
			}
			if star, ok := expr.(*ast.StarExpr); ok { // Handle *Struct
				if ident, ok := star.X.(*ast.Ident); ok {
					return ident.Name
				}
				if sel, ok := star.X.(*ast.SelectorExpr); ok {
					return sel.Sel.Name
				}
			}
			return ""
		}

		// Track: var req models.UserRequest
		if vspec, ok := inner.(*ast.ValueSpec); ok {
			typeName := getTypeName(vspec.Type)
			if typeName != "" {
				for _, name := range vspec.Names {
					localVarTypes[name.Name] = typeName
				}
			}
		}
		
		// Track: req := models.UserRequest{}
		if assign, ok := inner.(*ast.AssignStmt); ok {
			for i, rhs := range assign.Rhs {
				typeName := ""
				if comp, ok := rhs.(*ast.CompositeLit); ok {
					typeName = getTypeName(comp.Type)
				} else if call, ok := rhs.(*ast.CallExpr); ok {
					// Handle cases like new(UserRequest)
					if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "new" {
						if len(call.Args) > 0 {
							typeName = getTypeName(call.Args[0])
						}
					}
				}

				if typeName != "" && i < len(assign.Lhs) {
					if lhsIdent, ok := assign.Lhs[i].(*ast.Ident); ok {
						localVarTypes[lhsIdent.Name] = typeName
					}
				}
			}
		}

		call, ok := inner.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		name := sel.Sel.Name
		if name == "BodyParser" || name == "BindJSON" || name == "Bind" || name == "ShouldBindJSON" {
			if len(call.Args) > 0 {
				arg := call.Args[0]
				typeName := ""

				// Case: .BodyParser(&UserRequest{}) or .BodyParser(&req)
				if unary, ok := arg.(*ast.UnaryExpr); ok && unary.Op == token.AND {
					if comp, ok := unary.X.(*ast.CompositeLit); ok {
						typeName = getTypeName(comp.Type)
					} else if ident, ok := unary.X.(*ast.Ident); ok {
						typeName = localVarTypes[ident.Name]
					}
				} else if ident, ok := arg.(*ast.Ident); ok {
					typeName = localVarTypes[ident.Name]
				}

				if typeName != "" {
					// ALWAYS match against short name
					parts := strings.Split(typeName, ".")
					shortName := parts[len(parts)-1]
					if fields, exists := s.structs[shortName]; exists {
						foundFields = map[string]interface{}{"body": fields}
					}
				}
			}
		}
		return true
	})
	return foundFields
}

func getTagValue(tag, key string) string {
	search := key + ":\""
	idx := strings.Index(tag, search)
	if idx == -1 {
		return ""
	}
	start := idx + len(search)
	val := tag[start:]
	end := strings.Index(val, "\"")
	if end == -1 {
		return ""
	}
	return val[:end]
}

func isHTTPMethod(m string) bool {
	switch strings.ToUpper(m) {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return true
	}
	return false
}

func AutoScan(srcPath string) {
	fmt.Printf("[Wapbolt Scanner] Starting AutoScan in %s...\n", srcPath)
	s := NewScanner()
	found, err := s.ScanDir(srcPath)
	if err != nil {
		fmt.Printf("[Wapbolt Scanner] Error: %v\n", err)
		return
	}
	
	count := 0
	for key, meta := range found {
		parts := strings.Split(key, ":")
		if len(parts) == 2 {
			method := parts[0]
			path := "/" + strings.Trim(parts[1], "/")
			RegisterMetadata(method, path, meta)
			count++
		}
	}
	fmt.Printf("[Wapbolt Scanner] Scan complete. Registered %d routes to metadata registry.\n", count)
}

func GenerateSampleBody(fields map[string]interface{}) string {
	sample := make(map[string]interface{})
	for k, v := range fields {
		if rule, ok := v.(ValidationRule); ok {
			var val interface{} = "string"
			isEmail := false
			for _, r := range rule.Rules {
				if r == "email" {
					isEmail = true
				}
			}
			if isEmail {
				val = "example@mail.com"
			} else if rule.Min > 0 {
				val = "value_" + strconv.Itoa(int(rule.Min))
			}
			sample[k] = val
		}
	}
	b, _ := json.MarshalIndent(sample, "", "  ")
	return string(b)
}
