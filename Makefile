version=$(shell git describe --tags --abbrev=0)
build_date=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
commit_hash=$(shell git rev-parse HEAD)

build:
	go build -ldflags="-s -X 'github.com/Tubular-Bytes/tf-runner/pkg/version.CommitHash=$(commit_hash)' -X 'github.com/Tubular-Bytes/tf-runner/pkg/version.Version=$(version)' -X 'github.com/Tubular-Bytes/tf-runner/pkg/version.BuildTime=$(build_date)'" -o bin/runner ./cmd/...

release:
	git tag $(shell svu)
	git push
	git push --tags