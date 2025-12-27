package utils

import (
	"archive/zip"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func Join(base string, elems ...string) string {
	return filepath.Join(append([]string{base}, elems...)...)
}

func SaveUploadedTemp(baseTmp string, fh *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(baseTmp, 0755); err != nil {
		return "", err
	}
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	p := filepath.Join(baseTmp, fh.Filename)
	out, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, f); err != nil {
		return "", err
	}
	return p, nil
}

func SaveBytesTemp(baseTmp, filename string, data []byte) (string, error) {
	if err := os.MkdirAll(baseTmp, 0755); err != nil {
		return "", err
	}
	p := filepath.Join(baseTmp, filename)
	if err := os.WriteFile(p, data, 0644); err != nil {
		return "", err
	}
	return p, nil
}

func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func RemoveQuiet(path string) { _ = os.Remove(path) }

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func UnzipPDFs(dstDir, zipPath string) ([]string, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, err
	}
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var pdfs []string
	for _, f := range r.File {
		name := f.Name
		if strings.HasPrefix(name, "__MACOSX") {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		p := filepath.Join(dstDir, filepath.Base(name))
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			rc.Close()
			return nil, err
		}
		out, err := os.Create(p)
		if err != nil {
			rc.Close()
			return nil, err
		}
		if _, err := io.Copy(out, rc); err != nil {
			rc.Close()
			out.Close()
			return nil, err
		}
		rc.Close()
		out.Close()
		pdfs = append(pdfs, p)
	}
	if len(pdfs) == 0 {
		return nil, errors.New("no pdf found in zip")
	}
	return pdfs, nil
}
