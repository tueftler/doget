package transform

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func mkzip(src, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	w := zip.NewWriter(f)
	defer w.Close()

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			_, err := w.CreateHeader(&zip.FileHeader{Name: filepath.ToSlash(path), ExternalAttrs: 0x10})
			return err
		} else {
			o, err := os.Open(path)
			if err != nil {
				return err
			}
			defer o.Close()

			h, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// modify to provide full path name
			h.Name = filepath.ToSlash(path)
			z, err := w.CreateHeader(h)
			if err != nil {
				return err
			}

			_, err = io.Copy(z, o)
			return err
		}
	})

	return nil
}
