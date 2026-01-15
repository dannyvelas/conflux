package main

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

var _ Reader = (*fileReader)(nil)

type fileReader struct {
	fileSystem fs.FS
	paths      []string
}

// NewFileReader creates a new reader which gets key-value pairs from YML files from specified directories/files
func NewFileReader(path string, opts ...func(*fileReader)) *fileReader {
	fileReader := fileReader{
		fileSystem: os.DirFS("."),
		paths:      []string{path},
	}

	for _, opt := range opts {
		opt(&fileReader)
	}

	return &fileReader
}

func (r *fileReader) Read() (ReadResult, error) {
	configMap, diagnostics := make(map[string]string), make(map[string]string)
	for _, path := range r.paths {
		info, err := fs.Stat(r.fileSystem, path)
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics[path] = "Skipped: Not Found"
			continue
		} else if err != nil {
			return nil, fmt.Errorf("error getting info for path (%s): %v", path, err)
		}

		pathConfig, err := r.getConfigByInfo(info, path)
		if err != nil {
			return nil, fmt.Errorf("error reading config from path (%s): %v", path, err)
		}

		maps.Copy(configMap, pathConfig)
	}
	return NewSimpleReadResult(configMap), nil
}

func (r *fileReader) getConfigByInfo(info fs.FileInfo, path string) (map[string]string, error) {
	mode := info.Mode()
	switch {
	case mode.IsDir():
		return r.readDirectory(path)
	case mode.IsRegular():
		return r.readFile(path)
	default:
		return nil, fmt.Errorf("path %s is a special file (%v) and cannot be read as config", path, mode)
	}
}

func (r *fileReader) readDirectory(dir string) (map[string]string, error) {
	dirConfig := make(map[string]string)

	err := fs.WalkDir(r.fileSystem, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Support both .yml and .yaml
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}

		fileConfig, err := r.readFile(path)
		if err != nil {
			return err
		}

		maps.Copy(dirConfig, fileConfig)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking config directory (%s): %v", dir, err)
	}

	return dirConfig, nil
}

func (r *fileReader) readFile(file string) (map[string]string, error) {
	fileConfig := make(map[string]string)
	data, err := fs.ReadFile(r.fileSystem, file)
	if err != nil {
		return nil, fmt.Errorf("error reading config file(%s): %v", file, err)
	}
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file (%s): %v", file, err)
	}
	return fileConfig, nil
}

// WithFileSystem allows specifying a custom file system for the fileReader
// By default, it uses the OS file system, os.DirFS(".")
func WithFileSystem(fileSystem fs.FS) func(*fileReader) {
	return func(fileReader *fileReader) {
		fileReader.fileSystem = fileSystem
	}
}

// WithPath allows specifying an additional path (file or directory) to read config from
func WithPath(path string) func(*fileReader) {
	return func(fileReader *fileReader) {
		fileReader.paths = append(fileReader.paths, path)
	}
}
