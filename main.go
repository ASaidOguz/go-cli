package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Process struct {
	Name           string
	ProcessId      int
	WorkingSetSize int
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "memtop",
		Short: "A CLI to list the most memory-consuming services",
		Long: `This CLI demonstrates how to create a CLI application with Cobra
that lists the most memory-consuming services on Windows systems.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			folderName := args[0]
			fileName := args[1]

			err := os.MkdirAll(folderName, 0755)
			if err != nil {
				fmt.Printf("Error creating folder: %v\n", err)
				os.Exit(1)
			}

			outputPath := fmt.Sprintf("%s/%s.txt", folderName, fileName)

			psCmd := exec.Command("tasklist", "/v")
			output, err := psCmd.Output()
			if err != nil {
				fmt.Printf("Error executing tasklist command: %v\n", err)
				os.Exit(1)
			}

			processes := parsePsOutput(string(output))
			sort.Slice(processes, func(i, j int) bool {
				return processes[i].WorkingSetSize > processes[j].WorkingSetSize
			})

			file, err := os.Create(outputPath)
			if err != nil {
				fmt.Printf("Error creating file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()

			file.WriteString("PID\t\tCommand\t\tRSS (KB)\n")
			for _, p := range processes {
				line := fmt.Sprintf("%-10d\t%-30s\t%d\n", p.ProcessId, p.Name, p.WorkingSetSize)
				file.WriteString(line)
			}

			fmt.Printf("Process information written to: %s\n", outputPath)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parsePsOutput(output string) []Process {
	lines := strings.Split(output, "\n")
	processes := make([]Process, 0, len(lines)-1)

	for i, line := range lines {
		if i == 0 {
			continue // Skip header line
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		ProcessId, _ := strconv.Atoi(fields[1])
		WorkingSetSize, _ := strconv.Atoi(strings.Replace(fields[4], ",", "", -1))

		processes = append(processes, Process{
			ProcessId:      ProcessId,
			Name:           fields[0],
			WorkingSetSize: WorkingSetSize,
		})
	}

	return processes
}
