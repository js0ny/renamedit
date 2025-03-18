package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 2 {
		fmt.Println("Usage: renamedit <directory path>")
		os.Exit(1)
	}

	dirPath := os.Args[1]

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
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			fileNames = append(fileNames, fileName)
			tmpFile.WriteString(fileName + "\n")
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
		newName := newFileNames[i]
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
