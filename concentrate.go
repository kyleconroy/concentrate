package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type HostsFile struct {
	Entries []string
	Started bool
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
func (hostsfile *HostsFile) Start() {
	hostsfile.Started = true
}

// Comment the host file entries
func (hostsfile *HostsFile) Stop() {
	hostsfile.Started = false
}

func (hostsfile HostsFile) Write(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	finfo, err := os.Stat(path)

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	lines := []string{}
	mark := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# CONCENTRATE") {
			lines = append(lines, line)

			for _, entry := range hostsfile.Entries {
				if hostsfile.Started {
					lines = append(lines, "127.0.0.1 "+entry)
				} else {
					lines = append(lines, "# "+entry)
				}
			}

			mark = true
			continue
		}

		if strings.HasPrefix(line, "# END") {
			mark = false
		}

		if !mark {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	content := strings.Join(lines, "\n")

	return ioutil.WriteFile(path, []byte(content+"\n"), finfo.Mode())
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
			if fields[0] == "127.0.0.1" {
				hostsfile.Started = true
			}
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
		log.Fatal(err)
		return
	}

	switch cmd {
	case "add":
		hostfile.Add(flag.Arg(1))
		err = hostfile.Write("/etc/hosts")
	case "remove":
		hostfile.Remove(flag.Arg(1))
		err = hostfile.Write("/etc/hosts")
	case "start":
		hostfile.Start()
		err = hostfile.Write("/etc/hosts")
	case "stop":
		hostfile.Stop()
		err = hostfile.Write("/etc/hosts")
	case "help":
	default:
		fmt.Println(help)
	}

	if err != nil {
		log.Fatal(err)
	}
}
