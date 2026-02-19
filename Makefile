HOSTNAME=local.com
NAMESPACE=icdc-io
NAME=icdc
BINARY=terraform-provider-${NAME}
VERSION=1.0.0
OS_ARCH=darwin_arm64

default: install

check:
	golangci-lint run -v

build:
	go build -o ${BINARY}

#install: check build
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

clean:
	rm -f ${BINARY}
	rm -rf ./bin