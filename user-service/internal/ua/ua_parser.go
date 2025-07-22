package ua

import (
	"fmt"

	"github.com/mileusna/useragent"
)

type UserAgentInfo struct {
	Name   string
	Os     string
	Device string
}

type UserAgentParser interface {
	ParseUserAgent(ua string) UserAgentInfo
}

type userAgentParser struct{}

func NewUserAgentParser() UserAgentParser {
	return &userAgentParser{}
}

func (p *userAgentParser) ParseUserAgent(ua string) UserAgentInfo {
	parsed := useragent.Parse(ua)
	return UserAgentInfo{
		Name:   fmt.Sprintf("%s %s", parsed.Name, parsed.Version),
		Os:     fmt.Sprintf("%s %s", parsed.OS, parsed.OSVersion),
		Device: parsed.Device,
	}
}
