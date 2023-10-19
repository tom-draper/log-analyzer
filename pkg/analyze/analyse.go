package analyze

import "internal/server"

func Run(lines []map[string]string) {
	server.Start(lines)
}
