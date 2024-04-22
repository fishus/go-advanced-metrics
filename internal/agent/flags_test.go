package agent

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
		"POLL_INTERVAL",
		"REPORT_INTERVAL",
		"KEY",
		"CRYPTO_KEY",
		"RATE_LIMIT",
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
				"serverAddr":     "localhost:8080",
				"pollInterval":   2 * time.Second,
				"reportInterval": 10 * time.Second,
				"secretKey":      "",
				"publicKeyPath":  "",
				"rateLimit":      uint(3),
			},
		},
		{
			name: "Positive case: Set flag -a",
			args: []string{"-a=example.com:8181"},
			want: map[string]interface{}{"serverAddr": "example.com:8181"},
		},
		{
			name: "Positive case: Set flag -p",
			args: []string{"-p=3"},
			want: map[string]interface{}{"pollInterval": 3 * time.Second},
		},
		{
			name: "Positive case: Set flag -r",
			args: []string{"-r=7"},
			want: map[string]interface{}{"reportInterval": 7 * time.Second},
		},
		{
			name: "Positive case: Set flag -k",
			args: []string{"-k=secret"},
			want: map[string]interface{}{"secretKey": "secret"},
		},
		{
			name: "Positive case: Set flag -l",
			args: []string{"-l=15"},
			want: map[string]interface{}{"rateLimit": uint(15)},
		},
		{
			name: "Positive case: Set flag -crypto-key",
			args: []string{"-crypto-key=/tmp/key.pub"},
			want: map[string]interface{}{"publicKeyPath": "/tmp/key.pub"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			config := newConfig()
			config = parseFlags(config)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'. Value is: '%v'(%v), Should be: '%v'(%v)", k, field, field.Type(), reflect.ValueOf(want), reflect.ValueOf(want).Type())
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
			name: "Positive case: Default values",
			envs: nil,
			want: map[string]interface{}{
				"serverAddr":     "",
				"pollInterval":   0 * time.Second,
				"reportInterval": 0 * time.Second,
				"secretKey":      "",
				"publicKeyPath":  "",
				"rateLimit":      uint(0),
			},
		},
		{
			name: "Positive case: Set env ADDRESS",
			envs: []string{"ADDRESS=example.com:8181"},
			want: map[string]interface{}{"serverAddr": "example.com:8181"},
		},
		{
			name: "Positive case: Set env POLL_INTERVAL",
			envs: []string{"POLL_INTERVAL=3"},
			want: map[string]interface{}{"pollInterval": 3 * time.Second},
		},
		{
			name: "Positive case: Set env REPORT_INTERVAL",
			envs: []string{"REPORT_INTERVAL=7"},
			want: map[string]interface{}{"reportInterval": 7 * time.Second},
		},
		{
			name: "Positive case: Set env KEY",
			envs: []string{"KEY=secret"},
			want: map[string]interface{}{"secretKey": "secret"},
		},
		{
			name: "Positive case: Set env RATE_LIMIT",
			envs: []string{"RATE_LIMIT=15"},
			want: map[string]interface{}{"rateLimit": uint(15)},
		},
		{
			name: "Positive case: Set env CRYPTO_KEY",
			envs: []string{"CRYPTO_KEY=/tmp/key.pub"},
			want: map[string]interface{}{"publicKeyPath": "/tmp/key.pub"},
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

			config := newConfig()
			config = parseEnvs(config)

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'. Value is: '%v'(%v), Should be: '%v'(%v)", k, field, field.Type(), reflect.ValueOf(want), reflect.ValueOf(want).Type())
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
				"serverAddr":     "localhost:8080",
				"pollInterval":   2 * time.Second,
				"reportInterval": 10 * time.Second,
				"secretKey":      "",
				"publicKeyPath":  "",
				"rateLimit":      uint(3),
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
			name: "Positive case: Set flag -p and env POLL_INTERVAL",
			args: []string{"-p=21"},
			envs: []string{"POLL_INTERVAL=31"},
			want: map[string]interface{}{"pollInterval": 31 * time.Second},
		},
		{
			name: "Positive case: Set flag -p only",
			args: []string{"-p=21"},
			envs: nil,
			want: map[string]interface{}{"pollInterval": 21 * time.Second},
		},
		{
			name: "Positive case: Set env POLL_INTERVAL only",
			args: nil,
			envs: []string{"POLL_INTERVAL=31"},
			want: map[string]interface{}{"pollInterval": 31 * time.Second},
		},
		{
			name: "Positive case: Set flag -r and env REPORT_INTERVAL",
			args: []string{"-r=51"},
			envs: []string{"REPORT_INTERVAL=71"},
			want: map[string]interface{}{"reportInterval": 71 * time.Second},
		},
		{
			name: "Positive case: Set flag -r only",
			args: []string{"-r=51"},
			envs: nil,
			want: map[string]interface{}{"reportInterval": 51 * time.Second},
		},
		{
			name: "Positive case: Set env REPORT_INTERVAL only",
			args: nil,
			envs: []string{"REPORT_INTERVAL=71"},
			want: map[string]interface{}{"reportInterval": 71 * time.Second},
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
			name: "Positive case: Set flag -l and env RATE_LIMIT",
			args: []string{"-l=15"},
			envs: []string{"RATE_LIMIT=25"},
			want: map[string]interface{}{"rateLimit": uint(25)},
		},
		{
			name: "Positive case: Set flag -l only",
			args: []string{"-l=15"},
			envs: nil,
			want: map[string]interface{}{"rateLimit": uint(15)},
		},
		{
			name: "Positive case: Set env RATE_LIMIT only",
			args: nil,
			envs: []string{"RATE_LIMIT=15"},
			want: map[string]interface{}{"rateLimit": uint(15)},
		},
		{
			name: "Positive case: Set flag -crypto-key and env CRYPTO_KEY",
			args: []string{"-crypto-key=/tmp/key1.pub"},
			envs: []string{"CRYPTO_KEY=/tmp/key2.pub"},
			want: map[string]interface{}{"publicKeyPath": "/tmp/key2.pub"},
		},
		{
			name: "Positive case: Set flag -crypto-key only",
			args: []string{"-crypto-key=/tmp/key.pub"},
			envs: nil,
			want: map[string]interface{}{"publicKeyPath": "/tmp/key.pub"},
		},
		{
			name: "Positive case: Set env CRYPTO_KEY only",
			args: nil,
			envs: []string{"CRYPTO_KEY=/tmp/key.pub"},
			want: map[string]interface{}{"publicKeyPath": "/tmp/key.pub"},
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

			config := loadConfig()

			configFields := reflect.ValueOf(config)

			for k, want := range tc.want {
				field := configFields.FieldByName(k)
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'. Value is: '%v'(%v), Should be: '%v'(%v)", k, field, field.Type(), reflect.ValueOf(want), reflect.ValueOf(want).Type())
			}
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
