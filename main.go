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
	gfw(lines)
	clash(lines)
	rocket(lines)
}

func rocket(lines []string) {
	f, err := os.OpenFile("rocket.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	table := make([]string, 0)
	table = append(table, "[General]")
	table = append(table, "ipv6 = false")
	table = append(table, "bypass-system = true")
	table = append(table, "skip-proxy = 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, localhost, *.local, e.crashlytics.com, captive.apple.com")
	table = append(table, "bypass-tun = 10.0.0.0/8,100.64.0.0/10,127.0.0.0/8,169.254.0.0/16,172.16.0.0/12,192.0.0.0/24,192.0.2.0/24,192.88.99.0/24,192.168.0.0/16,198.18.0.0/15,198.51.100.0/24,203.0.113.0/24,224.0.0.0/4,255.255.255.255/32")
	table = append(table, "dns-server = system, 114.114.114.114, 112.124.47.27, 8.8.8.8, 8.8.4.4")
	table = append(table, "[Rule]\n")

	mark := make(map[string]struct{}, 0)
	for _, v := range lines {
		h1 := v[0:1]
		h2 := v[0:2]
		if h1 == "!" { // 注释
			continue
		} else if h1 == "." { // 域名前缀
			v = fmt.Sprintf("DOMAIN-SUFFIX,%s,PROXY", v[1:])
		} else if h2 == "ip" { // ip范围
			v = fmt.Sprintf("IP-CIDR,%s,PROXY", v[3:])
		} else { // 固定域名
			v = fmt.Sprintf("DOMAIN,%s,PROXY", v)
		}

		if _, ok := mark[v]; ok || v == "" {
			continue
		}

		mark[v] = struct{}{}
		table = append(table, v)
	}

	table = append(table, "FINAL,DIRECT\n")
	table = append(table, "[URL Rewrite]")
	table = append(table, "^http://(www.)?google.cn https://www.google.com 302")

	if _, err = f.WriteString(strings.Join(table, "\n")); err != nil {
		panic(err)
	}

	if err = f.Close(); err != nil {
		panic(err)
	}

	fmt.Println("rocket done")
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
		if h1 == "!" { // 注释
			continue
		} else if h1 == "." { // 域名前缀
			v = " - '+." + v[1:] + "'"
		} else if h2 == "ip" { // ip范围
			v = " - '" + v[3:] + "'"
		} else { // 固定域名
			v = " - " + v[1:]
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
	table = append(table, "! HomePage: https://github.com/mygfw/mygfw\n")

	f, err := os.OpenFile("gfwlist.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	mark := make(map[string]struct{}, 0)
	for _, v := range lines {
		h1 := v[0:1]
		h2 := v[0:2]
		if h1 == "!" { // 注释
			continue
		} else if h1 == "." { // 域名前缀
			v = "||" + v[1:]
		} else if h2 == "ip" { // ip范围
			continue
		}

		if _, ok := mark[v]; ok || v == "" {
			continue
		}

		mark[v] = struct{}{}
		table = append(table, v)
	}

	encoder := base64.NewEncoder(base64.StdEncoding, f)
	if _, err = encoder.Write([]byte(strings.Join(table, "\n"))); err != nil {
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
