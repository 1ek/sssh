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
	hostRegex := regexp.MustCompile(`(?m)^Host\s+([\w.-]+)$`)
	// Find all matches of host names in the config file
	matches := hostRegex.FindAllStringSubmatch(string(content), -1)
	// Extract host names from the matches
	for _, match := range matches {
		hosts = append(hosts, match[1])
	}
	return hosts, nil
}

func main() {
	var base *huh.Theme = huh.ThemeBase()
	var base16 *huh.Theme = huh.ThemeBase16()
	var dracula *huh.Theme = huh.ThemeDracula()
	var charm *huh.Theme = huh.ThemeCharm()
	var catppuccin *huh.Theme = huh.ThemeCatppuccin()

	themeEnv := os.Getenv("SSSH_THEME")
	var selectedTheme *huh.Theme
	switch themeEnv {
	case "base":
		selectedTheme = base
	case "base16":
		selectedTheme = base16
	case "dracula":
		selectedTheme = dracula
	case "charm":
		selectedTheme = charm
	case "catppuccin":
		selectedTheme = catppuccin
	default:
		selectedTheme = dracula
	}
	// Retrieve the list of hosts from the SSH config file
	hosts, err := getHosts()
	if err != nil {
		fmt.Println("Error:", err)
	}

	var height int
	l := len(hosts)
	switch {
	case l > 10:
		height = 15
	default:
		height = l + 3
	}

	// Prompt the user to select a host
	var host string
	s := huh.NewSelect[string]().
		Title("Connect to: ").
		Options(huh.NewOptions(hosts...)...).
		Value(&host)
	huh.NewForm(huh.NewGroup(s)).WithHeight(height).WithTheme(selectedTheme).Run()

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
