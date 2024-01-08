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
	osArgs    []string
	osEnviron map[string]string
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
		args []string
		want map[string]interface{}
	}{
		{
			name: "Positive case #1",
			args: nil,
			want: map[string]interface{}{
				"serverAddr":     "localhost:8080",
				"pollInterval":   2 * time.Second,
				"reportInterval": 10 * time.Second,
			},
		},
		{
			name: "Positive case #2",
			args: []string{"-a=example.com:8181"},
			want: map[string]interface{}{"serverAddr": "example.com:8181"},
		},
		{
			name: "Positive case #3",
			args: []string{"-p=3"},
			want: map[string]interface{}{"pollInterval": 3 * time.Second},
		},
		{
			name: "Positive case #4",
			args: []string{"-r=7"},
			want: map[string]interface{}{"reportInterval": 7 * time.Second},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.args != nil && len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			config := NewConfig()
			config = parseFlags(config)

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
		envs []string
		want map[string]interface{}
	}{
		{
			name: "Positive case #1",
			envs: nil,
			want: map[string]interface{}{
				"serverAddr":     "",
				"pollInterval":   0 * time.Second,
				"reportInterval": 0 * time.Second,
			},
		},
		{
			name: "Positive case #2",
			envs: []string{"ADDRESS=example.com:8181"},
			want: map[string]interface{}{"serverAddr": "example.com:8181"},
		},
		{
			name: "Positive case #3",
			envs: []string{"POLL_INTERVAL=3"},
			want: map[string]interface{}{"pollInterval": 3 * time.Second},
		},
		{
			name: "Positive case #4",
			envs: []string{"REPORT_INTERVAL=7"},
			want: map[string]interface{}{"reportInterval": 7 * time.Second},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.envs != nil && len(tc.envs) > 0 {
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
			config = parseEnvs(config)

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
		args []string
		envs []string
		want map[string]interface{}
	}{
		{
			name: "Positive case #1",
			args: nil,
			envs: nil,
			want: map[string]interface{}{
				"serverAddr":     "localhost:8080",
				"pollInterval":   2 * time.Second,
				"reportInterval": 10 * time.Second,
			},
		},
		{
			name: "Positive case #2A",
			args: []string{"-a=aaa.com:3333"},
			envs: []string{"ADDRESS=bbb.com:5555"},
			want: map[string]interface{}{"serverAddr": "bbb.com:5555"},
		},
		{
			name: "Positive case #2B",
			args: []string{"-a=aaa.com:3333"},
			envs: nil,
			want: map[string]interface{}{"serverAddr": "aaa.com:3333"},
		},
		{
			name: "Positive case #2C",
			args: nil,
			envs: []string{"ADDRESS=bbb.com:5555"},
			want: map[string]interface{}{"serverAddr": "bbb.com:5555"},
		},
		{
			name: "Positive case #3A",
			args: []string{"-p=21"},
			envs: []string{"POLL_INTERVAL=31"},
			want: map[string]interface{}{"pollInterval": 31 * time.Second},
		},
		{
			name: "Positive case #3B",
			args: []string{"-p=21"},
			envs: nil,
			want: map[string]interface{}{"pollInterval": 21 * time.Second},
		},
		{
			name: "Positive case #3C",
			args: nil,
			envs: []string{"POLL_INTERVAL=31"},
			want: map[string]interface{}{"pollInterval": 31 * time.Second},
		},
		{
			name: "Positive case #4A",
			args: []string{"-r=51"},
			envs: []string{"REPORT_INTERVAL=71"},
			want: map[string]interface{}{"reportInterval": 71 * time.Second},
		},
		{
			name: "Positive case #4B",
			args: []string{"-r=51"},
			envs: nil,
			want: map[string]interface{}{"reportInterval": 51 * time.Second},
		},
		{
			name: "Positive case #4C",
			args: nil,
			envs: []string{"REPORT_INTERVAL=71"},
			want: map[string]interface{}{"reportInterval": 71 * time.Second},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.args != nil && len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			if tc.envs != nil && len(tc.envs) > 0 {
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
				suite.Assert().Truef(field.Equal(reflect.ValueOf(want)), "Invalid value for '%s'", k)
			}
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
