package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: merge-coverprofiles <output-file> <input-files...>")
		os.Exit(1)
	}

	output := os.Args[1]
	inputFiles := os.Args[2:]

	combined := make(map[string]struct{})
	var header string

	for _, file := range inputFiles {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", file, err)
			os.Exit(1)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "mode:") {
				if header == "" {
					header = line
				} else if header != line {
					fmt.Println("Error: inconsistent coverage modes")
					os.Exit(1)
				}
			} else {
				combined[line] = struct{}{}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			os.Exit(1)
		}
	}

	outFile, err := os.Create(output)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	fmt.Fprintln(writer, header)
	for line := range combined {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}
