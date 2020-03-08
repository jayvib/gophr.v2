//+build integration,mysql

package driver_test

import (
	"github.com/stretchr/testify/require"
	"gophr.v2/config"
	"gophr.v2/user/repository/mysql/driver"
	"testing"
)

func TestInitializeDriver(t *testing.T) {
	conf, err := config.New(config.DevelopmentEnv)
	require.NoError(t, err)
	db, err := driver.InitializeDriver(conf)
	defer func() {
		e := db.Close()
		require.NoError(t, e)
	}()
	require.NoError(t, err)

}
