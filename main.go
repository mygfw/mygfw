package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	lines := scan("rules")
	gfw(lines)
	clash(lines)
}

func clash(lines []string) {
	table := make([]string, 0)

	table = append(table, "payload:")

	f, err := os.OpenFile("clash.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	mark := make(map[string]struct{}, 0)
	for _, v := range lines {
		h1 := v[0:1]
		h2 := v[0:2]
		e1 := v[len(v)-1:]

		if h1 == "\\" { // 正则表达式
			continue
		} else if h2 == "@@" || h1 == "!" { // 例外规则 => @@, 注释规则 => !
			continue
		} else if h2 == "||" { // 全匹配规则 => ||
			v = "." + v[2:]
		} else if v[0:2] == "*." { // 通配符支持 => *
			v = "." + v[2:]
		} else if h1 == "|" || e1 == "|" { // 匹配地址开始和结尾规则 => |
			v = strings.Trim(v, "|")
			parse, err := url.Parse(v)
			if err == nil && parse != nil {
				if parse.Host != "" {
					v = "." + parse.Host
				} else {
					v = "." + v
				}
			}
		}

		h1 = v[0:1]
		if h1 == "." {
			v = fmt.Sprintf("  - '+%s'", v)
		} else {
			v = fmt.Sprintf("  - %s", v)
		}

		if _, ok := mark[v]; ok || v == "" {
			continue
		}

		mark[v] = struct{}{}
		table = append(table, v)
	}

	if _, err = f.WriteString(strings.Join(table, "\n")); err != nil {
		panic(err)
	}

	if err = f.Close(); err != nil {
		panic(err)
	}

	fmt.Println("clash done")
}

func gfw(lines []string) {
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

	fmt.Println("gfwlist done")
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
