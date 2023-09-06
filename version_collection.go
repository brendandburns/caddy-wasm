package caddy_wasm

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemInterface interface {
	Glob(p string) ([]string, error)
	ReadFile(f string) ([]byte, error)
}

type osfs struct{}

func (o osfs) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (o osfs) ReadFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

type VersionCollection interface {
	GetVersions() ([]string, error)
	GetWebAssembly(version string) ([]byte, error)
}

type fileSystemVersionCollection struct {
	fs   FileSystemInterface
	dir  string
	glob string
}

func ForFilesystemGlob(dir, glob string) VersionCollection {
	return &fileSystemVersionCollection{osfs{}, dir, glob}
}

func (f *fileSystemVersionCollection) GetVersions() ([]string, error) {
	files, err := f.fs.Glob(filepath.Join(f.dir, f.glob))
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		files[i] = filepath.Base(file)
	}
	return files, nil
}

func (f *fileSystemVersionCollection) GetWebAssembly(version string) ([]byte, error) {
	return f.fs.ReadFile(filepath.Join(f.dir, version))
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
