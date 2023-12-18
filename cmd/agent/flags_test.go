package main

import (
	"flag"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
	"time"
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
		want Config
	}{
		{
			name: "Positive case #1",
			args: nil,
			want: Config{serverAddr: "localhost:8080", pollInterval: 2 * time.Second, reportInterval: 10 * time.Second},
		},
		{
			name: "Positive case #2",
			args: []string{"-a=example.com:8181"},
			want: Config{serverAddr: "example.com:8181", pollInterval: 2 * time.Second, reportInterval: 10 * time.Second},
		},
		{
			name: "Positive case #3",
			args: []string{"-p=3"},
			want: Config{serverAddr: "localhost:8080", pollInterval: 3 * time.Second, reportInterval: 10 * time.Second},
		},
		{
			name: "Positive case #4",
			args: []string{"-r=7"},
			want: Config{serverAddr: "localhost:8080", pollInterval: 2 * time.Second, reportInterval: 7 * time.Second},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.args != nil && len(tc.args) > 0 {
				os.Args = append(os.Args, tc.args...)
			}

			config := NewConfig()
			config = parseFlags(config)

			suite.Assert().EqualValues(tc.want, config)
		})
	}
}

func (suite *FlagsTestSuite) TestParseEnvs() {
	testCases := []struct {
		name string
		envs []string
		want Config
	}{
		{
			name: "Positive case #1",
			envs: nil,
			want: Config{serverAddr: "", pollInterval: 0, reportInterval: 0},
		},
		{
			name: "Positive case #2",
			envs: []string{"ADDRESS=example.com:8181"},
			want: Config{serverAddr: "example.com:8181", pollInterval: 0, reportInterval: 0},
		},
		{
			name: "Positive case #3",
			envs: []string{"POLL_INTERVAL=3"},
			want: Config{serverAddr: "", pollInterval: 3 * time.Second, reportInterval: 0},
		},
		{
			name: "Positive case #4",
			envs: []string{"REPORT_INTERVAL=7"},
			want: Config{serverAddr: "", pollInterval: 0, reportInterval: 7 * time.Second},
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

			suite.Assert().EqualValues(tc.want, config)
		})
	}
}

func (suite *FlagsTestSuite) TestLoadConfig() {
	testCases := []struct {
		name string
		args []string
		envs []string
		want Config
	}{
		{
			name: "Positive case #1",
			args: nil,
			envs: nil,
			want: Config{serverAddr: "localhost:8080", pollInterval: 2 * time.Second, reportInterval: 10 * time.Second},
		},
		{
			name: "Positive case #2",
			args: []string{"-a=aaa.com:3333"},
			envs: []string{"ADDRESS=bbb.com:5555"},
			want: Config{serverAddr: "bbb.com:5555", pollInterval: 2 * time.Second, reportInterval: 10 * time.Second},
		},
		{
			name: "Positive case #3",
			args: []string{"-p=21"},
			envs: []string{"POLL_INTERVAL=31"},
			want: Config{serverAddr: "localhost:8080", pollInterval: 31 * time.Second, reportInterval: 10 * time.Second},
		},
		{
			name: "Positive case #4",
			args: []string{"-r=51"},
			envs: []string{"REPORT_INTERVAL=71"},
			want: Config{serverAddr: "localhost:8080", pollInterval: 2 * time.Second, reportInterval: 71 * time.Second},
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

			suite.Assert().EqualValues(tc.want, config)
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
