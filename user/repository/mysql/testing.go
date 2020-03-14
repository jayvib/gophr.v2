package mysql

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func CheckErr(t *testing.T, err error) {
	checkErr(t, err)
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}
