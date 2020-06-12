package main

import "fmt"

func main() {
	c := 'r'
	n := c | 0x20
	fmt.Println(string(n))
}
