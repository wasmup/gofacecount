package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type InterfaceInfo struct {
	Name        string
	FilePath    string
	LineNumber  int
	MethodCount int
}

func main() {
	top := flag.Int("top", 0, "number of top items to print out")
	minMethods := flag.Int("min", 0, "min number of methods")
	vendor := flag.Bool("vendor", false, "include vendor folder")
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	var interfaces []InterfaceInfo
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			// Skip vendor directories if not requested
			if !*vendor && info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Use parser.DeclarationErrors to get better feedback
		f, err := parser.ParseFile(fset, path, nil, parser.DeclarationErrors)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", path, err)
			return nil
		}

		ast.Inspect(f, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if it, ok := ts.Type.(*ast.InterfaceType); ok {
				count := 0
				if it.Methods != nil {
					for _, field := range it.Methods.List {
						// Each field can represent multiple names, though rare in interfaces
						if len(field.Names) > 0 {
							count += len(field.Names)
						} else {
							// It's an embedded interface
							count++
						}
					}
				}

				if count >= *minMethods {
					pos := fset.Position(ts.Pos())
					interfaces = append(interfaces, InterfaceInfo{
						Name:        ts.Name.Name,
						FilePath:    path,
						LineNumber:  pos.Line,
						MethodCount: count,
					})
				}
			}
			return true
		})
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking path: %v\n", err)
		os.Exit(1)
	}

	sort.Slice(interfaces, func(i, j int) bool {
		if interfaces[i].MethodCount == interfaces[j].MethodCount {
			return interfaces[i].Name < interfaces[j].Name
		}
		return interfaces[i].MethodCount > interfaces[j].MethodCount
	})

	// Print Header
	fmt.Printf("%-5s | %-5s | %-40s | %s\n", "Rank", "Meths", "Interface Name", "Location")
	fmt.Println(strings.Repeat("-", 90))

	for i, inf := range interfaces {
		if *top > 0 && i >= *top {
			break
		}
		fmt.Printf("%-5d | %-5d | %-40s | %s:%d\n", i+1, inf.MethodCount, inf.Name, inf.FilePath, inf.LineNumber)
	}
}
