package analyze

import "internal/server"

func Run(lines []map[string]any) {
	server.Start(lines)
}
