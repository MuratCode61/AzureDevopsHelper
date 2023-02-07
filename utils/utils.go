package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ZipFolder(folderDir string, folderPath string) string {
	zipPath := fmt.Sprintf("%s.zip", folderPath)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Println("Error has been occurred while creating archieve file.", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// remove exportDir path prefix for not creating exportDir in zip file.
		path = strings.TrimPrefix(path, fmt.Sprintf("%s%c", folderDir, os.PathSeparator))
		f, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(folderPath, walker)
	if err != nil {
		log.Println("Error has been occurred while zipping workitemsexport folder", err)
	}

	defer zipWriter.Close()

	return zipPath
}
