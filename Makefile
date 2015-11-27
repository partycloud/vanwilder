BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOLDFLAGS="-X main.branch $(BRANCH) -X main.commit $(COMMIT)"

default: build

get:
	@go get -d ./...

build: get
	@mkdir -p bin
	@go build -ldflags=$(GOLDFLAGS) -a -o bin/vw ./cmd/vw
