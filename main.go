package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/charmbracelet/huh"
)

// getHosts retrieves the list of hosts from the SSH config file using regular expressions.
func getHosts() ([]string, error) {
	// Get the user's home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// Construct the path to the SSH config file
	sshConfigPath := homedir + "/.ssh/config"
	// Read the content of the SSH config file
	content, err := os.ReadFile(sshConfigPath)
	if err != nil {
		return nil, err
	}

	var hosts []string
	// Define a regular expression to match host names
	hostRegex := regexp.MustCompile(`(?m)^Host\s+([^\s]+)`)
	// Find all matches of host names in the config file
	matches := hostRegex.FindAllStringSubmatch(string(content), -1)
	// Extract host names from the matches
	for _, match := range matches {
		hosts = append(hosts, match[1])
	}
	return hosts, nil
}

func main() {
	// Retrieve the list of hosts from the SSH config file
	hosts, err := getHosts()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Prompt the user to select a host
	var host string
	s := huh.NewSelect[string]().
		Title("Connect to: ").
		Options(huh.NewOptions(hosts...)...).
		Value(&host)

	huh.NewForm(huh.NewGroup(s)).WithTheme(huh.ThemeDracula()).Run()

	if host == "" {
		fmt.Println("No host selected. Exiting...")
		return
	}

	// Run the SSH command to connect to the selected host
	cmd := exec.Command("ssh", host)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
