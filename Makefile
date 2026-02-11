generate:
	@go generate ./...

build:
	@go generate ./...
	@GOOS=linux go build -ldflags="-s -w" -o ./bin/openapi-tsgen main.go

lint:
	@golangci-lint run

changelog-lint:
	@changelog-lint

version-lint:
	@./scripts/lint-version.sh

fix:
	@golangci-lint run --fix

install-hooks:
	@pre-commit install

install-changelog-lint:
	@go install github.com/chavacava/changelog-lint@master

dev-version:
	@./scripts/dev-version.sh

release:
	@./scripts/release.sh $(v)

update-snapshots:
	@./scripts/update-snapshots.sh

test:
	@go test -v ./tests

ci-build:
	@go generate ./...
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./build/openapi-tsgen-windows-amd64.exe main.go
	@GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o ./build/openapi-tsgen-windows-arm64.exe main.go
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/openapi-tsgen-linux-amd64 main.go
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./build/openapi-tsgen-linux-arm64 main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./build/openapi-tsgen-darwin-amd64 main.go
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./build/openapi-tsgen-darwin-arm64 main.go

git-tag:
	@./scripts/release.sh $(v)
	@./scripts/lint-version.sh
	@git tag --sign v$(v) -m v$(v)
	@git push origin v$(v)

git-tag-remove:
	@git tag -d v$(v)
	@git push --delete origin v$(v)
