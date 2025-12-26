This Go tool is a static analysis utility that scans your source code for interfaces and reports their complexity based on method count. It is particularly useful for identifying big interfaces that might violate the Interface Segregation Principle.

---

# Go Interface Scanner

A command-line tool that parses Go source files to locate interface definitions, counts their methods, and displays them in a sorted list.

## Features

* **Recursive Scanning**: Walks through directories to find all `.go` files.
* **Static Analysis**: Uses the standard `go/ast` and `go/parser` libraries (no reflection required).
* **Filtering**: Filter results by minimum method count or limit to a "top N" list.
* **Smart Defaults**: Automatically ignores the `vendor/` directory unless specified otherwise.
* **Detailed Output**: Shows the interface name, method count, and the exact file/line number.

---

# Clone the repository

```bash
git clone https://github.com/wasmup/gofacecount.git
cd gofacecount
```

## Installation

```bash
go install github.com/wasmup/gofacecount@latest

# or
CGO_ENABLED=0 go install -x -ldflags=-s 
go version -m  $(which gofacecount)
```

# Build the binary

```bash
go build 
```

## Usage

Run the tool against a directory (defaults to current directory if none provided):

```bash

gofacecount [flags] [path]
gofacecount
gofacecount -top 10
gofacecount -top 10 -min 5
gofacecount -vendor -top 10
```

### Flags

| Flag | Type | Description |
| --- | --- | --- |
| `-top` | `int` | Only print the top X results (sorted by highest method count). |
| `-min` | `int` | Only print interfaces with at least X methods. |
| `-vendor` | `bool` | Include the `vendor/` folder in the scan (default: false). |

### Examples

**Find all interfaces in the current directory:**

```bash
gofacecount .

```

**Find the top 5 largest interfaces:**

```bash
gofacecount -top 5

```

**Find interfaces with 10 or more methods:**

```bash
gofacecount -min 10

```

---

## Sample Output

```text
index | Count | Interface                      | Location
--------------------------------------------------------------------------------
1     | 37    | Type                           | reflect/type.go:40
2     | 23    | TB                             | testing/testing.go:881
3     | 22    | Node                           | cmd/compile/internal/ir/node.go:19

```

## How it Works

The tool utilizes the Abstract Syntax Tree (AST) to traverse the code structure:

1. **Walk**: Recursively visits files in the provided path.
2. **Parse**: Converts Go source into an AST via `go/parser`.
3. **Inspect**: Targets `ast.TypeSpec` nodes specifically looking for `*ast.InterfaceType`.
4. **Sort**: Orders results primarily by method count (descending) and secondarily by name.
