TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=autocloud.io
NAMESPACE=autoclouddev
NAME=autocloud
BINARY=terraform-provider-${NAME}
CITIZEN_ARCHIVE_NAME=${NAMESPACE}-${NAME}
OS_ARCH=darwin_amd64
## uncomment the following line if you are working locally
#VERSION=0.2
# provider source = autocloud.io/autocloud/iac
default: install

build:
	go build -o ${BINARY}
	chmod +x ${BINARY}

install:build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools:download
	@echo Installing tools from tools.go
	@cat autocloud_provider/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	rm ./binaries/${CITIZEN_ARCHIVE_NAME}*
	rm ./binaries/${BINARY}*

release: release-darwin_amd64 release-darwin_arm release-darwin_arm64 release-freebsd_386 release-freebsd_amd64 release-freebsd_arm release-freebsd_arm64 release-linux_386 release-linux_amd64 release-linux_arm release-linux_arm64 release-openbsd_386 release-openbsd_amd64 release-openbsd_arm release-openbsd_arm64 release-windows_386 release-windows_amd64 release-windows_arm release-windows_arm64

release-darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_darwin_amd64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_darwin_amd64.zip ./binaries/${BINARY}_v${VERSION}_darwin_amd64
	rm ./binaries/${BINARY}_v${VERSION}_darwin_amd64

release-darwin_arm:
	GOOS=darwin GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_darwin_arm
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_darwin_arm.zip ./binaries/${BINARY}_v${VERSION}_darwin_arm
	rm ./binaries/${BINARY}_v${VERSION}_darwin_arm

release-darwin_arm64:
	GOOS=darwin GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_darwin_arm64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_darwin_arm64.zip ./binaries/${BINARY}_v${VERSION}_darwin_arm64
	rm ./binaries/${BINARY}_v${VERSION}_darwin_arm64

release-freebsd_386:
	GOOS=freebsd GOARCH=386 go build -o ./binaries/${BINARY}_v${VERSION}_freebsd_386
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_freebsd_386.zip ./binaries/${BINARY}_v${VERSION}_freebsd_386
	rm ./binaries/${BINARY}_v${VERSION}_freebsd_386

release-freebsd_amd64:
	GOOS=freebsd GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_freebsd_amd64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_freebsd_amd64.zip ./binaries/${BINARY}_v${VERSION}_freebsd_amd64
	rm ./binaries/${BINARY}_v${VERSION}_freebsd_amd64

release-freebsd_arm:
	GOOS=freebsd GOARCH=arm go build -o ./binaries/${BINARY}_v${VERSION}_freebsd_arm
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_freebsd_arm.zip ./binaries/${BINARY}_v${VERSION}_freebsd_arm
	rm ./binaries/${BINARY}_v${VERSION}_freebsd_arm

release-freebsd_arm64:
	GOOS=freebsd GOARCH=arm go build -o ./binaries/${BINARY}_v${VERSION}_freebsd_arm64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_freebsd_arm64.zip ./binaries/${BINARY}_v${VERSION}_freebsd_arm64
	rm ./binaries/${BINARY}_v${VERSION}_freebsd_arm64

release-linux_386:
	GOOS=linux GOARCH=386 go build -o ./binaries/${BINARY}_v${VERSION}_linux_386
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_linux_386.zip ./binaries/${BINARY}_v${VERSION}_linux_386
	rm ./binaries/${BINARY}_v${VERSION}_linux_386

release-linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_linux_amd64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_linux_amd64.zip ./binaries/${BINARY}_v${VERSION}_linux_amd64
	rm ./binaries/${BINARY}_v${VERSION}_linux_amd64

release-linux_arm:
	GOOS=linux GOARCH=arm go build -o ./binaries/${BINARY}_v${VERSION}_linux_arm
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_linux_arm.zip ./binaries/${BINARY}_v${VERSION}_linux_arm
	rm ./binaries/${BINARY}_v${VERSION}_linux_arm

release-linux_arm64:
	GOOS=linux GOARCH=arm go build -o ./binaries/${BINARY}_v${VERSION}_linux_arm64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_linux_arm64.zip ./binaries/${BINARY}_v${VERSION}_linux_arm64
	rm ./binaries/${BINARY}_v${VERSION}_linux_arm64

release-openbsd_386:
	GOOS=openbsd GOARCH=386 go build -o ./binaries/${BINARY}_v${VERSION}_openbsd_386
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_openbsd_386.zip ./binaries/${BINARY}_v${VERSION}_openbsd_386
	rm ./binaries/${BINARY}_v${VERSION}_openbsd_386

release-openbsd_amd64:
	GOOS=openbsd GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_openbsd_amd64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_openbsd_amd64.zip ./binaries/${BINARY}_v${VERSION}_openbsd_amd64
	rm ./binaries/${BINARY}_v${VERSION}_openbsd_amd64

release-openbsd_arm:
	GOOS=openbsd GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_openbsd_arm
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_openbsd_arm.zip ./binaries/${BINARY}_v${VERSION}_openbsd_arm
	rm ./binaries/${BINARY}_v${VERSION}_openbsd_arm

release-openbsd_arm64:
	GOOS=openbsd GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_openbsd_arm64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_openbsd_arm64.zip ./binaries/${BINARY}_v${VERSION}_openbsd_arm64
	rm ./binaries/${BINARY}_v${VERSION}_openbsd_arm64

release-windows_386:
	GOOS=windows GOARCH=386 go build -o ./binaries/${BINARY}_v${VERSION}_windows_386
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_windows_386.zip ./binaries/${BINARY}_v${VERSION}_windows_386
	rm ./binaries/${BINARY}_v${VERSION}_windows_386

release-windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_windows_amd64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_windows_amd64.zip ./binaries/${BINARY}_v${VERSION}_windows_amd64
	rm ./binaries/${BINARY}_v${VERSION}_windows_amd64

release-windows_arm:
	GOOS=windows GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_windows_arm
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_windows_arm.zip ./binaries/${BINARY}_v${VERSION}_windows_arm
	rm ./binaries/${BINARY}_v${VERSION}_windows_arm

release-windows_arm64:
	GOOS=windows GOARCH=amd64 go build -o ./binaries/${BINARY}_v${VERSION}_windows_arm64
	zip -j ./binaries/${CITIZEN_ARCHIVE_NAME}_${VERSION}_windows_arm64.zip ./binaries/${BINARY}_v${VERSION}_windows_arm64
	rm ./binaries/${BINARY}_v${VERSION}_windows_arm64



.PHONY: release-darwin_amd64 release-darwin_arm release-darwin_arm64 release-freebsd_386 release-freebsd_amd64 release-freebsd_arm release-freebsd_arm64 release-linux_386 release-linux_amd64 release-linux_arm release-linux_arm64 release-openbsd_386 release-openbsd_amd64 release-openbsd_arm release-openbsd_arm64 release-windows_386 release-windows_amd64 release-windows_arm release-windows_arm64
