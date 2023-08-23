package caddy_wasm

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type VersionCollection interface {
	GetVersions() ([]string, error)
	GetWebAssembly(version string) ([]byte, error)
}

type fileSystemVersionCollection struct {
	dir  string
	glob string
}

func ForFilesystemGlob(dir, glob string) VersionCollection {
	return &fileSystemVersionCollection{dir, glob}
}

func (f *fileSystemVersionCollection) GetVersions() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(f.dir, f.glob))
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		files[i] = filepath.Base(file)
	}
	return files, nil
}

func (f *fileSystemVersionCollection) GetWebAssembly(version string) ([]byte, error) {
	return os.ReadFile(filepath.Join(f.dir, version))
}

type urlVersionCollection struct {
	url string
}

func ForURL(url string) VersionCollection {
	return &urlVersionCollection{url}
}

func (u *urlVersionCollection) GetVersions() ([]string, error) {
	url, err := url.Parse(u.url)
	if err != nil {
		return nil, err
	}
	ix := strings.LastIndex(url.Path, "/")
	return []string{url.Path[ix+1:]}, nil
}

func (u *urlVersionCollection) GetWebAssembly(version string) ([]byte, error) {
	res, err := http.Get(u.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
