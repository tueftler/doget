package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	repositories = map[string]*template.Template{
		"github.com": template.Must(template.New("github.com").Parse("https://github.com/{{.Vendor}}/{{.Name}}/archive/{{.Version}}.zip")),
	}
)

type Origin struct {
	Host    string
	Vendor  string
	Name    string
	Version string
}

type track struct {
	io.Reader
	total    int64
	length   int64
	progress func(transferred, total int64)
}

func (t *track) Read(p []byte) (int, error) {
	n, err := t.Reader.Read(p)
	t.total += int64(n)
	t.progress(t.total, t.length)
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

		size, err := io.Copy(out, &track{resp.Body, 0, resp.ContentLength, progress})
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
	pos := strings.LastIndex(reference, ":")
	if pos == -1 {
		parsed := strings.SplitN(reference, "/", 3)
		return Origin{Host: parsed[0], Vendor: parsed[1], Name: parsed[2], Version: "master"}
	} else {
		parsed := strings.SplitN(reference[0:pos], "/", 3)
		return Origin{Host: parsed[0], Vendor: parsed[1], Name: parsed[2], Version: reference[pos+1 : len(reference)]}
	}
}

func fetch(reference string, progress func(transferred, total int64)) (string, error) {
	origin := origin(reference)

	if delegate, ok := repositories[origin.Host]; ok {
		var uri bytes.Buffer
		if err := delegate.Execute(&uri, origin); err != nil {
			return "", err
		}

		target := filepath.Join("vendor", origin.Host, origin.Vendor, origin.Name)
		zip := filepath.Join(target, origin.Version+".zip")

		if err := os.MkdirAll(target, 0700); err != nil {
			return "", err
		}

		_, err := download(uri.String(), zip, progress)
		if err != nil {
			return "", err
		}

		if err := unzip(zip, target, strings.NewReplacer(origin.Name+"-"+origin.Version+"/", "")); err != nil {
			return "", err
		}

		return target, nil
	} else {
		return "", fmt.Errorf("No repository %s", origin.Host)
	}
}
