## dnscheck script
Just query A record in a indefinite loop, save the log to CSV file for later

`go mod init dnscheck`

`go mod tidy`

`go build dnscheck.go`

To run:

`dnscheck.exe google.com 8.8.8.8`
