package cache

import (
	"ec2-test/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// TODO: Doc comments
// TODO: Use

type Cache struct {
	Filepath string                `json:"-"`
	Entries  map[string]CacheEntry `json:"entries"`
}

type CacheEntry struct {
	SetDate          time.Time `json:"setAt"`
	InvalidationDate time.Time `json:"invalidFrom"`
	Filepath         string    `json:"filepath"`
}

// Creates and returns a new cache using the given filepath as the location for the cache file.
// Returns an error if an error occurred when writing to the given filepath.
func New(cacheFilepath string) (*Cache, error) {
	err := utils.WriteBytesToFile(make([]byte, 0), cacheFilepath)
	if err != nil {
		return nil, utils.PrependToError(err, "failed to write to cache file")
	}
	return &Cache{
		Filepath: cacheFilepath,
		Entries:  make(map[string]CacheEntry),
	}, nil
}

func ParseCache(cacheFilepath string) (*Cache, error) {
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

	c.Filepath = cacheFilepath
	return &c, nil
}

func ParseIfExistsElseNew(cacheFilepath string) (*Cache, error) {
	exists, err := utils.FileExists(cacheFilepath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return New(cacheFilepath)
	}
	return ParseCache(cacheFilepath)
}

// TODO: Test
// TODO: Doc
// Returns true if the given filepath is in the cache.
// Returns false otherwise.
// Returns an error if the given
func (cache Cache) Contains(filepath string) (bool, error) {
	absFilepath, err := utils.AbsoluteFilepath(filepath)
	if err != nil {
		return false, err
	}
	return cache.containsAbsPath(absFilepath), nil
}

// Returns true if the given absolute filepath is in the cache (valid OR invalid).
// Returns false otherwise.
func (cache Cache) containsAbsPath(absFilepath string) bool {
	_, exists := cache.Entries[absFilepath]
	return exists
}

// TODO: Test
// Returns true if the given filepath is in the cache and is valid.
// Returns false otherwise.
// Returns an error if the filepath is not a valid UNIX filepath.
func (cache Cache) IsValid(filepath string) (bool, error) {
	absFilepath, err := utils.AbsoluteFilepath(filepath)
	if err != nil {
		return false, err
	}
	return cache.isValidAbsPath(absFilepath), nil
}

// Returns true if the given absolute filepath is in the cache and is valid.
// Returns false otherwise.
func (cache Cache) isValidAbsPath(absFilepath string) bool {
	if !cache.containsAbsPath(absFilepath) {
		return false
	}
	entry := cache.getEntryAbsPath(absFilepath)
	return entry.InvalidationDate.After(time.Now())
}

// TODO: Test & convert return to stream
// Returns the cached file's content if the given filepath is in the cache and is valid.
// Returns an error otherwise or if the filepath is not a valid UNIX filepath.
func (cache Cache) Get(filepath string) (string, error) {
	absFilepath, err := utils.AbsoluteFilepath(filepath)
	if err != nil {
		return "", err
	}
	if !cache.isValidAbsPath(absFilepath) {
		return "", errors.New("cache entry is invalid")
	}
	entry := cache.getEntryAbsPath(absFilepath)
	return utils.FileToString(entry.Filepath)
}

// Returns the CacheEntry associated with the given filepath.
// Returns an error if the given filepath is not a valid UNIX path
// or if the entry does not exist in the cache.
func (cache Cache) GetEntry(filepath string) (*CacheEntry, error) {
	absFilepath, err := utils.AbsoluteFilepath(filepath)
	if err != nil {
		return nil, err
	}
	if !cache.containsAbsPath(absFilepath) {
		return nil, fmt.Errorf("entry does not exist: %s", filepath)
	}
	entry := cache.getEntryAbsPath(absFilepath)
	return &entry, nil
}

func (cache Cache) getEntryAbsPath(absFilepath string) CacheEntry {
	return cache.Entries[absFilepath]
}

// TODO: Use stream, not string
// Writes fileContent to a file at filepath and creates a valid entry in the cache.
// The entry is valid for the duration of lifetime.
// Returns an error if the file could not be written or value could not be set.
func (cache Cache) Set(filepath string, fileContent string, lifetime time.Duration) error {
	absFilepath, err := utils.AbsoluteFilepath(filepath)
	if err != nil {
		return err
	}

	err = utils.WriteStringToFile(fileContent, absFilepath)
	if err != nil {
		return utils.PrependToError(err, "failed to write cache entry to file")
	}

	oldEntry := cache.getEntryAbsPath(absFilepath)

	now := time.Now()
	invalid := now.Add(lifetime)

	cache.Entries[absFilepath] = CacheEntry{
		SetDate:          now,
		InvalidationDate: invalid,
		Filepath:         absFilepath,
	}

	// TODO: Improve this consistency
	err = cache.writeToFile()
	if err != nil {
		cache.Entries[filepath] = oldEntry
		return utils.PrependToError(
			err,
			"failed to write cache to disk. WARNING: file data may have been written",
		)
	}

	return nil
}

// Writes cache to its given filepath in JSON format
func (cache Cache) writeToFile() error {
	cacheAsJsonBytes, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	err = utils.WriteBytesToFile(cacheAsJsonBytes, cache.Filepath)
	return err
}
