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
				Credentials: CredentialsList{
					Production:  Credentials{AwsKeyId: "PROD_KEY_ID", AwsSecretKey: "PROD_SECRET_KEY"},
					Development: Credentials{AwsKeyId: "DEV_KEY_ID", AwsSecretKey: "DEV_SECRET_KEY"},
					Test:        Credentials{AwsKeyId: "TEST_KEY_ID", AwsSecretKey: "TEST_SECRET_KEY"},
				},
				Constraints: Constraints{
					Regions: []string{"eu-west-2", "ca-central-1"},
					Services: []ServiceDescription{
						{
							Name:        "Test Service 1",
							MinMemory:   16,
							MaxVcpu:     4,
							Focus:       "performance",
							FocusWeight: 0.9,
							Instances: ServiceDescriptionInstances{
								MinimumCount: 0,
								TotalCount:   2,
							},
						},
					},
				},
				ApiConfig: ApiConfig{
					Port: 123456,
				},
				AwsApiConfig: AwsApiConfig{
					Endpoints: Endpoints{
						AwsSpotInstanceInfoUrl: "TEST_URL",
					},
					MaxInstancesToFetch:   1000,
					DownloadsDir:          "TEST_DOWNLOADS_DIR",
					ConsiderFreeInstances: true,
				},
				CacheConfig: CacheConfig{
					Dirpath: "TEST_CACHE_DIRPATH",
				},
			},
		},
		{
			filepath: "testdata/valid/config-2.json",
			expected: Config{
				Credentials: CredentialsList{
					Production:  Credentials{AwsKeyId: "PROD_KEY_ID", AwsSecretKey: "PROD_SECRET_KEY"},
					Development: Credentials{AwsKeyId: "DEV_KEY_ID", AwsSecretKey: "DEV_SECRET_KEY"},
					Test:        Credentials{AwsKeyId: "TEST_KEY_ID", AwsSecretKey: "TEST_SECRET_KEY"},
				},
				Constraints: Constraints{
					Regions: []string{"us-west-1"},
					Services: []ServiceDescription{
						{
							Name:        "Test Service 1",
							MinMemory:   4,
							MaxVcpu:     1,
							Focus:       "cost",
							FocusWeight: 0.4,
							Instances: ServiceDescriptionInstances{
								MinimumCount: 1,
								TotalCount:   3,
							},
						},
						{
							Name:        "Test Service 2",
							MinMemory:   16,
							MaxVcpu:     4,
							Focus:       "balanced",
							FocusWeight: 0.9,
							Instances: ServiceDescriptionInstances{
								MinimumCount: 0,
								TotalCount:   2,
							},
						},
					},
				},
				ApiConfig: ApiConfig{
					Port: 54321,
				},
				AwsApiConfig: AwsApiConfig{
					Endpoints: Endpoints{
						AwsSpotInstanceInfoUrl: "TEST_URL",
					},
					DownloadsDir:          "TEST_DOWNLOADS_DIR",
					MaxInstancesToFetch:   1000,
					ConsiderFreeInstances: false,
				},
				CacheConfig: CacheConfig{
					Dirpath: "TEST_CACHE_DIRPATH",
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
		"duplicate service names": {filepath: "testdata/invalid/duplicate-service-names.json"},
		"empty service names":     {filepath: "testdata/invalid/empty-service-names.json"},
		"invalid focus":           {filepath: "testdata/invalid/invalid-focus.json"},
		"invalid regions":         {filepath: "testdata/invalid/invalid-regions.json"},
		"no API config":           {filepath: "testdata/invalid/no-aws-api-config.json"},
		"no credentials":          {filepath: "testdata/invalid/no-credentials.json"},
		"no regions":              {filepath: "testdata/invalid/no-regions.json"},
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
