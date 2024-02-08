package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents from source to destination
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Flush any buffered data to ensure file is written completely
	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// CopyDir copies a directory recursively from src to dst.
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Construct the absolute path for the destination file or directory
		destPath := filepath.Join(dst, path[len(src):])

		if info.IsDir() {
			// If it's a directory, create it in the destination directory
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy the file
		return CopyFile(path, destPath)
	})
}

// CalculateFileInfo calculates the total size and number of files in the directory
func CalculateFileInfo(path string) (int, int64, error) {
	var totalSize int64
	var fileCount int // Counter for total number of files
	// Walk the directory tree recursively
	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %s\n", path, err)
			return nil
		}
		// Check if it's a regular file
		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				fmt.Printf("Error getting file info for %s: %s\n", path, err)
				return nil
			}

			totalSize += info.Size()
			fileCount++
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return fileCount, totalSize, nil
}
