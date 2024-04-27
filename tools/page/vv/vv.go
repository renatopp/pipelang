package vv

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type ProcessorFn func(*Ctx, *File) error

type File struct {
	// Represents the relative path to the file from the root directory. Includes
	// the file name and extension.
	SourcePath string

	// Represents the relative path that will be used when copying to dist
	// directory. Modify this to change the destination path.
	Dir string

	// Represents the file name without extension. Modify this to change the
	// destination file name.
	Name string

	// Represents the file extension. Modify this to change the destination file
	// extension.
	Extension string

	// Represents the file content. Modify this to change the destination file.
	Body []byte

	// If true, the file will be processed by the processors, but it won't be
	// written to the dist directory.
	Ignored bool

	// If true, the file won't be processed by the processors and won't be
	// written to the dist directory.
	Hidden bool
}

type Ctx struct {
	// Root dir is applied as prefix to all paths in the context, including dist.
	RootDir string

	// DistDir is the directory where the processed files will be written.
	DistDir string

	// Files in SourceDirs will be copied to DistDir, keeping the directory
	// structure. The file will be processed by processors.
	SourceDirs []string

	// Processors will be applied to each file in SourceDirs. Each processor
	// decides if the file should be processed or not.
	Processors []ProcessorFn
}

func Run(config *Ctx) error {
	files, err := loadFiles(config)
	if err != nil {
		return err
	}

	err = processFiles(config, files)
	if err != nil {
		return err
	}

	return copyFiles(config, files)
}

func ListFiles(dir string) ([]string, error) {
	files := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q (%v)", path, err)
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, nil
}

func LoadFile(path string) (*File, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %q (%v)", path, err)
	}

	dir := filepath.Dir(path)
	name := filepath.Base(path)
	extension := filepath.Ext(path)

	name = name[:len(name)-len(extension)]

	return &File{
		SourcePath: path,
		Dir:        dir,
		Name:       name,
		Extension:  extension,
		Body:       body,
	}, nil
}

func loadFiles(config *Ctx) ([]*File, error) {
	files := []*File{}

	for _, sourceDir := range config.SourceDirs {
		srcPath := filepath.Join(config.RootDir, sourceDir)

		paths, err := ListFiles(srcPath)
		if err != nil {
			return nil, err
		}

		for _, path := range paths {
			file, err := LoadFile(path)
			if err != nil {
				return nil, err
			}
			file.Dir = strings.Trim(file.Dir[len(srcPath):], "/\\")
			files = append(files, file)
		}
	}

	return files, nil
}

func processFiles(config *Ctx, files []*File) error {
	for _, processor := range config.Processors {
		for _, file := range files {
			if file.Hidden {
				continue
			}

			err := processor(config, file)
			if err != nil {
				name := runtime.FuncForPC(reflect.ValueOf(processor).Pointer()).Name()
				return fmt.Errorf("error processing file %q with processor %s (%v)", file.SourcePath, name, err)
			}
		}
	}
	return nil
}

func copyFiles(config *Ctx, files []*File) error {
	for _, file := range files {
		if file.Ignored || file.Hidden {
			continue
		}

		dist := filepath.Join(config.RootDir, config.DistDir, file.Dir, file.Name+file.Extension)
		err := os.MkdirAll(filepath.Dir(dist), 0755)
		if err != nil {
			return err
		}

		err = os.WriteFile(dist, file.Body, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
