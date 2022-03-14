package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	lines := scan("rules")
	generate(lines)
}

func generate(lines []string) {
	table := make([]string, 0)

	table = append(table, "[AutoProxy 0.2.9]")
	table = append(table, fmt.Sprintf("! Last Modified: %s", time.Now()))
	table = append(table, "! Expires: 24h")
	table = append(table, "! HomePage: https://github.com/mygfw/mygfw\n\n")

	f, err := os.OpenFile("gfwlist.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	encoder := base64.NewEncoder(base64.StdEncoding, f)
	if _, err = encoder.Write([]byte(strings.Join(table, "\n"))); err != nil {
		panic(err)
	}

	if _, err = encoder.Write([]byte(strings.Join(lines, "\n"))); err != nil {
		panic(err)
	}

	if err = encoder.Close(); err != nil {
		panic(err)
	}

	fmt.Println("done")
}

func scan(dir string) []string {
	lines := make([]string, 0)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		lines = append(lines, load(path)...)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return lines
}

func load(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) < 1 {
			continue
		}

		if line[0:1] == "!" || line[0:1] == "#" {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}
