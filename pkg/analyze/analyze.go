package analyze

import (
	"internal/server"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func Run(extraction []parse.Extraction, config *parse.Config) {
	server.Start(extraction, config)
}
