## dnscheck script
Just query A record in a indefinite loop, save the log to CSV file for later

`go mod init dnscheck`

`go mod tidy`

`go build dnscheck.go`

To run (cltr+c to interupt):

`dnscheck.exe google.com 1.1.1.1 10.0.0.0/24`

Output and Log file looks like this (last value is TTL)
```
2026-08-14 09:51:08.273,1.1.1.1:53,google.com,10.0.0.0,24,142.250.197.142,253
2026-08-14 09:51:09.300,1.1.1.1:53,google.com,10.0.0.0,24,142.250.197.238,66
2026-08-14 09:51:10.339,1.1.1.1:53,google.com,10.0.0.0,24,142.250.199.78,197
2026-08-14 09:51:11.364,1.1.1.1:53,google.com,10.0.0.0,24,142.251.179.138,142
2026-08-14 09:51:12.391,1.1.1.1:53,google.com,10.0.0.0,24,142.250.71.174,278
2026-08-14 09:51:13.429,1.1.1.1:53,google.com,10.0.0.0,24,142.250.197.238,127
```
