package cmd

import (
	"fmt"
	"os"

	"github.com/grepplabs/aws-sso/pkg/credentials"
	"github.com/spf13/cobra"
)

var credentialsRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh SSO credentials and add a profile to your AWS credential file ~/.aws/credentials",
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		roleCredentials, err := credentials.RetrieveRoleCredentials(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// fmt.Println(time.Unix(roleCredentials.Expiration/1000, 0).Format(time.RFC3339))

		err = credentials.RefreshProfileCredentials(profile, roleCredentials)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	credentialsCmd.AddCommand(credentialsRefreshCmd)
}
