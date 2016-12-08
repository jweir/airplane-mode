package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
)

// host file names
type hosts []string

// do not block any hosts
var shieldsdown = hosts{"system"}

// block only the privacy hosts
var off = append(shieldsdown, "privacy")

// block the services
var on = append(off, "services")

func enable(set hosts) error {
	data := []byte{}

	for _, name := range set {
		body, err := ioutil.ReadFile("/etc/hosts.d/" + name)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, body...)
	}

	return write(data)
}

// writes out the hosts file
func write(hosts []byte) error {
	err := ioutil.WriteFile("/etc/hosts", hosts, 0644)
	if err != nil {
		fmt.Printf("Error %s\nDid you run with 'sudo'?\n", err)
	}
	return err
}

// ensures the necessary files exist
// TODO
func setup() {}

type setting struct {
	hosts hosts
	desc  string
}

var settings = map[string]setting{
	"on":          setting{on, "Have a safe flight."},
	"off":         setting{off, "Please clean up after yourself."},
	"shieldsdown": setting{shieldsdown, "Commander are you crazy?!"},
}

func xinit() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	home := path.Join(u.HomeDir, ".airplane-mode")

	fmt.Println(home)

	if _, err := os.Stat(home); os.IsNotExist(err) {
		fmt.Printf(`
Now creating the airplane-mode home directory: %s

In this directory will be 3 files:

system   - The base host file that will always be present.
           Your existing /etc/hosts will be copied to this file now.
privacy  - A hosts file of services to block such as advertising and tracking.
           For example https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
services - Hosts that you want to block when airplane-mode is 'on'.
           These likely include services like Twitter, Facebook, etc.

		`, home)
	}

}

func main() {
	mode := "on"

	if len(os.Args) == 1 {
		fmt.Printf("Defaulting to 'on'.\nOptions:\n> on\n> off\n> shieldsdown # let in everything\n\n")
	} else {
		mode = os.Args[1]
	}

	set, ok := settings[mode]

	if ok == false {
		fmt.Print("This flight can not take off.  Please give us a command.")
	}

	err := enable(set.hosts)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", set.desc)

	// TODO detect if this is macOS or not
	exec.Command("/usr/bin/killall", "-HUP mDNSResponder").Output()
}
