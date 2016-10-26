package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// host file names
type hosts []string

// do not block any hosts
var shieldsdown = hosts{"system"}

// block only the privacy hosts
var off = append(shieldsdown, "privacy")

// block the services
var on = append(off, "services")

func enable(set hosts) {
	data := []byte{}

	for _, name := range set {
		body, err := ioutil.ReadFile("/etc/hosts.d/" + name)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, body...)
	}

	write(data)
}

// writes out the hosts file
func write(hosts []byte) {
	err := ioutil.WriteFile("/etc/hosts", hosts, 0644)
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}
}

// ensures the necessary files exist
// TODO
func setup() {}

func main() {
	mode := "on"

	if len(os.Args) == 1 {
		fmt.Printf("Defaulting to 'on'. Options:\non\noff\nshieldsdown # let in everything\n\n")
	} else {
		mode = os.Args[1]
	}

	fmt.Printf("Switching to %s\n", mode)
	var set hosts
	switch mode {
	case "on":
		set = on
	case "off":
		set = off
	case "shieldsdown":
		set = shieldsdown
	default:
		set = on
	}

	enable(set)
	fmt.Printf("Using %v\n", set)

	// TODO detect if this is macOS or not
	exec.Command("/usr/bin/killall", "-HUP mDNSResponder").Output()
}
