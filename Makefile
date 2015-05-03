.PHONY: build test python-test server-test clean

build:
	cd server && go build -o ../bin/airecalld ./cmd/airecalld

test: python-test server-test

python-test:
	cd sdk-python && pip install -e ".[dev]" && pytest

server-test:
