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

type getResult struct {
  usr *user.User
  err error
  id string
}

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
    var results []*user.User

    res := make(chan *getResult)
    for _, id := range args {
      go func(i string) {
        usr, err := userService.GetByUserID(context.Background(), i)
        if err != nil {
          res <-&getResult{err: err, id: i}
          return
        }
        res <- &getResult{usr: usr, id: i}
      } (id)
    }

    for i := 0; i < len(args); i++ {
      r := <-res
      if r.err == nil {
        results = append(results, r.usr)
      } else {
        if r.err == user.ErrNotFound {
          fmt.Printf("User with '%s' not exists 🙂🙂🙂🙂", r.id)
        } else {
          fmt.Println("Unexpected error:", r.err)
        }
      }
    }

    if results == nil {
      return
    }

    payload, err := json.MarshalIndent(results, "", "  ")
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(string(payload))
  },
}

