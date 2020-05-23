package cli

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/jayvib/golog"
  "github.com/spf13/cobra"
  "gophr.v2/user"
  "gophr.v2/user/service"
  "gophr.v2/user/service/proxy/remote"
  "log"
  "os"
)

var userService user.Service

var defaultOut = os.Stdout

var writeToFilePath string

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

  UserCmd.PersistentFlags().StringVar(&writeToFilePath, "to-file","", "Write result to file")
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

    usrs, err := service.GetByUserIDs(context.Background(), userService, args...)
    if err != nil {
      fmt.Println(err)
    }

    payload, err := json.MarshalIndent(usrs, "", "  ")
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(string(payload))

    golog.Debug("filepath", writeToFilePath)
    switch {
    case writeToFilePath != "":
      f, err := os.Create(writeToFilePath)
      if err != nil {
        log.Fatal(err)
      }
      defer f.Close()
      _, err = f.Write(payload)
      if err != nil {
        log.Fatal(err)
      }
    default:
      fmt.Println(string(payload))
    }
  },
}

