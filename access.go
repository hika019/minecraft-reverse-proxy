package main

import (
	"net"
)

func isAllowedIP(allowed []string, remoteAddr string) bool {
	if len(allowed) == 0 {
		return true
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return false
	}
	for _, ip := range allowed {
		if host == ip {
			return true
		}
	}
	return false
}
