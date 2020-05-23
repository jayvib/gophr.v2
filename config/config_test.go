//+build unit

package config

import (
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Clone(t *testing.T) {
		conf := &Config{
			Gophr: Gophr{
				Port: "8080",
				Environment: "ENV",
				Debug: false,
			},
			MySQL: MySQL{
				User: "pitchy",
				Password: "pitchylovespapa",
				Host: "localhost",
				Port: "3306",
				Database: "pitchy_db",
			},
		}

		clonedConf, err := conf.Clone()
		require.NoError(t, err)

		if clonedConf == conf {
			t.Error("cloned config should not the same address with the original config")
		}

		t.Run("Modifying the value", func(t *testing.T){
			clonedConf.Gophr.Port = "8081"
			assert.NotEqual(t, clonedConf, conf)
		})
}
