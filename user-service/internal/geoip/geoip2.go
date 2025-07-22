package geoip

import (
	"fmt"

	"github.com/iamvladw/GeoIP2-go/cmd/geoip2"
)

type geoIp2Client struct {
	api *geoip2.Api
}

func NewGeoIp2Client(accountId, licenseKey, host string) GeoIp {
	api := geoip2.New(accountId, licenseKey, host)
	return &geoIp2Client{api: api}
}

func (g *geoIp2Client) GetLocation(ip string) (string, error) {
	city, err := g.api.City(ip)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, %s", city.City.Names.En, city.Country.Names.En), nil
}
