BUILD_VERSION=$(shell git tag -l --contains HEAD)
BUILD_DATE=$(shell date +'%Y/%m/%d %H:%M:%S')
BUILD_COMMIT=$(shell git log -1 --pretty=%H)

build-staticlint:
	go build -o bin/staticlint cmd/staticlint/main.go

run-staticlint:
	./bin/staticlint ./...

run-staticlint-with-ignore-tests:
	./bin/staticlint -test=false ./...

run-server-with-build-info:
	go run -ldflags \
	"-X 'main.BuildVersion=$(BUILD_VERSION)' \
	-X 'main.BuildDate=$(BUILD_DATE)' \
	-X 'main.BuildCommit=$(BUILD_COMMIT)'" \
	cmd/server/main.go

run-agent-with-build-info:
	go run -ldflags \
	"-X 'main.BuildVersion=$(BUILD_VERSION)' \
	-X 'main.BuildDate=$(BUILD_DATE)' \
	-X 'main.BuildCommit=$(BUILD_COMMIT)'" \
	cmd/agent/main.go

web-doc:
	godoc -http=:8080