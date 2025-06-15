package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func findDest(domains []DomainConfig, domain string) (string, *DomainConfig, bool) {
	for _, d := range domains {
		if d.Domain == domain {
			return fmt.Sprintf("%s:%d", d.IP, d.Port), &d, true
		}
	}
	return "", nil, false
}

func isDomainChar(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '-' || c == '.'
}

func getDomain(data []byte) string {
	endBits := []byte{0x63, 0xdd, 0x02}
	domainStartIndex := 0
	domainEndIndex := len(data)

	for end := len(data) - len(endBits); end >= 0; end-- {
		if len(data) >= end+len(endBits) && string(data[end:end+len(endBits)]) == string(endBits) {
			domainEndIndex = end - 1
			break
		}
	}

	for start := domainEndIndex; start >= 0; start-- {
		c := data[start]
		if isDomainChar(c) {
			domainStartIndex = start
		} else {
			break
		}
	}

	if domainStartIndex < domainEndIndex && domainEndIndex < len(data) && domainStartIndex >= 0 {
		return string(data[domainStartIndex : domainEndIndex+1])
	}
	return ""
}

func handleConn(client net.Conn, cfg *Config) {
	defer client.Close()

	clientAddr := client.RemoteAddr().String()
	if !isAllowedIP(cfg.AllowedIPs, clientAddr) {
		log.Printf("access denied (global): %s\n", clientAddr)
		return
	}
	log.Printf("connected: %s\n", clientAddr)
	defer log.Printf("disconnected: %s\n", clientAddr)

	buf := make([]byte, 512)
	n, err := client.Read(buf)
	if err != nil {
		log.Println("read error:", err)
		return
	}

	data := buf[:n]
	domain := getDomain(data)
	if domain == "" {
		log.Println("failed to extract host")
		return
	}
	fmt.Printf("client ip: %s, client port: %s, host: %s\n", client.RemoteAddr().(*net.TCPAddr).IP.String(), fmt.Sprintf("%d", client.RemoteAddr().(*net.TCPAddr).Port), domain)

	dest, domainCfg, ok := findDest(cfg.Domains, domain)
	if !ok {
		log.Printf("unknown domain: %s\n", domain)
		return
	}
	if !isAllowedIP(domainCfg.AllowedIPs, clientAddr) {
		log.Printf("access denied (domain): %s for %s\n", clientAddr, domain)
		return
	}

	server, err := net.Dial("tcp", dest)
	if err != nil {
		log.Println("connect error:", err)
		return
	}
	defer server.Close()

	if _, err := server.Write(buf[:n]); err != nil {
		log.Println("write error:", err)
		return
	}

	go io.Copy(server, client)
	io.Copy(client, server)
}
