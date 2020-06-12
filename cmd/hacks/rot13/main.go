package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// rot13Reader reads the bytes from r and
// apply the ROT13 algorithm to each bytes
// reads from the reader.
type rot13Reader struct {
	r io.Reader
}

func (r *rot13Reader) Read(b []byte) (n int, err error) {
	n, err = r.r.Read(b)
	if err != nil {
		return n, err
	}

	// Loop runes and convert it to lower case.
	for i := 0; i < n; i++ {
		lc := b[i] | 0x20 // lowered case
		if 'a' <= lc && lc <= 'm' {
			b[i] += 13
		} else if 'n' <= lc && lc <= 'z' {
			b[i] -= 13
		}
	}
	return n, nil
}

func main() {
	input := "ABCDEFGabcdefg"
	r := strings.NewReader(input)
	rr := &rot13Reader{r: r}

	c, err := ioutil.ReadAll(rr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(c))
}
