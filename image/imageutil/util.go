package imageutil

import "gophr.v2/util/randutil"

func GenerateID() string {
	return randutil.GenerateID("image")
}
