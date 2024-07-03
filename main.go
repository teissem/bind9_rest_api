package main

import (
	"fmt"

	"teissem.fr/bind9_rest_api/bindfile"
)

func main() {
	zones, err := bindfile.NamedConfLocalParser("example/named.conf.local")
	if err != nil {
		panic(0)
	}
	fmt.Println(len(zones))
	fmt.Println(zones)
	DNSZones, err := bindfile.DBZoneParser("example/db.home.lab")
	if err != nil {
		panic(0)
	}
	fmt.Println(len(DNSZones))
	fmt.Println(DNSZones)
}
