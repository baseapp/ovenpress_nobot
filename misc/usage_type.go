package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	// Load the MaxMind GeoIP2 Anonymous IP database
	// (Note: You need to download the actual GeoIP2-Anonymous-IP.mmdb file)
	db, err := geoip2.Open("../assets/usagetype.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// The IP you want to check
	ipStr := "106.219.155.28"
	if len(os.Args) > 1 {
		ipStr = os.Args[1]
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Fatalf("Invalid IP address: %s", ipStr)
	}

	record, err := db.AnonymousIP(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("IP: %s\n", ip.String())
	fmt.Printf("Is VPN: %t\n", record.IsAnonymousVPN)
	fmt.Printf("Is Hosting/Datacenter: %t\n", record.IsHostingProvider)
	fmt.Printf("Is Tor Exit Node: %t\n", record.IsTorExitNode)
	fmt.Printf("Is Public Proxy: %t\n", record.IsPublicProxy)
}
