TARGET=shellz

all: build

build:
	@go build -ldflags="-s -w" -o $(TARGET) cmd/shellz/*.go

clean:
	@rm -rf $(TARGET)
	@rm -rf build

install: build
	@mv $(TARGET) $(GOPATH)/bin/
