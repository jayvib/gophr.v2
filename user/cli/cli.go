package cli

import (
  "github.com/jayvib/golog"
  "github.com/spf13/cobra"
  "gophr.v2/user"
  "gophr.v2/user/service/proxy/remote"
)

var userService user.Service


var writeToFilePath string

type getResult struct {
  usr *user.User
  err error
  id string
}

func init() {
  UserCmd.AddCommand(getCmd, registerCmd, getAllCmd)
  client, err := remote.NewClient()
  if err != nil {
    golog.Fatal(err)
  }
  userService = remote.New(client)

  UserCmd.PersistentFlags().StringVar(&writeToFilePath, "to-file","", "Write result to file")

  // Get All Flags
  getAllCmd.Flags().StringVar(&cursor, "cursor", "", "Base-64 time-encoded cursor")
  getAllCmd.Flags().IntVar(&num, "num", 1, "Number of item to fetch")
  _ = getAllCmd.MarkFlagRequired("num")

}

var UserCmd = &cobra.Command{
  Use: "user",
  Short: "A subcommand for interact with user service",
}


