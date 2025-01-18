package command

import (
	"context"

	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "konbi",
		Short: "CLI to manage project secrets in .env form and stored in Konbini.",
	}

	return rootCmd.ExecuteContext(context.Background())
}
