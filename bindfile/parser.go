package bindfile

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"teissem.fr/bind9_rest_api/bindmodel"
)

func NamedConfLocalParser(namedConfLocalFile string) ([]bindmodel.Zone, error) {
	// Compiling REGEX
	EMPTY_LINE := `^\s*$`
	COMMENT_LINE := `^\s*//.*$`
	BEGIN_ZONE_LINE := `^\s*zone\s*\"(.*?)\"\s*{\s*$`
	TYPE_ZONE_LINE := `^\s*type\s(.*?)\s*;\s*$`
	FILE_ZONE_LINE := `^\s*file\s*\"(.*?)\"\s*;\s*$`
	END_ZONE_LINE := `^\s*}\s*;\s*$`
	emptyLineRegex := regexp.MustCompile(EMPTY_LINE)
	commentLineRegex := regexp.MustCompile(COMMENT_LINE)
	beginZoneLineRegex := regexp.MustCompile(BEGIN_ZONE_LINE)
	typeZoneLineRegex := regexp.MustCompile(TYPE_ZONE_LINE)
	fileZoneLineRegex := regexp.MustCompile(FILE_ZONE_LINE)
	endZoneLineRegex := regexp.MustCompile(END_ZONE_LINE)
	// Parsing file
	file, err := os.Open(namedConfLocalFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Creates zones object from file
	var zones []bindmodel.Zone
	var currentZone *bindmodel.Zone = nil
	scanner := bufio.NewScanner(file)
	// Analyze line by line the file
	for scanner.Scan() {
		line := scanner.Text()
		if commentLineRegex.MatchString(line) || emptyLineRegex.MatchString(line) {
			continue
		} else if beginZoneLineRegex.MatchString(line) {
			currentZone = &bindmodel.Zone{
				Name:         beginZoneLineRegex.FindStringSubmatch(line)[1],
				FileLocation: "",
			}
		} else if endZoneLineRegex.MatchString(line) {
			zones = append(zones, *currentZone)
		} else if typeZoneLineRegex.MatchString(line) {
			continue
		} else if fileZoneLineRegex.MatchString(line) {
			currentZone.FileLocation = fileZoneLineRegex.FindStringSubmatch(line)[1]
		} else {
			fmt.Printf("Line not understand : %s\n", line)
		}
	}
	return zones, nil
}

func DBZoneParser(dbZoneFile string) (*bindmodel.DatabaseDNSZone, error) {
	// Compiling REGEX
	EMPTY_LINE := `^\s*$`
	COMMENT_LINE := `^\s*;.*$`                                           // Example : ; BIND data file for local loopback interface
	TTL_LINE := `^\s*\$TTL\s*(\d*?)\s*(;.*)*$`                           // Example : $TTL    604800
	MULTILINE_ENTRY_LINE := `^(.*?)\s+IN\s+(.*?)\s+(.*?)\s*\(\s*(;.*)*$` // Example : @       IN      SOA     localhost. root.localhost. (
	ADDITIONAL_INFORMATION_LINE := `^\s*(\d*?)\s*(;.*)*$`                // Example :                               2         ; Serial
	ENDING_ADDITIONAL_INFORMATION_LINE := `^\s*(\d*?)\s*\)\s*(;.*)*$`    // Example :                          604800 )       ; Negative Cache TTL
	SINGLE_ENTRY_LINE := `^(.*?)\s+IN\s+(.*?)\s+(.*?)\s*(;.*)*$`         // Example : home.lab.       IN      A       10.0.0.1
	emptyLineRegex := regexp.MustCompile(EMPTY_LINE)
	commentLineRegex := regexp.MustCompile(COMMENT_LINE)
	ttlLineRegex := regexp.MustCompile(TTL_LINE)
	multilineEntryLineRegex := regexp.MustCompile(MULTILINE_ENTRY_LINE)
	additionalInformationLineRegex := regexp.MustCompile(ADDITIONAL_INFORMATION_LINE)
	endingAdditionalInformationLineRegex := regexp.MustCompile(ENDING_ADDITIONAL_INFORMATION_LINE)
	singleEntryLineRegex := regexp.MustCompile(SINGLE_ENTRY_LINE)
	// Parsing file
	file, err := os.Open(dbZoneFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Creates dns zones object from file
	var databaseDNSZone bindmodel.DatabaseDNSZone
	var currentDNSZone bindmodel.DNSZone
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if commentLineRegex.MatchString(line) || emptyLineRegex.MatchString(line) {
			continue
		} else if ttlLineRegex.MatchString(line) {
			extractFromLine := ttlLineRegex.FindStringSubmatch(line)
			searchedNumber, err := strconv.Atoi(extractFromLine[1])
			if err == nil {
				databaseDNSZone.TTL = searchedNumber
			}
		} else if multilineEntryLineRegex.MatchString(line) {
			extractFromLine := multilineEntryLineRegex.FindStringSubmatch(line)
			currentDNSZone = bindmodel.DNSZone{
				DNS:       extractFromLine[1],
				EntryType: extractFromLine[2],
				IP:        extractFromLine[3],
			}
		} else if endingAdditionalInformationLineRegex.MatchString(line) {
			extractFromLine := endingAdditionalInformationLineRegex.FindStringSubmatch(line)
			searchedNumber, err := strconv.Atoi(extractFromLine[1])
			if err == nil {
				currentDNSZone.AdditionalInformation = append(currentDNSZone.AdditionalInformation, searchedNumber)
			}
			databaseDNSZone.DNSZones = append(databaseDNSZone.DNSZones, currentDNSZone)
		} else if additionalInformationLineRegex.MatchString(line) {
			extractFromLine := additionalInformationLineRegex.FindStringSubmatch(line)
			searchedNumber, err := strconv.Atoi(extractFromLine[1])
			if err == nil {
				currentDNSZone.AdditionalInformation = append(currentDNSZone.AdditionalInformation, searchedNumber)
			}
		} else if singleEntryLineRegex.MatchString(line) {
			extractFromLine := singleEntryLineRegex.FindStringSubmatch(line)
			databaseDNSZone.DNSZones = append(databaseDNSZone.DNSZones, bindmodel.DNSZone{
				DNS:       extractFromLine[1],
				EntryType: extractFromLine[2],
				IP:        extractFromLine[3],
			})
		}
	}
	return &databaseDNSZone, nil
}
