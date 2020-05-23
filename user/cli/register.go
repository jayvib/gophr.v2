package cli

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/jayvib/golog"
  "github.com/spf13/cobra"
  "golang.org/x/crypto/ssh/terminal"
  "gophr.v2/user"
  "log"
  "syscall"
)

var registerCmd = &cobra.Command{
  Use: "register",
  Short: "Is a sub-command for registering a new user",
  Run: func(cmd *cobra.Command, args []string){
    // This will ask user information
    var username string
    fmt.Print("Enter username: ")
    _, err := fmt.Scanf("%s", &username)
    if err != nil {
      log.Fatal(err)
    }

    var email string
    fmt.Print("Enter email: ")
    _, err = fmt.Scanf("%s", &email)
    if err != nil {
      log.Fatal(err)
    }

    var password string
    fmt.Print("Enter Password: ")
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
      log.Fatal(err)
    }
    password = string(bytePassword)
    golog.Debug("Username:", username)
    golog.Debug("Email:", email)
    golog.Debug("Password:", password)

    usr := &user.User{
      Username: username,
      Email: email,
      Password: password,
    }
    err = userService.Register(context.Background(), usr)
    if err != nil {
      log.Fatal(err)
    }

    payload, err := json.MarshalIndent(usr, "", "  ")
    if err != nil {
      log.Fatal(err)
    }

    fmt.Print(string(payload))
  },
}