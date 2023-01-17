snot:
	go build -o /usr/local/bin/$@ ./cmd/cli/main.go

clean:
	rm -f /usr/local/bin/snot

.PHONY: clean
