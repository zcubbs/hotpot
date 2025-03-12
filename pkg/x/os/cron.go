package os

import (
	"fmt"
	"os/exec"
	"strings"
)

// AddCronJob adds a new cron job to the current user's crontab.
// The cronJob parameter should be the full string of the cron job, e.g., "* * * * * /path/to/script.sh".
func AddCronJob(cronJob string) error {
	// Get the current crontab
	currentCron, _ := exec.Command("crontab", "-l").CombinedOutput()
	cronOutput := string(currentCron)

	var updatedCron string

	// If there's no existing crontab, just set the new one
	if strings.Contains(cronOutput, "no crontab for") {
		updatedCron = cronJob + "\n"
	} else {
		// Check if the cron job already exists to avoid duplicates
		if strings.Contains(cronOutput, cronJob) {
			return fmt.Errorf("cron job already exists")
		}
		// Append the new cron job
		updatedCron = cronOutput + cronJob + "\n"
	}

	fmt.Println("Updated Cron Content:")
	fmt.Println(updatedCron)

	// Save the updated crontab
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(updatedCron)

	combinedOutput, execErr := cmd.CombinedOutput()
	if execErr != nil {
		return fmt.Errorf("failed to add cron job. Error: %v. Output: %s", execErr, combinedOutput)
	}

	return nil
}
