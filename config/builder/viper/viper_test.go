package viper

import (
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "gophr.v2/config"
  "testing"
)

func TestBuilder(t *testing.T) {

  testCases := []struct{
    testName string
    configName string
    configPath string
    configType string
    want *config.Config
  }{
    {
      testName: "Basic",
      configName: "config-test",
      configType: "yaml",
      configPath: "testdata",
      want: &config.Config{
        Gophr: config.Gophr{
          Port: "8080",
          Env: "DEV",
        },
      },
    },
  }

  for _, tc := range testCases {
    t.Run(tc.testName, func(t *testing.T){
      b := New(config.StageEnv,
        SetConfigName(tc.configName),
        SetConfigPath(tc.configPath),
        SetConfigType(tc.configType),
      )
      got, err := config.Build(b)
      require.NoError(t, err)
      assert.Equal(t, tc.want, got)
    })
  }
}