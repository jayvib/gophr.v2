//+build integration,mysql

package mysql_test

import (
	"github.com/jayvib/golog"
	"github.com/stretchr/testify/require"
	"gophr.v2/config"
	"gophr.v2/user/repository/mysql/driver"
	"testing"
)

func TestInitializeDriver(t *testing.T) {
	golog.SetLevel(golog.DebugLevel)
	conf, err := config.New(config.DevelopmentEnv)
	require.NoError(t, err)
	db, err := driver.InitializeDriver(conf)
	defer func() {
		e := db.Close()
		require.NoError(t, e)
	}()
	require.NoError(t, err)

}
