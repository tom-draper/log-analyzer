.PHONY: build test dashboard clean

build:
	go build -o log-analyzer .

test:
	go test ./...

dashboard:
	cd dashboard && pnpm run build

clean:
	rm -f log-analyzer
