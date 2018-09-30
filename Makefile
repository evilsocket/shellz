TARGET=shellz

all: deps build

deps: godep
	@dep ensure

build:
	@go build -o $(TARGET) cmd/shellz/*.go

clean:
	@rm -rf $(TARGET)
	@rm -rf build

install:
	@cp $(TARGET) /usr/local/bin/

godep:
	@go get -u github.com/golang/dep/...
