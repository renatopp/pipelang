package runtime

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/renatopp/pipelang/internal/logs"
)

type FileCache struct {
	mutex sync.Mutex
}

func NewFileCache() *FileCache {
	return &FileCache{
		mutex: sync.Mutex{},
	}
}

// Load the file from the cache if it exists or from the disk if it doesn't.
// - Path should be the absolute path to the file
func (c *FileCache) Load(path string) (*SourceFile, error) {
	logs.Print("[filecache] loading file (%s)", path)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	file, err := c.createFileStruct(path)
	if err != nil {
		return nil, err
	}

	isCacheValid, err := c.checkCache(file)
	if err != nil {
		return nil, err
	}

	if isCacheValid {
		cacheFile, err := c.loadCacheFile(file)
		if err == nil {
			file.ast = cacheFile.Ast
		}
	}

	if file.ast == nil {
		_, err := file.LoadAst()
		if err != nil {
			return nil, err
		}

		c.saveCacheFile(file)
	}

	return file, nil
}

// Create a struct to hold the processed information
func (c *FileCache) createFileStruct(sourcePath string) (*SourceFile, error) {
	hash := hashString(sourcePath)
	dir, err := os.UserCacheDir()
	if err != nil {
		logs.Print("[runtime] error getting cache dir (%v)", err)
		return nil, err
	}

	dir = filepath.Join(dir, "pipelang")
	if err := os.MkdirAll(dir, 0755); err != nil {
		logs.Print("[runtime] error guaranteeing cache dir (%v)", err)
		return nil, err
	}

	cachePath := filepath.Join(dir, hash)

	return &SourceFile{
		hash:       hash,
		sourcePath: sourcePath,
		cachePath:  cachePath,
	}, nil
}

// Check if the cache is valid (if it exists and is up to date with the source
// file)
func (c *FileCache) checkCache(file *SourceFile) (bool, error) {
	cacheStat, err := os.Stat(file.CachePath())
	if err != nil {
		return false, nil
	}

	sourceStat, err := os.Stat(file.SourcePath())
	if err != nil {
		logs.Print("[runtime] error getting source stats (%v)", err)
		return false, err
	}

	// logs.Print("[filecache] checking cache (%s) (%s)", file.CachePath(), file.SourcePath())
	// logs.Print("[filecache] checking date (%s) (%s)", cacheStat.ModTime().String(), sourceStat.ModTime().String())
	return cacheStat.ModTime().After(sourceStat.ModTime()), nil
}

// Load the cache file from the disk
func (c *FileCache) loadCacheFile(file *SourceFile) (*CacheFile, error) {
	logs.Print("[filecache] loading cache file (%s)", file.CachePath())

	data, err := os.ReadFile(file.CachePath())
	if err != nil {
		logs.Print("[filecache] error reading cache file (%v)", err)
		return nil, err
	}

	cacheFile := &CacheFile{}
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(cacheFile)
	if err != nil {
		logs.Print("[filecache] error decoding cache file (%v)", err)
		return nil, err
	}

	return cacheFile, nil
}

// Save the AST to the cache
func (c *FileCache) saveCacheFile(file *SourceFile) error {
	logs.Print("[filecache] saving cache file (%s)", file.CachePath())

	cacheFile := &CacheFile{
		SourcePath: file.SourcePath(),
		Ast:        file.ast,
	}

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(cacheFile)
	if err != nil {
		logs.Print("[filecache] error encoding cache file (%v)", err)
		return err
	}

	err = os.WriteFile(file.CachePath(), buf.Bytes(), 0644)
	if err != nil {
		logs.Print("[filecache] error writing cache file (%v)", err)
		return err
	}

	return nil
}

func hashString(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}
