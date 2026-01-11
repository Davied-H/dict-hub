# MDX/MDD Parser for Go
[CN](README_ZH.md)
This is a high-performance MDict (.mdx/.mdd) file parsing library written in Go. It supports querying dictionary content, retrieving dictionary metadata, and provides a file system wrapper compliant with the `io/fs.FS` interface, making it easy to integrate with other Go ecosystem libraries (such as HTTP servers).

This library was originally based on the [terasum/medict](https://github.com/terasum/medict) project and has undergone extensive bug fixes, performance optimizations, and code refactoring.

## Features

- **High-Performance Queries**: Uses a binary search algorithm for fast O(log n) complexity queries of keywords.
- **MDX/MDD Support**: Supports both .mdx (text dictionaries) and .mdd (resource files) formats.
- **Standard Interface**: Implements the `io/fs.FS` interface, allowing dictionaries to be easily served as a file system.
- **Robust Error Handling**: Comprehensive error handling and logging.
- **Complete Metadata**: Provides an API to access all dictionary metadata (such as title, description, creation date, etc.).

## Installation

```bash
go get github.com/lib-x/mdx
```
*(Note: Please replace `github.com/lib-x/mdx` with the actual repository path)*

## Usage Examples

### Example 1: Querying an MDX Dictionary

Here is a simple example of how to load an MDX dictionary and query a word.

```go
package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mdx" // Assuming this is the module path
)

func main() {
	// 1. Create a new Mdict instance
	// Replace "path/to/your/dictionary.mdx" with your MDX file path
	mdict, err := mdx.New("path/to/your/dictionary.mdx")
	if err != nil {
		log.Fatalf("Failed to load dictionary file: %v", err)
	}

	// 2. Build the index (recommended to be done once at program startup)
	err = mdict.BuildIndex()
	if err != nil {
		log.Fatalf("Failed to build dictionary index: %v", err)
	}

	// 3. Print dictionary information
	fmt.Printf("Dictionary Title: %s\n", mdict.Title())
	fmt.Printf("Dictionary Description: %s\n", mdict.Description())

	// 4. Query a word
	word := "hello"
	definition, err := mdict.Lookup(word)
	if err != nil {
		log.Fatalf("Failed to look up word '%s': %v", word, err)
	}

	fmt.Printf("Definition of '%s':\n%s\n", word, string(definition))

	// 5. Query a non-existent word
	word = "nonexistentword"
	_, err = mdict.Lookup(word)
	if err != nil {
		fmt.Printf("As expected, an error occurred when querying a non-existent word '%s': %v\n", word, err)
	}
}
```

### Example 2: Listing MDD Resource Files

MDD files typically contain resources like audio and images. The following example shows how to list all resources in an MDD file.

```go
package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mdx" // Assuming this is the module path
)

func main() {
	// 1. Load the MDD file
	// Replace "path/to/your/resource.mdd" with your MDD file path
	mdd, err := mdx.New("path/to/your/resource.mdd")
	if err != nil {
		log.Fatalf("Failed to load MDD file: %v", err)
	}

	// 2. Build the index
	err = mdd.BuildIndex()
	if err != nil {
		log.Fatalf("Failed to build MDD index: %v", err)
	}

	// 3. Get and print all keyword entries (in MDD, this is usually the file path)
	entries, err := mdd.GetKeyWordEntries()
	if err != nil {
		log.Fatalf("Failed to get keyword entries: %v", err)
	}

	fmt.Printf("Found %d resource files in '%s':\n", len(entries), mdd.Name())
	for i, entry := range entries {
		// Print only the first 10 as an example
		if i >= 10 {
			break
		}
		fmt.Println(entry.KeyWord) // The KeyWord field stores the resource file path
	}
}
```

## Contributing

Issues and Pull Requests are welcome.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).