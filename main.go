package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type DomainConfig struct {
	Domain     string   `yaml:"domain"`
	IP         string   `yaml:"ip"`
	Port       int      `yaml:"port"`
	AllowedIPs []string `yaml:"allowed_ips,omitempty"` // ドメインごとのアクセス制限
}

type Config struct {
	Listen     string         `yaml:"listen"`
	Domains    []DomainConfig `yaml:"domains"`
	AllowedIPs []string       `yaml:"allowed_ips,omitempty"` // 全体のアクセス制限
}

// アクセス制限チェック
func isAllowedIP(allowed []string, remoteAddr string) bool {
	if len(allowed) == 0 {
		return true // 制限なし
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

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func findDest(domains []DomainConfig, domain string) (string, *DomainConfig, bool) {
	for _, d := range domains {
		if d.Domain == domain {
			return fmt.Sprintf("%s:%d", d.IP, d.Port), &d, true
		}
	}
	return "", nil, false
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
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '-' || c == '.' {
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
	log.Printf("client ip: %s, client port: %s, host: %s\n", client.RemoteAddr().(*net.TCPAddr).IP.String(), fmt.Sprintf("%d", client.RemoteAddr().(*net.TCPAddr).Port), domain)

	dest, domainCfg, ok := findDest(cfg.Domains, domain)
	if !ok {
		log.Printf("unknown domain: %s\n", domain)
		return
	}
	// ドメインごとのアクセス制限
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

func main() {
	cfg, err := loadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s\n", cfg.Listen)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn, cfg)
	}
}
