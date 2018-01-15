BINARY = go-visitor
GOARCH = amd64

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION}"

# Build the project
all: linux darwin windows

linux:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	GOOS=linux GOARCH=${GOARCH} go build -o ${OUTPUT}/${BINARY}-linux-${GOARCH} . ;

darwin:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${OUTPUT}/${BINARY}-darwin-${GOARCH} . ;

windows:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${OUTPUT}/${BINARY}-windows-${GOARCH}.exe . ;