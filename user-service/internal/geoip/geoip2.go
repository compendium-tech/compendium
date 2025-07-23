package geoip

import (
	"fmt"

	"github.com/iamvladw/GeoIP2-go/cmd/geoip2"
)

type geoIP2Client struct {
	api *geoip2.Api
}

func NewGeoIP2Client(accountID, licenseKey, host string) GeoIP {
	api := geoip2.New(accountID, licenseKey, host)
	return &geoIP2Client{api: api}
}

func (g *geoIP2Client) GetLocation(ip string) (string, error) {
	city, err := g.api.City(ip)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, %s", city.City.Names.En, city.Country.Names.En), nil
}
