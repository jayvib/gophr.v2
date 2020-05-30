// +Build integration

package configutil_test

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) (teardown func() error) {
	home := os.Getenv("HOME")
	gophrPath := filepath.Join(home, ".gophr", "testenv")
	err := os.MkdirAll(gophrPath, 0777)
	require.NoError(t, err)

	confText := `
mysql:
  user: root
  password: test
  host: 127.0.0.1
  port: 6607
  name: user
`
	configPath := filepath.Join(gophrPath, "config.yaml")
	err = ioutil.WriteFile(configPath, []byte(confText), 0777)
	require.NoError(t, err)

	return func() error {
		return os.Remove(configPath)
	}
}
