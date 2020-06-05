//+build integration

package mysql

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/config"
	"gophr.v2/config/builder/viper"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	builder := viper.New(0,
		viper.SetConfigType("yaml"),
		viper.SetConfigPath("testdata"),
		viper.SetConfigName("mysql-conf.yaml"),
	)
	conf, err := config.New(builder)
	t.Run("Single Initialization", func(t *testing.T){
		require.NoError(t, err)

		conn, err := New(conf, "gophr_integration_test")
		require.NoError(t, err)
		if assert.NotNil(t, conn) {
			conn.Close()
		}
		assert.Len(t, dbPool, 1)
	})

  t.Run("Two the same database initialization should return the same pointer address", func(t *testing.T){
		require.NoError(t, err)

		first, err := New(conf, "gophr_integration_test")
		require.NoError(t, err)
		if assert.NotNil(t, first) {
			first.Close()
		}

		second, err := New(conf, "gophr_integration_test")
		require.NoError(t, err)
		assert.Equal(t, first, second)
	})

  t.Run("Parallel same database initialization should initialize only once", func(t *testing.T){
  	blockChan := make(chan interface{})
  	var wg sync.WaitGroup
  	for i := 0; i < 10; i++ {
  		wg.Add(1)
  		go func() {
  			<-blockChan
  			defer wg.Done()
				_, err = New(conf, "gophr_integration_test")
				require.NoError(t, err)
			}()
		}
  	close(blockChan)
  	wg.Wait()
		assert.Equal(t, 1, _testInitCount)
	})
}
