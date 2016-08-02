package main

import (
    "bufio"
    "io"
    "fmt"
    "os"
    "strings"
    "regexp"
    "net/http"
    "archive/zip"
    "path/filepath"
)

func Unzip(src, dest string, base *strings.Replacer) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extract := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(dest, base.Replace(f.Name))

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extract(f)
        if err != nil {
            return err
        }
    }

    return nil
}

func Download(uri, file string) (int64, error) {
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

func Transform(input string, base string, yield func(string)) (error) {
    transformations := map[string]func([]string) error{
        "FROM (.+)"    : func(matches []string) error {
            if "" == base {
                base = matches[1]
                yield("FROM " + base)
                return nil
            } else if (matches[1] == base) {
                return nil
            }
            return fmt.Errorf("Expecting %q, have %q in %q", base, matches[1], input)
        },
        "INCLUDE (.+)" : func(matches []string) error {
            yield("# Included from " + matches[1])

            include := regexp.MustCompile("github.com/([^/]+)/([^/]+)").FindStringSubmatch(matches[1])
            if len(include) > 0 {
                vendor := include[1]
                name := include[2]

                target := filepath.Join("vendor", vendor, name)
                zip := filepath.Join(target, "master.zip")

                if err := os.MkdirAll(target, 0700); err != nil {
                    return err
                }

                _, err := Download("https://github.com/" + vendor + "/" + name + "/archive/master.zip", zip)
                if err != nil {
                    return err
                }

                if err := Unzip(zip, target, strings.NewReplacer(name + "-master/", "")); err != nil {
                    return err
                }

                Transform(filepath.Join(target, "Dockerfile.in"), base, yield)
            }
            return nil
        },
    }

    file, err := os.Open(input)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        transformed := false

        for pattern, action := range transformations {
            matches := regexp.MustCompile(pattern).FindStringSubmatch(line)
            if len(matches) > 0 {
                if err := action(matches); err != nil {
                    return err
                }
                transformed = true
                break
            }
        }
        
        if !transformed {
            yield(line)
        }
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}

func main() {
    err := Transform("Dockerfile.in", "", func(line string) { fmt.Println(line) })

    if (err != nil) {
        fmt.Println(err)
        os.Exit(1)
    }
}