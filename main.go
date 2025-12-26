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
	top := flag.Int(`top`, 0, `number of to items to print out`)
	methodCount := flag.Int(`min`, 0, `min number of methods`)
	vendor := flag.Bool(`vendor`, false, `include vendor folder too`)
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	var interfaces []InterfaceInfo
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// skip vendor folder if flag is false
		if !*vendor && (strings.HasPrefix(path, "vendor/") || strings.Contains(path, "/vendor/")) {
			return nil
		}

		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		ast.Inspect(f, func(n ast.Node) bool {
			// Look for type declarations
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			// Check if the type is an interface
			if it, ok := ts.Type.(*ast.InterfaceType); ok {
				count := 0
				if it.Methods != nil {
					count = len(it.Methods.List)
				}

				pos := fset.Position(ts.Pos())
				interfaces = append(interfaces, InterfaceInfo{
					Name:        ts.Name.Name,
					FilePath:    path,
					LineNumber:  pos.Line,
					MethodCount: count,
				})
			}
			return true
		})

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	// Sort by MethodCount descending
	sort.Slice(interfaces, func(i, j int) bool {
		if interfaces[i].MethodCount == interfaces[j].MethodCount {
			return interfaces[i].Name < interfaces[j].Name
		}
		return interfaces[i].MethodCount > interfaces[j].MethodCount
	})

	fmt.Printf("%-5s | %-5s | %-30s | %s\n", "index", "Count", "Interface", "Location")
	fmt.Println(strings.Repeat("-", 80))

	for i, inf := range interfaces {
		if *top > 0 && i >= *top {
			break
		}
		if *methodCount > 0 && inf.MethodCount < *methodCount {
			break
		}

		fmt.Printf("%-5d | %-5d | %-30s | %s:%d\n", i+1, inf.MethodCount, inf.Name, inf.FilePath, inf.LineNumber)
	}
}
