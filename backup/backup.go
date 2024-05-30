package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var excludedExtensions = map[string]struct{}{
	".zip": {},
}

func Backup(w io.WriteCloser, dir string) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if _, ok := excludedExtensions[filepath.Ext(path)]; ok {
			fmt.Printf("skipping %s...\n", path)
			return nil
		}

		file, err := entry.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(file)
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		fileReader, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(writer, fileReader)
		return err
	})

	return err
}
