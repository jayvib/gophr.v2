package main

import (
	"gophr.v2/cli"
	usercli "gophr.v2/user/cli"
	"log"
)

func main() {
	cli.GophrApp.AddCommand(usercli.UserCmd)
	if err := cli.GophrApp.Execute(); err != nil {
		log.Fatal(err)
	}
}
