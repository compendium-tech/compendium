package geoip

type GeoIp interface {
	GetLocation(ip string) (string, error)
}
