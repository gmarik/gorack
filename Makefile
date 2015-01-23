GOPATH := $(GOPATH)
BUILDDIR := ./build

# target: all - Run tests and generate binary
all: test build

# target: help - Display targets
help:
	@egrep "^# target:" [Mm]akefile | sort - |sed 's/# target://'

# target: clean - Cleans build artifacts
clean:
	echo Cleaning build artifacts...
	go clean
	rm -rf ${BUILDDIR}
	echo

# target: test - Runs CLI tests
test:
	echo Testing packages:
	go test .

prep:
	mkdir -p ${BUILDDIR}/ruby/{bin,libexec,lib}

copyfiles:
	cp -r ruby/* rubygem/* ${BUILDDIR}/ruby/

binaries:
	GOARCH=amd64 GOOS=darwin go build -o ${BUILDDIR}/ruby/libexec/darwin_gorack main/gorack-server.go

gemfile:
	cd ${BUILDDIR}/ruby/ && gem build gorack.gemspec

# target: build - Build CLI binary
build: clean prep copyfiles binaries gemfile
	echo

.PHONY: all help clean build
