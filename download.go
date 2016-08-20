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

func download(uri, file string) (int64, error) {
	out, err := os.Create(file)
	if err != nil {
		return -1, err
	}
	defer out.Close()

	resp, err := http.Get(uri)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return -1, err
	}

	return size, nil
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

func fetch(reference string) (string, error) {
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

		// TODO: If-Modified-Since!
		_, err := download(uri.String(), zip)
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
