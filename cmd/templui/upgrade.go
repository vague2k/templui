package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// runUpgrade handles the 'upgrade' command logic.
func runUpgrade(args []string, commandArg string) {
	var ref string

	if strings.Contains(commandArg, "@") {
		parts := strings.SplitN(commandArg, "@", 2)
		if len(parts) == 2 && parts[0] == "upgrade" && parts[1] != "" {
			ref = parts[1]
			fmt.Printf("Updating templUI cli using specified ref: %s\n", ref)
		} else {
			fmt.Printf("Error: Invalid format '%s'. Use 'upgrade' or 'upgrade@<ref>'.\n", commandArg)
			return
		}
	}

	if err := updateCLI(ref); err != nil {
		fmt.Printf("Error upgrading templUI cli: %v\n", err)
	}
}

// updateCLI attempts to install a new version of the templUI cli based on the passed in ref.
func updateCLI(ref string) error {
	if ref == "" {
		ref = "latest"
	}
	cmd := exec.Command("go", "install", fmt.Sprintf("github.com/templui/templui/cmd/templui@%s", ref))
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	fmt.Printf("Updated templUI to ref '%s'\n", ref)
	return nil
}
