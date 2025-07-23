package geoip

type GeoIP interface {
	GetLocation(ip string) (string, error)
}
