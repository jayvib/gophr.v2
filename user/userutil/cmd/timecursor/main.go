package main

import (
	"fmt"
	"gophr.v2/user/userutil"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

func main() {
		t := time.Now().AddDate(0, 0, -1)
		cursor := userutil.EncodeCursor(t)
		fmt.Println(cursor)

}
