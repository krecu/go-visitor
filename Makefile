BINARY = go-visitor
GOARCH = amd64

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION}"

# Build the project
all: linux darwin windows

linux:
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \
	cd - >/dev/null