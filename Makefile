GOPATH := $(GOPATH)
BUILDDIR := ./build

# target: all - Run tests and generate binary
all: help

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

# target: gemfile - builds gemfile from gemspec
gemfile: build
	cd ${BUILDDIR}/ruby/ && gem build gorack.gemspec

# target: build - Build CLI binary
build: clean prep copyfiles binaries
	echo

# target: gem_install - builds and `gem install`s gemfile
gem_install: gemfile
	gem install ${BUILDDIR}/ruby/gorack-*.gem

.PHONY: all help clean build
