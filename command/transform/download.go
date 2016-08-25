package transform

import (
	"bytes"
	"fmt"
	"github.com/tueftler/doget/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Origin struct {
	Host    string
	Vendor  string
	Name    string
	Version string
	Dir     string
}

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

func origin(reference string) Origin {
	var parsed []string
	var version, dir string

	pos := strings.LastIndex(reference, ":")
	if pos == -1 {
		parsed = strings.Split(reference, "/")
		version = "master"
	} else {
		parsed = strings.Split(reference[0:pos], "/")
		version = reference[pos+1 : len(reference)]
	}

	if len(parsed) == 3 {
		dir = ""
	} else {
		dir = strings.Join(parsed[3:len(parsed)], "/")
	}

	return Origin{Host: parsed[0], Vendor: parsed[1], Name: parsed[2], Dir: dir, Version: version}
}

func fetch(reference string, config *config.Configuration, progress func(transferred, total int64)) (string, error) {
	origin := origin(reference)

	if repository, ok := config.Repositories[origin.Host]; ok {
		template, err := template.New(origin.Host).Parse(repository["url"])
		if err != nil {
			return "", err
		}

		var uri bytes.Buffer
		if err := template.Execute(&uri, origin); err != nil {
			return "", err
		}

		target := filepath.Join("vendor", origin.Host, origin.Vendor, origin.Name)
		zip := filepath.Join(target, origin.Version+".zip")

		if err := os.MkdirAll(target, 0755); err != nil {
			return "", err
		}

		_, err = download(uri.String(), zip, progress)
		if err != nil {
			return "", err
		}

		if err := unzip(zip, target, strings.NewReplacer(origin.Name+"-"+origin.Version+"/", "")); err != nil {
			return "", err
		}

		return filepath.Join(target, origin.Dir), nil
	} else {
		return "", fmt.Errorf("No repository %s", origin.Host)
	}
}
