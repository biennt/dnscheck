package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
	"bufio"
	"os"
)

func lookup(resolver string, domain string, ipversion string) string {
	var retstr string
	var stringValue string

	start := time.Now()
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, resolver)
		},
	}
	
	ip, _ := r.LookupIP(context.Background(), ipversion, domain)
	
	if len(ip) > 0 {
		elapsed := time.Since(start).Milliseconds()
		stringValue = strconv.FormatInt(elapsed, 10)
		retstr = "(" + resolver + ") -- " + domain + " --> " + ip[0].String() + ", response time: " + stringValue + "ms"
	} else {
		retstr = "(" + resolver + ") -- " + domain + " --> FAIL"
	}
	return retstr
}

func check_viettel_dns() {
	var list_of_resolver_v4 [4]string
	list_of_resolver_v4[0] = "116.97.90.124:53"
	list_of_resolver_v4[1] = "116.97.90.125:53"
	list_of_resolver_v4[2] = "116.97.90.126:53"
	list_of_resolver_v4[3] = "116.97.90.127:53"
	var list_of_resolver_v6 [4]string
	list_of_resolver_v6[0] = "[2402:800:20ff:109c::1]:53"
	list_of_resolver_v6[1] = "[2402:800:20ff:305a::3]:53"
	list_of_resolver_v6[2] = "[2402:800:20ff:2056::3]:53"
	list_of_resolver_v6[3] = "[2402:800:20ff:109c::3]:53"
	var list_of_domain_v4 [2]string
	list_of_domain_v4[0] = "v4.bienlab.com"
	list_of_domain_v4[1] = "google.com"
	var list_of_domain_v6 [2]string
	list_of_domain_v6[0] = "v6.bienlab.com"
	list_of_domain_v6[1] = "google.com"

	fmt.Println("------------ Checking DNS resolvers (IPv4) ------------ ")
	for i := 0; i < len(list_of_domain_v4); i++ {
		for j := 0; j < len(list_of_resolver_v4); j++ {
			fmt.Println(lookup(list_of_resolver_v4[j], list_of_domain_v4[i],"ip4"))
		}
	}
	fmt.Println("------------ Checking DNS resolvers (IPv6) ------------ ")
	for i := 0; i < len(list_of_domain_v6); i++ {
		for j := 0; j < len(list_of_resolver_v6); j++ {
			fmt.Println(lookup(list_of_resolver_v6[j], list_of_domain_v6[i],"ip6"))
		}
	}
	fmt.Println("------------------------------------------------------- ")
}

func check_current_dns(){
	var domain string = "bienlab.com"
	fmt.Println("-------------- Testing current DNS settings ----------- ")
	dnsrecord, err := net.LookupHost(domain)
	if err != nil {
		fmt.Println("Error when lookup with current DNS settings:", err)
	} else {
		fmt.Print("Lookup ",domain, " --> ", dnsrecord[0])
		ip := net.ParseIP(dnsrecord[0])
		if ip.To4() != nil {
			fmt.Println(" (IPv4)")
		} else {
			fmt.Println(" (IPv6)")
		}
	}
	fmt.Println("------------------------------------------------------- ")
}

func header() {
	fmt.Println("------------------------------------------------------- ")
	fmt.Println("By Bien <bien.nguyen@f5.com>")
	fmt.Println("------------------------------------------------------- ")
	fmt.Println("")
}

func pause(){
	fmt.Println("Press Enter to finish...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
func main() {
	header()
	check_current_dns()
	check_viettel_dns()
	pause()
}
