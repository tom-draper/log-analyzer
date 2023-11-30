module github.com/tom-draper/log-analyzer

go 1.21

require github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de

replace internal/display => ./internal/display

require internal/server v1.0.0

require internal/location v1.0.0

require (
	github.com/adrg/strutil v0.3.1 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/go-chi/chi/v5 v5.0.10 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/oschwald/geoip2-golang v1.9.0 // indirect
	github.com/oschwald/maxminddb-golang v1.12.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace internal/server => ./internal/server

replace internal/location => ./internal/location
