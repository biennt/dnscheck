package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func lookup(resolver string, domain string) string {
	// Automatically append the DNS port if the user didn't provide one
	if !strings.Contains(resolver, ":") {
		resolver = net.JoinHostPort(resolver, "53")
	}

	start := time.Now()
	// Format a clean timestamp string
	timestamp := start.Format("2006-01-02 15:04:05.000")

	// 1. Configure the DNS client
	c := new(dns.Client)
	c.Timeout = 10 * time.Second

	// 2. Prepare the query message
	m := new(dns.Msg)
	// dns.Fqdn ensures the domain has a trailing dot, required for raw queries
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

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
	return fmt.Sprintf("%s,%s,%s,%s,%d", timestamp, resolver, domain, returnip, ttl)
}

func main() {
	filename := "dnscheck-output-" + strings.ReplaceAll(strings.ReplaceAll(time.Now().String()[:19], " ", "_"), ":", "_") + ".log"
	if len(os.Args) == 3 {
		fqdn := os.Args[1]
		recursive_resolver := os.Args[2]
		file, fileErr := os.Create(filename)
		if fileErr != nil {
			fmt.Println(fileErr)
			return
		}
		defer file.Close()
		for {
			result := lookup(recursive_resolver, fqdn)
			fmt.Fprintf(file, "%v\n", result)
			fmt.Println(result)

			time.Sleep(1 * time.Second)
		}
	} else {
		fmt.Println("Usage: dnscheck <fqdn> <recursive_resolver>")
		fmt.Println("Example: dnscheck google.com 8.8.8.8")
	}
}