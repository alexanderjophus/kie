.PHONY: buf
buf:
	buf lint
	buf generate

.PHONY: lint
lint:
	go list -f '{{.Dir}}' -m | xargs golangci-lint run --verbose