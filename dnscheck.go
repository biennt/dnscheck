package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func lookup(resolver string, domain string, client_network_addr string, client_network_prefix uint8) string {
	// Automatically append the DNS port if the user didn't provide one
	if !strings.Contains(resolver, ":") {
		resolver = net.JoinHostPort(resolver, "53")
	}

	start := time.Now()
	// Format a clean timestamp string
	timestamp := start.Format("2006-01-02 15:04:05.000")
	// Create the EDNS0 Client Subnet (ECS) option
	ecs := &dns.EDNS0_SUBNET{
		Code:          dns.EDNS0SUBNET,
		Family:        1,                                // 1 for IPv4, 2 for IPv6
		SourceNetmask: client_network_prefix,            // The subnet mask length we are sending (/24)
		SourceScope:   0,                                // Set to 0 in queries; the server sets this in replies
		Address:       net.ParseIP(client_network_addr), // The proxy/client IP address
	}

	// 1. Configure the DNS client
	c := new(dns.Client)
	c.Timeout = 10 * time.Second

	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetUDPSize(4096)
	o.Option = append(o.Option, ecs)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.Extra = append(m.Extra, o)

	// 3. Execute the lookup
	in, _, err := c.Exchange(m, resolver)

	// Check for network errors or empty responses
	if err != nil || in == nil || len(in.Answer) == 0 {
		return fmt.Sprintf("%s,%s,%s,FAIL,0", timestamp, resolver, domain)
	}

	var returnip string
	var ttl uint32

	// 4. Extract the IP and TTL from the first A record found
	for _, ans := range in.Answer {
		if a, ok := ans.(*dns.A); ok {
			returnip = a.A.String()
			ttl = a.Header().Ttl
			break // We grab the first A record, matching your original logic
		}
	}

	// If no A record was found in the answer section (e.g., only CNAMEs)
	if returnip == "" {
		return fmt.Sprintf("%s,%s,%s,FAIL,0", timestamp, resolver, domain)
	}

	// Return formatted string: Timestamp, Resolver, Domain, IP, TTL
	return fmt.Sprintf("%s,%s,%s,%s,%d,%s,%d", timestamp, resolver, domain, client_network_addr, client_network_prefix, returnip, ttl)
}

func main() {
	filename := "dnscheck-output-" + strings.ReplaceAll(strings.ReplaceAll(time.Now().String()[:19], " ", "_"), ":", "_") + ".log"

	var fqdn, recursive_resolver, client_network_addr string
	var client_network_prefix_uint8 uint8

	if len(os.Args) == 4 {
		fqdn = os.Args[1]
		recursive_resolver = os.Args[2]
		client_network_subnet := os.Args[3]

		ip, ipnet, err := net.ParseCIDR(client_network_subnet)
		if err != nil {
			fmt.Println("Error parsing client network subnet CIDR:", err)
			return
		}
		ones, _ := ipnet.Mask.Size()
		client_network_addr = ip.String()
		client_network_prefix_uint8 = uint8(ones)
	} else if len(os.Args) == 5 {
		fqdn = os.Args[1]
		recursive_resolver = os.Args[2]
		client_network_addr = os.Args[3]
		client_network_prefix := os.Args[4]
		val64, err := strconv.ParseUint(client_network_prefix, 10, 8)
		if err != nil {
			fmt.Println("Error parsing string:", err)
			return
		}
		client_network_prefix_uint8 = uint8(val64)
	} else {
		fmt.Println("dnscheck version 0.1 <bien.nguyen@f5.com>")
		fmt.Println("Usage: dnscheck <fqdn> <recursive_resolver> <client_network_addr/client_network_prefix>")
		fmt.Println("Example: dnscheck google.com 1.1.1.1 10.0.0.0/24")
		return
	}

	file, fileErr := os.Create(filename)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
	defer file.Close()

	for {
		result := lookup(recursive_resolver, fqdn, client_network_addr, client_network_prefix_uint8)
		fmt.Fprintf(file, "%v\n", result)
		fmt.Println(result)

		time.Sleep(1 * time.Second)
	}
}
