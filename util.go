package caddy_wasm

import (
	"os"
	"path/filepath"
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
