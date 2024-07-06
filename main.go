package main

import (
	"teissem.fr/bind9_rest_api/bindfile"
)

func main() {
	zones, err := bindfile.NamedConfLocalParser("example/named.conf.local")
	if err != nil {
		panic(0)
	}
	bindfile.GenerateNamedConfLocal(zones, "example/named.conf.local.generated")
	DNSZones, err := bindfile.DBZoneParser("example/db.home.lab")
	if err != nil {
		panic(0)
	}
	bindfile.GenerateDBZone(DNSZones, "example/db.home.lab.generated")
}
