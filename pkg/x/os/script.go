package os

import (
	"fmt"
	"os"
	"strings"
)

// GenerateBashScript creates a bash script file with the given commands.
// It returns the path of the created script or an error.
func GenerateBashScript(outputPath string, commands ...string) error {
	// Start with the bash shebang
	scriptContent := "#!/bin/bash\n\n"

	// Append commands
	scriptContent += strings.Join(commands, "\n")

	// Write the script content to the specified output path
	// #nosec G306 -- Script file needs execute permissions to be runnable
	err := os.WriteFile(outputPath, []byte(scriptContent), 0700) // Give execute permission
	if err != nil {
		return fmt.Errorf("failed to write bash script path=%s error=%w", outputPath, err)
	}

	return nil
}
