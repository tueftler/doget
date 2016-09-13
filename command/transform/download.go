package transform

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/use"
)

type Track struct {
	io.Reader
	total    int64
	length   int64
	progress func(transferred, total int64)
}

func (t *Track) Read(p []byte) (int, error) {
	n, err := t.Reader.Read(p)
	if n > 0 {
		t.total += int64(n)
		t.progress(t.total, t.length)
	}
	return n, err
}

func download(uri, file string, progress func(transferred, total int64)) (int64, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return -1, err
	}

	stat, err := os.Stat(file)
	if err == nil {
		req.Header.Add("If-Modified-Since", stat.ModTime().UTC().Format(http.TimeFormat))
	}

	// DEBUG fmt.Printf(">>> %+v\n", req)

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	// DEBUG fmt.Printf("<<< %+v\n", resp)

	switch resp.StatusCode {
	case 200:
		out, err := os.Create(file)
		if err != nil {
			return -1, err
		}
		defer out.Close()

		size, err := io.Copy(out, &Track{resp.Body, 0, resp.ContentLength, progress})
		if err != nil {
			return -1, err
		}

		return size, nil

	case 304:
		progress(stat.Size(), stat.Size())
		return stat.Size(), nil

	default:
		return -1, fmt.Errorf("Could not download %q, response %+v", uri, resp)
	}
}

func fetch(origin *use.Origin, useCache bool, progress func(transferred, total int64)) (string, error) {
	target := filepath.Join(config.Vendordir, origin.Host, origin.Vendor, origin.Name)
	zip := filepath.Join(target, origin.Version+".zip")

	doDownload := !useCache
	if _, err := os.Stat(target); err != nil {
		doDownload = true
	}

	if doDownload {
		if err := os.MkdirAll(target, 0755); err != nil {
			return "", err
		}

		if _, err := download(origin.Uri, zip, progress); err != nil {
			return "", err
		}

		if err := unzip(zip, target, strings.NewReplacer(origin.Name+"-"+origin.Version+"/", "")); err != nil {
			return "", err
		}
	} else {
		fmt.Printf("> Using %s", origin.String())
	}

	return filepath.Join(target, origin.Dir), nil
}
