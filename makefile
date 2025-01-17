SHELL=/bin/bash

packages = ./core ./prompts

run:
	@go run ./examples
play:
	@go run ./playground
test:
	go clean -testcache
	@until [ $$RET -eq 0 ]; do \
		go test $(packages) -cover ; \
		RET=$$? ; \
	done
profile:
	go clean -testcache
	@until [ $$RET -eq 0 ]; do \
		go test $(packages) -cover -coverprofile cover.out ; \
		RET=$$? ; \
	done
	go tool cover -html cover.out -o cover.html
	rm cover.out
profile-core:
	go clean -testcache
	@until [ $$RET -eq 0 ]; do \
		go test ./core -cover -coverprofile cover.out ; \
		RET=$$? ; \
	done
	go tool cover -html cover.out -o cover.html
	rm cover.out
profile-prompts:
	go clean -testcache
	@until [ $$RET -eq 0 ]; do \
		go test ./prompts -cover -coverprofile cover.out ; \
		RET=$$? ; \
	done
	go tool cover -html cover.out -o cover.html
	rm cover.out
snap:
	go clean -testcache
	@UPDATE_SNAPSHOTS=true go test ./prompts
format:
	gofmt -w .
ci: test format
clog:
	git-chglog -o CHANGELOG.md
