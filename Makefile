BINARY=bookr
LINUX_AMD=$(BINARY)_linux_amd64
LINUX_ARM=$(BINARY)_linux_arm64
DARWIN_AMD=$(BINARY)_darwin_amd64
DARWIN_ARM=$(BINARY)_darwin_arm64


$(LINUX_AMD):
	env GOOS=linux GOARCH=amd64 go build -o bin/$(LINUX_AMD) -ldflags="-s -w" ./main.go
$(LINUX_ARM):
	env GOOS=linux GOARCH=arm64 go build -o bin/$(LINUX_ARM) -ldflags="-s -w" ./main.go

linux: $(LINUX_AMD) $(LINUX_ARM)

$(DARWIN_AMD):
	env GOOS=darwin GOARCH=amd64 go build -o bin/$(DARWIN_AMD) -ldflags="-s -w" ./main.go

$(DARWIN_ARM):
	env GOOS=darwin GOARCH=arm64 go build -o bin/$(DARWIN_ARM) -ldflags="-s -w" ./main.go

darwin: $(DARWIN_AMD) $(DARWIN_ARM)

build: linux darwin

.PHONY: test
test:
	@go test ./...

.PHONY: clean
clean:
	@rm -rf bin
