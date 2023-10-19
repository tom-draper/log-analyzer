module github.com/tom-draper/log-analyzer

go 1.21

require github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de

require internal/display v1.0.0

replace internal/display => ./internal/display

require internal/server v1.0.0

require github.com/go-chi/chi/v5 v5.0.10 // indirect

replace internal/server => ./internal/server
