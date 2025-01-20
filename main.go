package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/term"
)

const (
	greenBgWhiteText = "\033[1;42;37m"
	resetColor       = "\033[0m"
)

func main() {
	// Command-line flags
	signal := flag.String("s", "", "Signal to send (e.g., -s 9 for SIGKILL)")
	yes := flag.Bool("y", false, "Assume yes; kill all matching processes without confirmation")
	flag.Parse()

	// Handle signals like -9
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: ka [options] process_name")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Extract process name and additional signals
	var processName string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Handle signals provided without -s flag
			if num, err := strconv.Atoi(strings.TrimPrefix(arg, "-")); err == nil {
				*signal = strconv.Itoa(num)
			} else {
				log.Fatalf("Invalid signal: %s", arg)
			}
		} else {
			processName = arg
		}
	}

	if processName == "" {
		log.Fatal("Process name is required")
	}

	// Default signal is SIGTERM (15)
	if *signal == "" {
		*signal = "15"
	}

	// Get the current process ID to exclude it later
	currentPID := os.Getpid()

	// Use pgrep to find matching PIDs
	pgrepCmd := exec.Command("pgrep", "-f", processName)
	var pgrepOut bytes.Buffer
	pgrepCmd.Stdout = &pgrepOut
	if err := pgrepCmd.Run(); err != nil {
		fmt.Printf("No processes found matching '%s'\n", processName)
		os.Exit(0)
	}

	// Parse PIDs
	pidStrings := strings.Fields(pgrepOut.String())
	var pids []int
	for _, pidStr := range pidStrings {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}
		// Exclude the current process
		if pid == currentPID {
			continue
		}
		pids = append(pids, pid)
	}

	if len(pids) == 0 {
		fmt.Printf("No processes found matching '%s'\n", processName)
		os.Exit(0)
	}

	// If -y flag is provided, kill all matching processes without confirmation
	if *yes {
		for _, pid := range pids {
			err := exec.Command("kill", "-"+*signal, strconv.Itoa(pid)).Run()
			if err != nil {
				fmt.Printf("Failed to kill process %d: %v\n", pid, err)
			} else {
				fmt.Printf("Killed process %d\n", pid)
			}
		}
		return
	}

	// If only one matching process, kill it without interactive dialog
	if len(pids) == 1 {
		pid := pids[0]
		err := exec.Command("kill", "-"+*signal, strconv.Itoa(pid)).Run()
		if err != nil {
			fmt.Printf("Failed to kill process %d: %v\n", pid, err)
		} else {
			fmt.Printf("Killed process %d\n", pid)
		}
		return
	}

	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Default to 80x24 if terminal size cannot be determined
		width = 80
		height = 24
	}

	// Adjust PageSize to use full terminal height before scrolling
	pageSize := height - 4 // Subtract for prompt and padding
	if pageSize < 1 {
		pageSize = 1
	}

	// Prepare options for interactive selection
	pidMap := make(map[string]int)
	var options []string

	// Use ps to get command lines for the PIDs
	psArgs := append([]string{"-o", "pid=,comm=,args=", "-p"}, pidStrings...)
	psCmd := exec.Command("ps", psArgs...)
	var psOut bytes.Buffer
	psCmd.Stdout = &psOut
	if err := psCmd.Run(); err != nil {
		log.Fatalf("Failed to get process information: %v", err)
	}

	scanner := bufio.NewScanner(&psOut)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		pidStr, name, cmdline := fields[0], fields[1], strings.Join(fields[2:], " ")

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Sanitize name and cmdline to remove newlines
		name = sanitizeString(name)
		cmdline = sanitizeString(cmdline)
		optionStr := formatOptionWithHighlight(pid, name, cmdline, width, processName)
		options = append(options, optionStr)
		pidMap[optionStr] = pid
	}

	if len(options) == 0 {
		fmt.Printf("No processes found matching '%s'\n", processName)
		os.Exit(0)
	}

	// Interactive selection
	selectedOptions := []string{}
	prompt := &survey.MultiSelect{
		Message:  "Select processes to kill:",
		Options:  options,
		Default:  options,
		PageSize: pageSize,
	}
	err = survey.AskOne(prompt, &selectedOptions)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Kill selected processes
	for _, option := range selectedOptions {
		pid := pidMap[option]
		err := exec.Command("kill", "-"+*signal, strconv.Itoa(pid)).Run()
		if err != nil {
			fmt.Printf("Failed to kill process %d: %v\n", pid, err)
		} else {
			fmt.Printf("Killed process %d\n", pid)
		}
	}
}

func formatOptionWithHighlight(pid int, name, cmdline string, width int, processName string) string {
	pidWidth := 8
	nameWidth := 25
	cmdWidth := width - pidWidth - nameWidth - 11
	if cmdWidth < 10 {
		cmdWidth = 10
	}

	name = truncateString(name, nameWidth)
	cmdline = truncateString(cmdline, cmdWidth)

	name = highlightText(name, processName)
	cmdline = highlightText(cmdline, processName)

	optionStr := fmt.Sprintf("%-*d  %-*s  %-*s",
		pidWidth, pid,
		nameWidth, name,
		cmdWidth, cmdline)

	return sanitizeString(optionStr)
}

func highlightText(text, search string) string {
	return strings.ReplaceAll(text, search, greenBgWhiteText+search+resetColor)
}

// truncateString truncates a string to a specified width, adding "..." if truncated
func truncateString(s string, maxWidth int) string {
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth > 3 {
		return string(runes[:maxWidth-3]) + ".."
	}
	return string(runes[:maxWidth])
}

// sanitizeString removes any newline characters from a string
func sanitizeString(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\r", " ")
}
