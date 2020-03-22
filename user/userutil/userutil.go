package userutil

import "gophr.v2/util/randutil"

func GenerateID() string {
	return randutil.GenerateID("user")
}
