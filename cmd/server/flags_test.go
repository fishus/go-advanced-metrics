package main

import (
	"flag"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"strings"
	"testing"
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
	for _, e := range []string{"ADDRESS"} {
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
		want map[string]string
	}{
		{
			name: "Positive case #1",
			args: nil,
			want: map[string]string{"serverAddr": "localhost:8080"},
		},
		{
			name: "Positive case #2",
			args: []string{"-a=example.com:8181"},
			want: map[string]string{"serverAddr": "example.com:8181"},
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
				suite.Assert().EqualValues(want, configFields.FieldByName(k).String())
			}
		})
	}
}

func (suite *FlagsTestSuite) TestParseEnvs() {
	testCases := []struct {
		name string
		envs []string
		want map[string]string
	}{
		{
			name: "Positive case #1",
			envs: nil,
			want: map[string]string{"serverAddr": ""},
		},
		{
			name: "Positive case #2",
			envs: []string{"ADDRESS=example.com:8181"},
			want: map[string]string{"serverAddr": "example.com:8181"},
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
				suite.Assert().EqualValues(want, configFields.FieldByName(k).String())
			}
		})
	}
}

func (suite *FlagsTestSuite) TestLoadConfig() {
	testCases := []struct {
		name string
		args []string
		envs []string
		want map[string]string
	}{
		{
			name: "Positive case #1",
			args: nil,
			envs: nil,
			want: map[string]string{"serverAddr": "localhost:8080"},
		},
		{
			name: "Positive case #2",
			args: []string{"-a=aaa.com:3333"},
			envs: []string{"ADDRESS=bbb.com:5555"},
			want: map[string]string{"serverAddr": "bbb.com:5555"},
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
				suite.Assert().EqualValues(want, configFields.FieldByName(k).String())
			}
		})
	}
}

func TestFlagsSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
