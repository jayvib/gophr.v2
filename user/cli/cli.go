package cli

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/jayvib/golog"
  "github.com/spf13/cobra"
  "gophr.v2/user"
  "gophr.v2/user/service/proxy/remote"
  "log"
)

var userService user.Service

func init() {
  UserCmd.AddCommand(get)
  client, err := remote.NewClient()
  if err != nil {
    golog.Fatal(err)
  }
  userService = remote.New(client)
}

var UserCmd = &cobra.Command{
  Use: "user",
  Short: "A subcommand for interact with user service",
}

var get = &cobra.Command{
  Use: "get",
  Short: "get is a command for getting user by user ids",
  Long: `
DESCRIPTION:
  get is a command for getting user by user ids.

EXAMPLE:
  gophr user get id1 id2 id3
`,
  Run: func(cmd *cobra.Command, args[]string) {
    golog.Debug("OP: user.Get")
    golog.Debug("ARGS:", args)

    var results []*user.User

    for _, id := range args {
      usr, err := userService.GetByUserID(context.Background(), id)
      if err != nil {
        if err == user.ErrNotFound {
          fmt.Printf("User with id '%s' not exists", id)
          return
        }
        fmt.Println(err)
        return
      }
      results = append(results, usr)
    }

    payload, err := json.MarshalIndent(results, "", "  ")
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(string(payload))
  },
}

