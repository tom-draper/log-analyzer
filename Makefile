.PHONY: build test check dashboard clean

build:
	go build -o log-analyzer .

test:
	go test ./...

check:
	go vet ./...
	cd dashboard && pnpm run check

dashboard:
	cd dashboard && pnpm run build

clean:
	rm -f log-analyzer
