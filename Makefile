lint:
	@golangci-lint run

changelog-lint:
	@changelog-lint

fix:
	@golangci-lint run --fix

install-hooks:
	@pre-commit install

install-changelog-lint:
	@go install github.com/chavacava/changelog-lint@master

dev-version:
	@./scripts/dev-version.sh

update-snapshots:
	@./scripts/update-snapshots.sh

test:
	@go test -v ./tests

git-tag:
	git tag --sign v$(v) -m v$(v)
	git push origin v$(v)

git-tag-remove-local:
	git tag -d v$(v)
	git fetch --prune --tags
