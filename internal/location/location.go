package location

import (
	"errors"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var (
	dbOnce   sync.Once
	dbReader *geoip2.Reader
	dbErr    error
)

func GetCountryCode(ipAddress net.IP) (string, error) {
	if ipAddress == nil {
		return "", errors.New("invalid IP address")
	}
	dbOnce.Do(func() {
		dbReader, dbErr = geoip2.Open("internal/location/GeoLite2-Country.mmdb")
	})
	if dbErr != nil {
		return "", dbErr
	}

	record, err := dbReader.Country(ipAddress)
	if err != nil {
		return "", err
	}
	location := record.Country.Names["en"]
	return location, nil
}
