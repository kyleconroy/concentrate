package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type HostsFile struct {
	Entries []string
}

// Add a new hostfile entry
func (hostsfile *HostsFile) Add(domain string) {
	for _, entry := range hostsfile.Entries {
		if entry == domain {
			return
		}
	}

	// TODO: Add domain parsing logic
	hostsfile.Entries = append(hostsfile.Entries, domain)
}

// Remove an existing hostfile entry
func (hostsfile *HostsFile) Remove(domain string) {
	for i, entry := range hostsfile.Entries {
		if entry == domain {
			hostsfile.Entries[i] = ""
		}
	}
}

// Uncomment the host file entries
func (hostsfile HostsFile) Start() error {
	return nil
}

// Comment the host file entries
func (hostsfile HostsFile) Stop() error {
	return nil
}

func (hostsfile HostsFile) Write(path string) error {
	return nil
}

func Parse(path string) (HostsFile, error) {
	hostsfile := HostsFile{}
	file, err := os.Open(path)

	if err != nil {
		return hostsfile, err
	}

	scanner := bufio.NewScanner(file)

	mark := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# CONCENTRATE") {
			mark = true
			continue
		}

		if strings.HasPrefix(line, "# END") {
			mark = false
			continue
		}

		if mark {
			fields := strings.Fields(line)
			hostsfile.Entries = append(hostsfile.Entries, fields[1])
		}

	}

	if err := scanner.Err(); err != nil {
		return HostsFile{}, err
	}

	return hostsfile, nil
}

func main() {
	flag.Parse()

	help := `
concentrate start
concentrate stop
concentrate add <domain>
concentrate remove <domain>
concentrate help
	`

	cmd := flag.Arg(0)

	hostfile, err := Parse("/etc/hosts")

	if err != nil {
		fmt.Println("ERROR", err)
		return
	}

	switch cmd {
	case "add":
		hostfile.Add("www.reddit.com")
	case "remove":
		hostfile.Remove("www.reddit.com")
	case "start":
		err = hostfile.Start()
	case "stop":
		err = hostfile.Stop()
	case "help":
	default:
		fmt.Println(help)
	}

	if err != nil {
		fmt.Println("ERROR", err)
	}
}
