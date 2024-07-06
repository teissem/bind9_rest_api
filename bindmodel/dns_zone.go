package bindmodel

type DNSZone struct {
	DNS                   string
	EntryType             string
	IP                    string
	AdditionalInformation []int
}

type DatabaseDNSZone struct {
	TTL      int
	DNSZones []DNSZone
}
