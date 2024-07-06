package bindfile

import (
	"fmt"
	"os"

	"teissem.fr/bind9_rest_api/bindmodel"
)

func openCleanFile(fileLocation string) (*os.File, error) {
	/* openCleanFile creates or open file ready to be used.
	If the file is open, the file is wiped completely */
	// Wipe file if exists
	_ = os.Truncate(fileLocation, 0)
	// Create or open file
	file, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, err
}

func GenerateNamedConfLocal(zones []bindmodel.Zone, fileLocation string) error {
	/* GenerateNamedConfLocal generates a named.conf.local at fileLocation
	The content is not analyzed / added, the file is entirely regenerated */
	file, err := openCleanFile(fileLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	// Write each zones in the named.conf.local
	for _, zone := range zones {
		file.WriteString(fmt.Sprintf("zone \"%s\" {\n", zone.Name))
		file.WriteString("\ttype master;\n")
		file.WriteString(fmt.Sprintf("\tfile \"%s\";\n", zone.FileLocation))
		file.WriteString("};\n\n")
	}
	return nil
}

func GenerateDBZone(databaseDNSZone *bindmodel.DatabaseDNSZone, fileLocation string) error {
	/* GenerateDBZone generates a db.zone.zone at fileLocation
	The content is not analyzed / added, the file is entirely regenerated */
	file, err := openCleanFile(fileLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	// Write TTL first
	file.WriteString(fmt.Sprintf("$TTL\t%d\n", databaseDNSZone.TTL))
	// Write each DNS Zone in the bind file
	for _, dnsZone := range databaseDNSZone.DNSZones {
		if len(dnsZone.AdditionalInformation) > 0 {
			file.WriteString(fmt.Sprintf(
				"%s\tIN\t%s\t%s\t(\n",
				dnsZone.DNS,
				dnsZone.EntryType,
				dnsZone.IP))
			for index, additionalInformation := range dnsZone.AdditionalInformation {
				if index == len(dnsZone.AdditionalInformation)-1 {
					file.WriteString(fmt.Sprintf("\t%d\t)\n", additionalInformation))
				} else {
					file.WriteString(fmt.Sprintf("\t%d\n", additionalInformation))
				}
			}
		} else {
			file.WriteString(fmt.Sprintf(
				"%s\tIN\t%s\t%s\n",
				dnsZone.DNS,
				dnsZone.EntryType,
				dnsZone.IP))
		}
	}
	return nil
}
