TARGET=shellz

all: deps build

deps: godep
	@dep ensure

build: deps
	@go build -o $(TARGET) cmd/shellz/*.go

clean:
	@rm -rf $(TARGET)
	@rm -rf build

install: build
	@mv $(TARGET) $(GOPATH)/bin/

godep:
	@go get -u github.com/golang/dep/...
