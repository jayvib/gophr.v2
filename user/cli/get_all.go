package cli

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/jayvib/golog"
  "github.com/spf13/cobra"
  "log"
)

var cursor string
var num int

var getAllCmd = &cobra.Command{
  Use: "getall",
  Short: "A command for getting all the users",
  Run: func(cmd *cobra.Command, args []string) {

    golog.Debugf("Cursor: %s Num: %d\n", cursor, num)

    usrs, next, err := userService.GetAll(context.Background(), cursor, num)
    if err != nil {
      log.Fatal(err)
    }

    payload, err := json.MarshalIndent(usrs, "", "  ")
    if err != nil {
      log.Fatal(err)
    }

    golog.Debug("next cursor:", next)
    if next != "" {
      fmt.Println("Next Cursor:", next)
    }

    fmt.Println(string(payload))
  },
}
