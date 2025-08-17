package geoip

// GeoIP is an interface for retrieving geographical location
// information based on an IP address.
type GeoIP interface {
	GetLocation(ip string) string
}
