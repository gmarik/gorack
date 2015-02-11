# Gorack

[gorack] is a [Go] backed fronted webserver for Ruby's [Rack] applications


# Current state

1. *alpha quality*
2. Gem file ships with only Darwin amd64 prebuilt binary. Why? See 1. Feel free to submit PRs for other OSes and ARCHes

# Why

An experiment; inspired by [node]'s [nack]

# How To
## Get Up And Running

1. `gem install gorack`
2. `gorack -config ./path/to/config.ru` 

ie:

1. `gorack -config=$(dirname $(gem which gorack))/../test/echo.ru`
2. `open http://localhost:3000`

## Develop

1. git clone http://github.com/gmarik/gorack
2. cd gorack
3. make

Builds gem file

## Testing

1. cd gorack/
2. go test -v .
3. go test -v -bench=. . # benchmarking

Requires [Go] installed. Developed with 1.4 version

## TODO

[x] fix weird zombie leaking
[ ] improve error handling: broken IPC results in malfunction of the parent ruby process
[ ] improve performance


[Go]: http://golang.org
[gorack]: http://github.com/gmarik/gorack
[nack]: http://github.com/josh/nack
[Rack]: http://rack.github.io
[node]: http://nodejs.org
