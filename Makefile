


PHONY: gomat
gomat:
	go fmt ./...
	go vet ./...
	find . -name \*.go -not -path ./.git -not -path ./static -exec goimports -w {} \;
	golangci-lint run