package main

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type FlagsTestSuite struct {
	suite.Suite
	osEnviron map[string]string
	osArgs    []string
}

func (suite *FlagsTestSuite) SetupSuite() {
	// Flags
	suite.osArgs = make([]string, 0)
	suite.osArgs = append(suite.osArgs, os.Args...)

	// ENV
	suite.osEnviron = make(map[string]string)
	for _, e := range []string{
		"ADDRESS",
		"STORE_INTERVAL",
		"FILE_STORAGE_PATH",
		"RESTORE",
		"DATABASE_DSN",
		"KEY",
		"CRYPTO_KEY",
		"TRUSTED_SUBNET",
	} {
		suite.osEnviron[e] = os.Getenv(e)
	}
}

func (suite *FlagsTestSuite) TearDownSuite() {
	// Flags
	os.Args = make([]string, 0)
	os.Args = append(os.Args, suite.osArgs...)

	// ENV
	for k, v := range suite.osEnviron {
		if v != "" {
			_ = os.Setenv(k, v)
		}
	}
}

func (suite *FlagsTestSuite) SetupSubTest() {
	// Flags
	os.Args = make([]string, 0)
	os.Args = append(os.Args, suite.osArgs[0])

	// Prepare a new default flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Clear ENV
	for k := range suite.osEnviron {
		_ = os.Unsetenv(k)
	}
}

func (suite *FlagsTestSuite) TestParseFlags() {
	testCases := []struct {
		name string
		want map[string]interface{}
		args []string
	}{
		{
			name: "Positive case: Default values",
			args: nil,
			want: map[string]interface{}{
				"serverAddr":      "localhost:8080",
				"storeInterval":   300 * time.Second,
				"fileStoragePath": "/tmp/metrics-db.json",
				"isReqRestore":    true,
				"databaseDSN":     "",
				"secretKey":       "",
				"privateKeyPath":  "",
			},
		},
		{
			name: "Positive case: Set flag -a",
			args: []string{"-a=example.com:8181"},
			want: map[string]interface{}{"serverAddr": "example.com:8181"},
		},
		{
			name: "Positive case: Set flag -i",
			args: []string{"-i=10"},
			want: map[string]interface{}{"storeInterval": 10 * time.Second},
		},
		{
			name: "Positive case: Set flag -f",
			args: []string{"-f=/temp/metrics-db.test.json"},
			want: map[string]interface{}{"fileStoragePath": "/temp/metrics-db.test.json"},
		},
		{
			name: "Positive case: Set flag -r (true)",
			args: []string{"-r=true"},
			want: map[string]interface{}{"isReqRestore": true},
		},
		{
			name: "Positive case: Set flag -r (false)",
			args: []string{"-r=false"},
			want: map[string]interface{}{"isReqRestore": false},
		},
		{
			name: "Positive case: Set flag -d",
			args: []string{"-d=postgres://username:password@localhost:5432/database_name"},
			want: map[string]interface{}{"databaseDSN": "postgres://username:password@localhost:5432/database_name"},
		},
		{
			name: "Positive case: Set flag -k",
			args: []string{"-k=secret"},
			want: map[string]interface{}{"secretKey": "secret"},
		},
		{
			name: "Positive case: Set flag -crypto-key",
			args: []string{"-crypto-key=/tmp/key"},
			want: map[string]interface{}{"privateKeyPath": "/tmp/key"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			config := NewConfig()
			config, err := parseFlags(config)
			suite.Require().NoError(err)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
			}
		})
	}
}

func (suite *FlagsTestSuite) TestParseEnvs() {
	testCases := []struct {
		name string
		want map[string]interface{}
		envs []string
	}{
		{
			name: "Positive case: Set env ADDRESS",
			envs: []string{"ADDRESS=example.com:8181"},
			want: map[string]interface{}{"serverAddr": "example.com:8181"},
		},
		{
			name: "Positive case: Set env STORE_INTERVAL",
			envs: []string{"STORE_INTERVAL=10"},
			want: map[string]interface{}{"storeInterval": 10 * time.Second},
		},
		{
			name: "Positive case: Set env FILE_STORAGE_PATH",
			envs: []string{"FILE_STORAGE_PATH=/temp/metrics-db.test.json"},
			want: map[string]interface{}{"fileStoragePath": "/temp/metrics-db.test.json"},
		},
		{
			name: "Positive case: Set env RESTORE (true)",
			envs: []string{"RESTORE=true"},
			want: map[string]interface{}{"isReqRestore": true},
		},
		{
			name: "Positive case: Set env RESTORE (false)",
			envs: []string{"RESTORE=false"},
			want: map[string]interface{}{"isReqRestore": false},
		},
		{
			name: "Positive case: Set env DATABASE_DSN",
			envs: []string{"DATABASE_DSN=postgres://username:password@localhost:5432/database_name"},
			want: map[string]interface{}{"databaseDSN": "postgres://username:password@localhost:5432/database_name"},
		},
		{
			name: "Positive case: Set env KEY",
			envs: []string{"KEY=secret"},
			want: map[string]interface{}{"secretKey": "secret"},
		},
		{
			name: "Positive case: Set env CRYPTO_KEY",
			envs: []string{"CRYPTO_KEY=/tmp/key"},
			want: map[string]interface{}{"privateKeyPath": "/tmp/key"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.envs) > 0 {
				for _, v := range tc.envs {
					e := strings.Split(v, "=")

					suite.Require().GreaterOrEqual(len(e), 2)

					if len(e) > 2 {
						e[1] = strings.Join(e[1:], "=")
						e = e[:2]
					}

					err := os.Setenv(e[0], e[1])
					suite.Require().NoError(err)
				}
			}

			config := NewConfig()
			config, err := parseEnvs(config)
			suite.Require().NoError(err)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
			}
		})
	}
}

func (suite *FlagsTestSuite) TestLoadConfig() {
	testCases := []struct {
		name string
		want map[string]interface{}
		args []string
		envs []string
	}{
		{
			name: "Positive case: Default values",
			args: nil,
			envs: nil,
			want: map[string]interface{}{
				"serverAddr":      "localhost:8080",
				"storeInterval":   300 * time.Second,
				"fileStoragePath": "/tmp/metrics-db.json",
				"isReqRestore":    true,
				"secretKey":       "",
				"privateKeyPath":  "",
			},
		},
		{
			name: "Positive case: Set flag -a and env ADDRESS",
			args: []string{"-a=aaa.com:3333"},
			envs: []string{"ADDRESS=bbb.com:5555"},
			want: map[string]interface{}{"serverAddr": "bbb.com:5555"},
		},
		{
			name: "Positive case: Set flag -a only",
			args: []string{"-a=aaa.com:3333"},
			envs: nil,
			want: map[string]interface{}{"serverAddr": "aaa.com:3333"},
		},
		{
			name: "Positive case: Set env ADDRESS only",
			args: nil,
			envs: []string{"ADDRESS=bbb.com:5555"},
			want: map[string]interface{}{"serverAddr": "bbb.com:5555"},
		},
		{
			name: "Positive case: Set flag -i and env STORE_INTERVAL",
			args: []string{"-i=100"},
			envs: []string{"STORE_INTERVAL=200"},
			want: map[string]interface{}{"storeInterval": 200 * time.Second},
		},
		{
			name: "Positive case: Set flag -i only",
			args: []string{"-i=100"},
			envs: nil,
			want: map[string]interface{}{"storeInterval": 100 * time.Second},
		},
		{
			name: "Positive case: Set env STORE_INTERVAL only",
			args: nil,
			envs: []string{"STORE_INTERVAL=200"},
			want: map[string]interface{}{"storeInterval": 200 * time.Second},
		},
		{
			name: "Positive case: Set flag -f and env FILE_STORAGE_PATH",
			args: []string{"-f=/temp/metrics-db.test1.json"},
			envs: []string{"FILE_STORAGE_PATH=/temp/metrics-db.test2.json"},
			want: map[string]interface{}{"fileStoragePath": "/temp/metrics-db.test2.json"},
		},
		{
			name: "Positive case: Set flag -f only",
			args: []string{"-f=/temp/metrics-db.test1.json"},
			envs: nil,
			want: map[string]interface{}{"fileStoragePath": "/temp/metrics-db.test1.json"},
		},
		{
			name: "Positive case: Set env FILE_STORAGE_PATH only",
			args: nil,
			envs: []string{"FILE_STORAGE_PATH=/temp/metrics-db.test2.json"},
			want: map[string]interface{}{"fileStoragePath": "/temp/metrics-db.test2.json"},
		},
		{
			name: "Positive case: Set flag -f (false) and env RESTORE",
			args: []string{"-r=false"},
			envs: []string{"RESTORE=true"},
			want: map[string]interface{}{"isReqRestore": true},
		},
		{
			name: "Positive case: Set flag -t (true) and env RESTORE",
			args: []string{"-r=true"},
			envs: []string{"RESTORE=false"},
			want: map[string]interface{}{"isReqRestore": false},
		},
		{
			name: "Positive case: Set flag -t (true) only",
			args: []string{"-r=true"},
			envs: nil,
			want: map[string]interface{}{"isReqRestore": true},
		},
		{
			name: "Positive case: Set flag -r (false) only",
			args: []string{"-r=false"},
			envs: nil,
			want: map[string]interface{}{"isReqRestore": false},
		},
		{
			name: "Positive case: Set env RESTORE (true) only",
			args: nil,
			envs: []string{"RESTORE=true"},
			want: map[string]interface{}{"isReqRestore": true},
		},
		{
			name: "Positive case: Set env RESTORE (false) only",
			args: nil,
			envs: []string{"RESTORE=false"},
			want: map[string]interface{}{"isReqRestore": false},
		},
		{
			name: "Positive case: Set flag -d and env DATABASE_DSN",
			args: []string{"-d=postgres://username1:password1@localhost:5432/database_name1"},
			envs: []string{"DATABASE_DSN=postgres://username2:password2@localhost:5432/database_name2"},
			want: map[string]interface{}{"databaseDSN": "postgres://username2:password2@localhost:5432/database_name2"},
		},
		{
			name: "Positive case: Set flag -d only",
			args: []string{"-d=postgres://username1:password1@localhost:5432/database_name1"},
			envs: nil,
			want: map[string]interface{}{"databaseDSN": "postgres://username1:password1@localhost:5432/database_name1"},
		},
		{
			name: "Positive case: Set env DATABASE_DSN only",
			args: nil,
			envs: []string{"DATABASE_DSN=postgres://username2:password2@localhost:5432/database_name2"},
			want: map[string]interface{}{"databaseDSN": "postgres://username2:password2@localhost:5432/database_name2"},
		},
		{
			name: "Positive case: Set flag -k and env KEY",
			args: []string{"-k=secret1"},
			envs: []string{"KEY=secret2"},
			want: map[string]interface{}{"secretKey": "secret2"},
		},
		{
			name: "Positive case: Set flag -k only",
			args: []string{"-k=secret"},
			envs: nil,
			want: map[string]interface{}{"secretKey": "secret"},
		},
		{
			name: "Positive case: Set env KEY only",
			args: nil,
			envs: []string{"KEY=secret"},
			want: map[string]interface{}{"secretKey": "secret"},
		},
		{
			name: "Positive case: Set flag -crypto-key and env CRYPTO_KEY",
			args: []string{"-crypto-key=/tmp/key1"},
			envs: []string{"CRYPTO_KEY=/tmp/key2"},
			want: map[string]interface{}{"privateKeyPath": "/tmp/key2"},
		},
		{
			name: "Positive case: Set flag -crypto-key only",
			args: []string{"-crypto-key=/tmp/key"},
			envs: nil,
			want: map[string]interface{}{"privateKeyPath": "/tmp/key"},
		},
		{
			name: "Positive case: Set env CRYPTO_KEY only",
			args: nil,
			envs: []string{"CRYPTO_KEY=/tmp/key"},
			want: map[string]interface{}{"privateKeyPath": "/tmp/key"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			if len(tc.envs) > 0 {
				for _, v := range tc.envs {
					e := strings.Split(v, "=")

					suite.Require().GreaterOrEqual(len(e), 2)

					if len(e) > 2 {
						e[1] = strings.Join(e[1:], "=")
						e = e[:2]
					}

					err := os.Setenv(e[0], e[1])
					suite.Require().NoError(err)
				}
			}

			config, err := loadConfig()
			suite.Require().NoError(err)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
			}
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
