package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"harvest-cli/internal/models"
	"harvest-cli/internal/prompt"
)

var (
	updateID    int
	updateHours string
	updateNotes string
	updateDate  string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing time entry",
	Long: `Update hours, notes, or date on an existing time entry.

Use 'harvest view' to find the entry ID. Only flags you pass are changed.

Examples:
  harvest update --id 12345 --hours 0.5
  harvest update --id 12345 --hours 1.25 --notes "Standup"
  harvest update --id 12345 --date 2026-04-08`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().IntVar(&updateID, "id", 0, "Time entry ID (required)")
	updateCmd.Flags().StringVar(&updateHours, "hours", "", "New hours (e.g. 1.5, 1:30, 90m)")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "New notes")
	updateCmd.Flags().StringVarP(&updateDate, "date", "d", "", "New date (YYYY-MM-DD)")
	_ = updateCmd.MarkFlagRequired("id")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if updateHours == "" && !cmd.Flags().Changed("notes") && updateDate == "" {
		return fmt.Errorf("at least one of --hours, --notes, or --date must be provided")
	}

	req := &models.UpdateTimeEntryRequest{}

	if updateHours != "" {
		h, err := prompt.ParseHours(updateHours)
		if err != nil {
			return fmt.Errorf("invalid hours value: %s", updateHours)
		}
		req.Hours = &h
	}

	if cmd.Flags().Changed("notes") {
		n := updateNotes
		req.Notes = &n
	}

	if updateDate != "" {
		if _, err := time.Parse("2006-01-02", updateDate); err != nil {
			return fmt.Errorf("invalid --date format, use YYYY-MM-DD: %s", updateDate)
		}
		d := updateDate
		req.SpentDate = &d
	}

	fmt.Printf("Updating time entry %d...\n", updateID)
	entry, err := apiClient.UpdateTimeEntry(updateID, req)
	if err != nil {
		return fmt.Errorf("failed to update time entry: %w", err)
	}

	fmt.Printf("\nTime entry updated successfully!\n")
	fmt.Printf("  ID:      %d\n", entry.ID)
	fmt.Printf("  Project: %s\n", entry.Project.Name)
	fmt.Printf("  Task:    %s\n", entry.Task.Name)
	fmt.Printf("  Date:    %s\n", entry.SpentDate)
	fmt.Printf("  Hours:   %.2f\n", entry.Hours)
	if entry.Notes != "" {
		fmt.Printf("  Notes:   %s\n", entry.Notes)
	}

	return nil
}
