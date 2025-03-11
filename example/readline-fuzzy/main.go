package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/petrjanda/readline"
)

func main() {

	fsys := os.DirFS("example")
	completionItems, err := listFiles(fsys, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	log.Println(completionItems)

	// Create a new fuzzy completer
	completer := readline.NewFuzzyCompleter(completionItems)

	// Create a new readline instance with the fuzzy completer
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "fruits> ",
		AutoComplete: completer,
		HistoryFile:  "/tmp/readline-fuzzy-demo",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Println("Fuzzy Completion Demo")
	fmt.Println("=====================")
	fmt.Println("Type part of a fruit name and press TAB for fuzzy completion")
	fmt.Println("Try typing 'ber' or 'an' and press TAB")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		fmt.Printf("You selected: %s\n", line)
	}
}

// listFiles recursively lists all files and directories in the given filesystem
func listFiles(fsys fs.FS, dir string) ([]string, error) {
	var paths []string

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryPath := entry.Name()
		if dir != "." {
			entryPath = filepath.Join(dir, entry.Name())
		}

		paths = append(paths, entryPath)

		if entry.IsDir() {
			subPaths, err := listFiles(fsys, entryPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		}
	}

	return paths, nil
}
