package analyze

import (
	"internal/server"
	"net"

	"reflect"

	"github.com/oschwald/geoip2-golang"
	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func dataTypeBreakdown(extraction []parse.Extraction) map[string]map[string]int {
	dataTypes := make(map[string]map[string]int)
	for _, e := range extraction {
		for k, p := range e.Params {
			dataType := reflect.TypeOf(p)
			if _, ok := dataTypes[k]; !ok {
				dataTypes[k] = make(map[string]int)
			}
			dataTypes[k][dataType.String()] += 1
		}
	}
	return dataTypes
}

func failedLines(extraction []parse.Extraction) map[int]string {
	failedLines := make(map[int]string)
	for _, e := range extraction {
		if len(e.Params) == 0 {
			failedLines[e.LineNumber] = e.Line
		}
	}
	return failedLines
}

func getCountryCode(ipAddress net.IP) string {
	if ipAddress == nil {
		return ""
	}
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return ""
	}
	defer db.Close()

	record, err := db.Country(ipAddress)
	if err != nil {
		return ""
	}
	location := record.Country.IsoCode
	return location
}

func ipLocations(extraction []parse.Extraction) map[string]string {
	ipLocations := make(map[string]string)
	for _, e := range extraction {
		for _, p := range e.Params {
			if ipAddress, ok := p.(net.IP); ok {
				ipAddressStr := ipAddress.String()
				// Check if country code already exists
				if _, ok := ipLocations[ipAddressStr]; ok {
					continue
				}
				location := getCountryCode(ipAddress)
				ipLocations[ipAddress.String()] = location
			}

		}
	}
	return ipLocations
}

func NewData(extraction []parse.Extraction, config *parse.Config) *server.Data {
	dataTypes := dataTypeBreakdown(extraction)
	failed := failedLines(extraction)
	data := server.Data{Extraction: extraction, DataTypes: dataTypes, Failed: failed, Config: config}
	return &data
}

func Run(extraction []parse.Extraction, config *parse.Config) {
	data := NewData(extraction, config)
	server.Start(data)
}
