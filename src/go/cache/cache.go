package cache

import (
	"ec2-test/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// TODO: Doc comments for types
const (
	CACHE_FILENAME = "cache.json"
)

type Cache struct {
	Dirpath string                `json:"-"`
	Entries map[string]CacheEntry `json:"entries"`
}

type CacheEntry struct {
	SetDate          time.Time `json:"setAt"`
	InvalidationDate time.Time `json:"invalidFrom"`
	Filename         string    `json:"file"`
}

// Creates and returns a new cache.
// The given cacheDirpath is used as the directory to store cached files and cache metadata.
// Previous caches using cacheDirpath are overwritten.
// Returns an error if an error occurred when writing to the given filepath.
func New(cacheDirpath string) (*Cache, error) {
	cacheFilepath, err := getCacheFileFilepath(cacheDirpath)
	if err != nil {
		return nil, utils.PrependToError(err, "failed to generate cache filepath")
	}

	err = utils.WriteBytesToFile(make([]byte, 0), cacheFilepath)
	if err != nil {
		return nil, utils.PrependToError(err, "failed to write to cache file")
	}
	return &Cache{
		Dirpath: cacheDirpath,
		Entries: make(map[string]CacheEntry),
	}, nil
}

// Parses and returns a cache created at the given cacheDirpath.
// Returns an error if no cache exists at cacheDirpath or if there was an issue reading files.
func ParseCache(cacheDirpath string) (*Cache, error) {
	cacheFilepath, err := getCacheFileFilepath(cacheDirpath)
	if err != nil {
		return nil, utils.PrependToError(err, "failed to generate cache filepath")
	}

	cacheFileBytes, err := utils.FileToBytes(cacheFilepath)
	if err != nil {
		return nil, utils.PrependToError(
			err,
			fmt.Sprintf("failed to read cache file from filepath %s", cacheFilepath),
		)
	}

	var c Cache
	err = json.Unmarshal(cacheFileBytes, &c)
	if err != nil {
		return nil, err
	}
	c.Dirpath = cacheDirpath

	return &c, nil
}

// Parses and returns a cache if one exists at the given cacheDirpath.
// Creates and returns a new cache otherwise.
// Returns an error if an error occurs during the checking or creation process.
func ParseIfExistsElseNew(cacheDirpath string) (*Cache, error) {
	cacheFilepath, err := getCacheFileFilepath(cacheDirpath)
	if err != nil {
		return nil, utils.PrependToError(err, "failed to generate cache filepath")
	}

	exists, err := utils.FileExists(cacheFilepath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return New(cacheDirpath)
	}
	return ParseCache(cacheDirpath)
}

// TODO: Test
// Returns true if the given file is in the cache.
// Returns false otherwise.
func (cache Cache) Contains(file string) bool {
	_, exists := cache.Entries[file]
	return exists
}

// Returns the CacheEntry for the file if it's in the cache.
// Returns an error otherwise.
func (cache Cache) GetEntry(file string) (*CacheEntry, error) {
	if !cache.Contains(file) {
		return nil, fmt.Errorf("entry does not exist: %s", file)
	}
	entry := cache.Entries[file]
	return &entry, nil
}

// TODO: Test
// Returns true if the given file is in the cache and is valid.
// Returns false otherwise.
func (cache Cache) IsValid(file string) bool {
	if !cache.Contains(file) {
		return false
	}
	entry, err := cache.GetEntry(file)
	if err != nil {
		return false
	}
	return entry.InvalidationDate.After(time.Now())
}

// TODO: Test & convert return to stream
// Returns the cached file's contents if the given file is in the cache and is valid.
// Returns an error otherwise or if an error occurs during the file reading process.
func (cache Cache) Get(file string) (string, error) {
	if !cache.IsValid(file) {
		return "", errors.New("cache entry is invalid")
	}
	contents, err := cache.GetIgnoringValidity(file)
	if err != nil {
		return "", err
	}
	return contents, nil
}

// TODO: Test
// Returns the cached file's contents if the given file is in the cache, regardless of its validity.
// Returns an error otherwise or if an error occurs during the file reading process.
func (cache Cache) GetIgnoringValidity(file string) (string, error) {
	if !cache.Contains(file) {
		return "", fmt.Errorf("cache does not contain file: %s", file)
	}

	path, err := cache.getFileFilepath(file)
	if err != nil {
		return "", err
	}

	return utils.FileToString(path)
}

// TODO: Use stream, not string
// Writes fileContent to a file named filename in the cache directory and creates a
// representative entry in the cache, which is valid for a duration of lifetime.
// Returns an error if the file could not be written or value could not be set.
func (cache Cache) Set(filename string, fileContent string, lifetime time.Duration) error {
	path, err := utils.CreateFilepath(cache.Dirpath, filename)
	if err != nil {
		return err
	}

	err = utils.WriteStringToFile(fileContent, path)
	if err != nil {
		return utils.PrependToError(err, "failed to write cache entry to file")
	}

	oldEntry, oldEntryErr := cache.GetEntry(filename)

	now := time.Now()
	invalid := now.Add(lifetime)

	cache.Entries[filename] = CacheEntry{
		SetDate:          now,
		InvalidationDate: invalid,
		Filename:         filename,
	}

	err = cache.writeToFile()
	if err != nil {
		if oldEntryErr == nil {
			cache.Entries[filename] = *oldEntry
		}
		return utils.PrependToError(
			err,
			"failed to write cache to disk. WARNING: file data may be inconsistent with in-memory cache",
		)
	}

	return nil
}

// TODO: Doc
func (cache Cache) getFileFilepath(file string) (string, error) {
	return utils.CreateFilepath(cache.Dirpath, file)
}

// TODO: Doc
func getCacheFileFilepath(cacheDirpath string) (string, error) {
	return utils.CreateFilepath(cacheDirpath, CACHE_FILENAME)
}

// TODO: Doc
func (cache Cache) getCacheFileFilepath() (string, error) {
	return getCacheFileFilepath(cache.Dirpath)
}

// Writes cache to its given filepath in JSON format.
func (cache Cache) writeToFile() error {
	cacheFilepath, err := cache.getCacheFileFilepath()
	if err != nil {
		return utils.PrependToError(err, "failed to generate cache filepath")
	}

	cacheAsJsonBytes, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = utils.WriteBytesToFile(cacheAsJsonBytes, cacheFilepath)
	return err
}
