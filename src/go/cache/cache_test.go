package cache

import (
	"ec2-test/utils"
	"errors"
	"testing"
	"time"
)

const (
	TEST_CACHE_DIRPATH = "./testdata"

	TEST_FILE_FILENAME       = "test-cache-file.txt"
	TEST_FILE_CONTENT        = "TEST CONTENT"
	TEST_FILE_ENTRY_SET_DATE = "2000-01-01T13:00:00Z"
	TEST_FILE_ENTRY_DURATION = "1h"
)

func setup() (
	testCacheDirpath string,
	testSetDate time.Time,
	testDuration time.Duration,
	err error,
) {
	cacheFileFilepath, err := utils.CreateFilepath(TEST_CACHE_DIRPATH, CACHE_FILENAME)
	if err != nil {
		err = utils.PrependToError(err, "could not create cache file's path")
		return
	}

	exists, err := utils.FileExists(cacheFileFilepath)
	if err != nil {
		err = utils.PrependToError(err, "failed to check if test cache file exists")
		return
	}
	if exists {
		err = errors.New("test cache file already exists")
		return
	}

	testSetDate, err = time.Parse(time.RFC3339, TEST_FILE_ENTRY_SET_DATE)
	if err != nil {
		err = utils.PrependToError(err, "failed to parse test entry set date")
		return
	}

	testDuration, err = time.ParseDuration(TEST_FILE_ENTRY_DURATION)
	if err != nil {
		err = utils.PrependToError(err, "failed to parse test entry duration")
		return
	}

	return
}

func cleanup() error {
	cacheFileFilepath, err := utils.CreateFilepath(TEST_CACHE_DIRPATH, CACHE_FILENAME)
	if err != nil {
		return utils.PrependToError(err, "could not create cache file's path")
	}

	exists, err := utils.FileExists(cacheFileFilepath)
	if err != nil {
		return utils.PrependToError(err, "error when checking if cache file exists")
	}
	if exists {
		err = utils.DeleteFile(cacheFileFilepath)
		if err != nil {
			return utils.PrependToError(err, "failed to delete cache file")
		}
	}
	return nil
}

func TestNew(t *testing.T) {
	_, _, _, err := setup()
	if err != nil {
		t.Fatalf("Setup failed: %s", err.Error())
	}

	cache, err := New(TEST_CACHE_DIRPATH)
	if err != nil || cache == nil {
		t.Fatalf("Error occurred when creating cache with no previous file: %s", err.Error())
	}

	cache, err = New(TEST_CACHE_DIRPATH)
	if err != nil || cache == nil {
		t.Fatalf("Error occurred when creating cache with previously existing file: %s", err.Error())
	}

	err = cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %s", err.Error())
	}
}

func TestGetEntry(t *testing.T) {
	_, testSetDate, testDuration, err := setup()
	if err != nil {
		t.Fatalf("Setup failed: %s", err.Error())
	}
	testInvalidationDate := testSetDate.Add(testDuration)

	cache, err := New(TEST_CACHE_DIRPATH)
	if err != nil {
		t.Fatalf("Error when creating cache: %s", err.Error())
	}

	setEntry := CacheEntry{
		Filename:         TEST_FILE_FILENAME,
		SetDate:          testSetDate,
		InvalidationDate: testInvalidationDate,
	}
	cache.Entries[TEST_FILE_FILENAME] = setEntry

	gotEntry, err := cache.GetEntry(TEST_FILE_FILENAME)
	if err != nil {
		t.Fatalf("Error when getting entry from cache: %s", err.Error())
	}

	if !gotEntry.SetDate.Equal(testSetDate) {
		t.Fatalf(
			"Set date retrieved from GetEntry does not match date set. Wanted: %v, got: %v",
			testSetDate,
			gotEntry.SetDate,
		)
	}
	if !gotEntry.InvalidationDate.Equal(testInvalidationDate) {
		t.Fatalf(
			"Invalidation date retrieved from GetEntry does not match date set. Wanted: %v, got: %v",
			testInvalidationDate,
			gotEntry.InvalidationDate,
		)
	}
	if gotEntry.Filename != TEST_FILE_FILENAME {
		t.Fatalf(
			"Filename returend from from GetEntry does not match set Filename. Wanted: %s, got: %s",
			TEST_FILE_FILENAME,
			gotEntry.Filename,
		)
	}

	err = cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %s", err.Error())
	}
}

func TestSet(t *testing.T) {
	_, _, testDuration, err := setup()
	if err != nil {
		t.Fatalf("Setup failed: %s", err.Error())
	}

	cache, err := New(TEST_CACHE_DIRPATH)
	if err != nil {
		t.Fatalf("Error when creating cache: %s", err.Error())
	}

	now := time.Now()
	earliestSetDate := now.Add(time.Second * -2)
	latestSetDate := now.Add(time.Second * 2)

	earliestInvalidationDate := earliestSetDate.Add(testDuration)
	latestInvalidationDate := latestSetDate.Add(testDuration)

	err = cache.Set(TEST_FILE_FILENAME, TEST_FILE_CONTENT, testDuration)
	if err != nil {
		t.Fatalf("Error when setting entry in cache: %s", err.Error())
	}

	entry, err := cache.GetEntry(TEST_FILE_FILENAME)
	if err != nil {
		t.Fatalf("Error when getting entry from cache: %s", err.Error())
	}

	if entry.Filename != TEST_FILE_FILENAME {
		t.Fatalf("Entry filename is incorrect. Wanted: %s, got: %s", TEST_FILE_FILENAME, entry.Filename)
	}
	if !(entry.SetDate.After(earliestSetDate) && entry.SetDate.Before(latestSetDate)) {
		t.Fatalf(
			"Entry set date is incorrect. Wanted between %s and %s, got: %s",
			earliestSetDate.String(),
			latestSetDate.String(),
			entry.SetDate.String(),
		)
	}
	if !(entry.InvalidationDate.After(earliestInvalidationDate) &&
		entry.InvalidationDate.Before(latestInvalidationDate)) {
		t.Fatalf(
			"Entry invalidation date is incorrect. Wanted between %s and %s, got: %s",
			earliestInvalidationDate.String(),
			latestInvalidationDate.String(),
			entry.InvalidationDate.String(),
		)
	}

	err = cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %s", err.Error())
	}
}

func TestIsValid(t *testing.T) {
	_, _, _, err := setup()
	if err != nil {
		t.Fatalf("Setup failed: %s", err.Error())
	}

	cache, err := New(TEST_CACHE_DIRPATH)
	if err != nil {
		t.Fatalf("Error creating cache: %s", err.Error())
	}

	isValid := cache.IsValid(TEST_FILE_FILENAME)
	if isValid {
		t.Fatalf(
			"Non-existent entry returned as valid from cache. Wanted: %t, got: %t",
			false,
			isValid,
		)
	}

	err = cache.Set(TEST_FILE_FILENAME, TEST_FILE_CONTENT, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("Error when setting value in cache: %s", err.Error())
	}

	time.Sleep(40 * time.Millisecond) // Wait until entry is invalid

	isValid = cache.IsValid(TEST_FILE_FILENAME)
	if isValid {
		t.Fatalf(
			"Invalid entry returned as valid from cache. Wanted: %t, got: %t",
			false,
			isValid,
		)
	}

	err = cache.Set(TEST_FILE_FILENAME, TEST_FILE_CONTENT, time.Hour)
	if err != nil {
		t.Fatalf("Error when setting value in cache: %s", err.Error())
	}

	isValid = cache.IsValid(TEST_FILE_FILENAME)
	if !isValid {
		t.Fatalf(
			"Valid entry returned as invalid from cache. Wanted: %t, got: %t",
			true,
			isValid,
		)
	}

	err = cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %s", err.Error())
	}
}

func TestGet(t *testing.T) {
	_, _, _, err := setup()
	if err != nil {
		t.Fatalf("Setup failed: %s", err.Error())
	}

	cache, err := New(TEST_CACHE_DIRPATH)
	if err != nil {
		t.Fatalf("Error creating cache: %s", err.Error())
	}

	_, err = cache.Get(TEST_FILE_FILENAME)
	if err == nil {
		t.Fatalf("Error not returned when getting a non-existent cache entry")
	}

	err = cache.Set(TEST_FILE_FILENAME, TEST_FILE_CONTENT, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Error when setting value in cache: %s", err.Error())
	}

	time.Sleep(55 * time.Millisecond) // Wait until entry is invalid
	_, err = cache.Get(TEST_FILE_FILENAME)
	if err == nil {
		t.Fatalf("Error not returned when getting invalid entry: %s", err.Error())
	}

	err = cache.Set(TEST_FILE_FILENAME, TEST_FILE_CONTENT, time.Hour)
	if err != nil {
		t.Fatalf("Error when setting value in cache: %s", err.Error())
	}

	content, err := cache.Get(TEST_FILE_FILENAME)
	if err != nil {
		t.Fatalf("Error when getting valid value from cache: %s", err.Error())
	}
	if content != TEST_FILE_CONTENT {
		t.Fatalf(
			"Incorrect file content returned from cache. Wanted: %s, got: %s",
			TEST_FILE_CONTENT,
			content,
		)
	}

	err = cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %s", err.Error())
	}
}

func TestParse(t *testing.T) {
	_, _, testDuration, err := setup()
	if err != nil {
		t.Fatalf("Setup failed: %s", err.Error())
	}

	_, err = ParseCache("./NON_EXISTENT_DIR")
	if err == nil {
		t.Fatalf(
			"No error returned when parsing cache from non-existent file. Filepath from test: %s",
			"./NON_EXISTENT_DIR",
		)
	}

	cache, err := New(TEST_CACHE_DIRPATH)
	if err != nil {
		t.Fatalf("Error creating new cache: %s", err.Error())
	}

	err = cache.Set(TEST_FILE_FILENAME, TEST_FILE_CONTENT, testDuration)
	if err != nil {
		t.Fatalf("Error when setting cache entry: %s", err.Error())
	}

	setDate := cache.Entries[TEST_FILE_FILENAME].SetDate
	invalidationDate := setDate.Add(testDuration)

	cache, err = ParseCache(TEST_CACHE_DIRPATH)
	if err != nil {
		t.Fatalf("Error when parsing cache from existing file: %s", err.Error())
	}

	content, err := cache.Get(TEST_FILE_FILENAME)
	if err != nil {
		t.Fatalf("Error when getting file content from cache: %s", err.Error())
	}
	if content != TEST_FILE_CONTENT {
		t.Fatalf("Incorrect file content returned. Wanted: %s, got: %s", TEST_FILE_CONTENT, content)
	}

	entry, err := cache.GetEntry(TEST_FILE_FILENAME)
	if err != nil {
		t.Fatalf("Error when getting entry for test file from cache: %s", err.Error())
	}
	if entry.Filename != TEST_FILE_FILENAME {
		t.Fatalf("Incorrect filepath retrieved. Wanted: %s, got: %s", TEST_FILE_FILENAME, entry.Filename)
	}
	if !entry.SetDate.Equal(setDate) {
		t.Fatalf(
			"Incorrect set date retrieved. Wanted: %s, got: %s",
			setDate.String(),
			entry.SetDate.String(),
		)
	}
	if !entry.InvalidationDate.Equal(invalidationDate) {
		t.Fatalf(
			"Incorrect invalidation date retrieved. Wanted: %s, got: %s",
			setDate.Add(testDuration).String(),
			entry.InvalidationDate.String(),
		)
	}

	err = cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %s", err.Error())
	}
}
