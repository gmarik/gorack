# Gorack

[gorack] is a [Go] backed frontend webserver for Ruby's [Rack] applications


## Current state

1. *alpha quality*
2. Gem file ships with only Darwin amd64 prebuilt binary. Why? See 1.

# Why

An experiment; inspired by [node]'s [nack]

# How To
## Building from sources

0. `export GORACKPATH=$GOPATH/src/github.com/gmarik`
1. `mkdir -p $GORACKPATH`
2. `git clone http://github.com/gmarik/gorack $GORACKPATH`
3. `cd $GORACKPATH`
4. `go run main/gorack-server.go -config ./ruby/test/echo.ru`
5. `open http://localhost:3000`


## Building Gemfile

0. `cd $GORACKPATH`
1. `make gemfile # builds to ./build/ruby/gorack-x.x.x.gem`


## Testing

0. `cd $GORACKPATH`
1. `go test -v .`
2. `go test -v -bench=. . # with benchmarking`


## Get Up And Running

NOTE: gem provides only for x86_64 OSX binary at the moment. See instructions how to build from source

1. `gem install gorack`
2. `gorack -config ./path/to/config.ru` 

ie:

1. `gorack -config=$(dirname $(gem which gorack))/../test/echo.ru`
2. `open http://localhost:3000`


## TODO

- [x] fix weird zombie leaking
- [ ] improve error handling: broken IPC results in malfunction of the parent ruby process
- [ ] improve performance


[Go]: http://golang.org
[gorack]: http://github.com/gmarik/gorack
[nack]: http://github.com/josh/nack
[Rack]: http://rack.github.io
[node]: http://nodejs.org
