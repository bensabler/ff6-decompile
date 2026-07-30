.PHONY: fmt test vet fuzz check build clean

fmt:
	gofmt -w .

test:
	go test ./...

vet:
	go vet ./...

fuzz:
	@for pkg in $$(go list ./...); do \
		for f in $$(go test -list '^Fuzz' $$pkg | grep '^Fuzz' || true); do \
			go test $$pkg -run='^$$' -fuzz="^$${f}\$$" -fuzztime=30s || exit 1; \
		done; \
	done

check: fmt test vet

build:
	go build ./cmd/ff6lab

clean:
	rm -rf out dist coverage.out ff6lab ff6lab.exe
