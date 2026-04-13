package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a time entry",
	Long: `Delete a time entry by ID.

Use 'harvest view' to find the entry ID.

Example:
  harvest delete 123456789`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid entry ID: %s", args[0])
	}

	fmt.Printf("Deleting time entry %d...\n", id)
	if err := apiClient.DeleteTimeEntry(id); err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}

	fmt.Printf("Time entry %d deleted.\n", id)
	return nil
}
