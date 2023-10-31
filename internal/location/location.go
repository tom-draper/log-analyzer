package location

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func GetCountryCode(ipAddress net.IP) (string, error) {
	if ipAddress == nil {
		return "", errors.New("invalid IP address")
	}
	db, err := geoip2.Open("internal/location/GeoLite2-Country.mmdb")
	if err != nil {
		return "", err
	}
	defer db.Close()

	record, err := db.Country(ipAddress)
	if err != nil {
		return "", err
	}
	location := record.Country.Names["en"]
	return location, nil
}
