package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ignoreExt bool

func init() {
	const ignoreExtDesc = "Ignore file extensions when renaming"
	flag.BoolVar(&ignoreExt, "ignore-ext", false, ignoreExtDesc)
	flag.BoolVar(&ignoreExt, "i", false, ignoreExtDesc+" (shorthand)")
}

func main() {
	flag.Parse()

	// Check command line arguments
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: renamedit [options] <directory path>")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	dirPath := args[0]

	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Printf("Error: %s is not a directory\n", dirPath)
		os.Exit(1)
	}

	// Get list of files in the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "rename-*.txt")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())

	// Write filenames to temporary file
	fileNames := make([]string, 0)
	fileExts := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			fileNames = append(fileNames, fileName)

			// If ignoring extensions, split filename and extension
			if ignoreExt {
				baseName := fileName
				ext := ""

				// Get extension only if there is one
				if dotIndex := strings.LastIndex(fileName, "."); dotIndex > 0 {
					baseName = fileName[:dotIndex]
					ext = fileName[dotIndex:]
				}

				fileExts = append(fileExts, ext)
				tmpFile.WriteString(baseName + "\n")
			} else {
				tmpFile.WriteString(fileName + "\n")
				fileExts = append(fileExts, "") // Empty placeholder
			}
		}
	}
	tmpFile.Close()

	// Open temporary file with editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	editorFlags := map[string]string{
		"code":    "--wait",
		"subl":    "--wait",
		"zeditor": "--wait",
		"atom":    "--wait",
		"gedit":   "--standalone",
	}

	var cmd *exec.Cmd
	if flag, exists := editorFlags[editor]; exists {
		cmd = exec.Command(editor, flag, tmpFile.Name())
	} else {
		// Default case for vim, nano, emacs, etc. that wait by default
		cmd = exec.Command(editor, tmpFile.Name())
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Editor error: %v\n", err)
		os.Exit(1)
	}

	// Read edited file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Printf("Error reading temporary file: %v\n", err)
		os.Exit(1)
	}

	// Parse new filenames
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	var newFileNames []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			newFileNames = append(newFileNames, line)
		}
	}

	// Check if the number of files matches
	if len(newFileNames) != len(fileNames) {
		fmt.Println("Error: Number of new and old filenames doesn't match")
		os.Exit(1)
	}

	// Perform renaming
	for i, oldName := range fileNames {
		// Re-add extension if we were ignoring it
		newName := newFileNames[i]
		if ignoreExt {
			newName = newName + fileExts[i]
		}

		if oldName != newName {
			oldPath := filepath.Join(dirPath, oldName)
			newPath := filepath.Join(dirPath, newName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Printf("Failed to rename %s to %s: %v\n", oldName, newName, err)
			} else {
				fmt.Printf("Renamed: %s -> %s\n", oldName, newName)
			}
		}
	}
}
