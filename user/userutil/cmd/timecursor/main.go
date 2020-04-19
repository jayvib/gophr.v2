package main

import (
	"bufio"
	"fmt"
	"github.com/jayvib/golog"
	"gophr.v2/user/userutil"
	"os"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter input: ")
		scanner.Scan()
		input := scanner.Text()
		t, err := time.Parse(timeFormat, input)
		if err != nil {
			golog.Error(err)
			continue
		}
		cursor := userutil.EncodeCursor(t)
		fmt.Println(cursor)
	}
	if err := scanner.Err(); err != nil {
		golog.Fatal(err)
	}

}
