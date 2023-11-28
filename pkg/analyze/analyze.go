package analyze

import (
	"internal/location"
	"internal/server"
	"net"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func ipLocations(extraction []parse.Extraction) map[string]string {
	ipLocations := make(map[string]string)
	for _, e := range extraction {
		for _, p := range e.Params {
			ipAddress, ok := p.Value.(net.IP)
			if !ok {
				continue
			}
			ipAddressStr := ipAddress.String()
			// Check if country code already exists
			if _, ok := ipLocations[ipAddressStr]; ok {
				continue
			}
			loc, err := location.GetCountryCode(ipAddress)
			if err != nil {
				continue
			}
			ipLocations[ipAddress.String()] = loc
		}
	}
	return ipLocations
}

func NewData(extraction []parse.Extraction, config *parse.Config) *server.Data {
	locations := ipLocations(extraction)
	data := server.Data{Extraction: extraction, Locations: locations, Config: config}
	return &data
}

func Run(extraction []parse.Extraction, config *parse.Config) {
	data := NewData(extraction, config)
	server.Start(data)
}
