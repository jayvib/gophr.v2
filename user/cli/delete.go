package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A sub-command for deleting a user",
	Run: func(cmd *cobra.Command, args []string) {
		for _, id := range args {
			err := userService.Delete(context.Background(), id)
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}
