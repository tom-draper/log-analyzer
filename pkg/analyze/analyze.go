package analyze

import (
	"internal/server"
	"net"

	"reflect"

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

func getCountryCode(IPAddress string) string {
	if IPAddress == "" {
		return ""
	}
	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return ""
	}
	defer db.Close()

	ip := net.ParseIP(IPAddress)
	if ip == nil {
		return ""
	}
	record, err := db.Country(ip)
	if err != nil {
		return ""
	}
	location := record.Country.IsoCode
	return location
}

func ipLocations(extraction []parse.Extraction) map[int]string {
	ipLocations := make(map[int]string)
	for _, e := range extraction {
		for _, p := range e.Params {
			if ipAddress, ok := p.(net.IPAddress); ok {
				// s is string here
				location := getCountryCode(ipAddress)
				ipLocations[ipAddress] = location
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
