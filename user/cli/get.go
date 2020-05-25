package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jayvib/golog"
	"github.com/spf13/cobra"
	usersvc "gophr.v2/user/service"
	"log"
	"os"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get is a command for getting user by user ids",
	Long: `
DESCRIPTION:
  get is a command for getting user by user ids.

EXAMPLE:
  gophr user get id1 id2 id3
`,

	Run: func(cmd *cobra.Command, args []string) {

		usrs, err := usersvc.GetByUserIDs(context.Background(), userService, args...)
		if err != nil {
			fmt.Println(err)
		}

		if len(usrs) == 0 {
			return
		}

		payload, err := json.MarshalIndent(usrs, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

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
