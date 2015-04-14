.PHONY: build test python-test server-test clean

build:
	cd server && go build -o ../bin/airecalld ./cmd/airecalld

