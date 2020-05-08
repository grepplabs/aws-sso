package cmd

import (
	"github.com/spf13/cobra"
)

var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "AWS credentials utilities",
}

func init() {
	rootCmd.AddCommand(credentialsCmd)

	credentialsCmd.PersistentFlags().String("profile", "", "Use a specific profile from your credential file.")
}
