package config

import (
	"testing"
)

type validConfigTest struct {
	filepath string
	expected Config
}

type invalidConfigTest struct {
	filepath string
}

func TestParseConfigValid(t *testing.T) {

	tests := []validConfigTest{
		{
			filepath: "testdata/valid/config-1.json",
			expected: Config{
				Credentials: Credentials{AwsKeyId: "KEY_ID", AwsSecretKey: "SECRET_KEY"},
				ApiConfig: ApiConfig{
					AllowedDomains: []string{"https://test.com:3000"},
					Port:           123456,
				},
				AwsApiConfig: AwsApiConfig{
					Endpoints: Endpoints{
						AwsSpotInstanceInfoUrl: "TEST_URL",
					},
					MaxInstancesToFetch: 1000,
					DownloadsDir:        "TEST_DOWNLOADS_DIR",
				},
				CacheConfig: CacheConfig{
					Dirpath:         "TEST_CACHE_DIRPATH",
					DefaultLifetime: 100,
				},
			},
		},
		{
			filepath: "testdata/valid/config-2.json",
			expected: Config{
				Credentials: Credentials{AwsKeyId: "KEY_ID_2", AwsSecretKey: "SECRET_KEY_2"},
				ApiConfig: ApiConfig{
					Port:           54321,
					AllowedDomains: []string{"http://some.domain.com"},
				},
				AwsApiConfig: AwsApiConfig{
					Endpoints: Endpoints{
						AwsSpotInstanceInfoUrl: "TEST_URL",
					},
					DownloadsDir:        "TEST_DOWNLOADS_DIR",
					MaxInstancesToFetch: 1000,
				},
				CacheConfig: CacheConfig{
					Dirpath:         "TEST_CACHE_DIRPATH",
					DefaultLifetime: 200,
				},
			},
		},
	}

	for _, test := range tests {
		config, err := ParseConfig(test.filepath)
		if err != nil {
			t.Fatalf("Failed to parse config. File: %s, error: %s", test.filepath, err.Error())
		}
		if config.String() != test.expected.String() {
			t.Fatalf(
				"Parsed config not equal to expected config. File %s \n\nParsed:\n%s\n\nExpected:\n%s",
				test.filepath,
				config.String(),
				test.expected.String(),
			)
		}
	}
}

func TestParseConfigInvalid(t *testing.T) {

	errorTests := map[string]invalidConfigTest{
		"no AWS API config":       {filepath: "testdata/invalid/no-aws-api-config.json"},
		"no credentials":          {filepath: "testdata/invalid/no-credentials.json"},
		"no API config":           {filepath: "testdata/invalid/no-api-config.json"},
		"invalid port API config": {filepath: "testdata/invalid/invalid-port-config.json"},
	}

	for name, test := range errorTests {
		_, err := ParseConfig(test.filepath)
		if err == nil {
			t.Fatalf(
				"Expected error, but did not received for \"%s\". Filepath: %s",
				name,
				test.filepath,
			)
		}
	}
}
