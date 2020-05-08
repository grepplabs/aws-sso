package cmd

import (
	"fmt"
	"os"

	"github.com/grepplabs/aws-sso/pkg/credentials"
	"github.com/spf13/cobra"
)

var credentialsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Get SSO credentials and print AWS environment variables to set",
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

		fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", roleCredentials.AccessKeyId)
		fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", roleCredentials.SecretAccessKey)
		fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", roleCredentials.SessionToken)
	},
}

func init() {
	credentialsCmd.AddCommand(credentialsExportCmd)
}
