.PHONY: build test python-test server-test clean

build:
	cd server && go build -o ../bin/airecalld ./cmd/airecalld

test: python-test server-test

python-test:
	cd sdk-python && pip install -e ".[dev]" && pytest

server-test:
	cd server && go test ./... -race -count=1

clean:
	rm -rf bin
<!-- draft note 1456 -->
