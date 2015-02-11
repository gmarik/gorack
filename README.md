# Gorack

[gorack] is a [Go] backed fronted webserver for Ruby's [Rack] applications


# Current state

*alpha quality*

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
2. go test .
3. go test -bench=. .

Requires [Go] installed. Developed with 1.4 version

## TODO

[x] fix weird zombie leaking
[ ] improve error handling: broken IPC results in malfunction of the parent ruby process


[Go]: http://golang.org
[gorack]: http://github.com/gmarik/gorack
[nack]: http://github.com/josh/nack
[Rack]: http://rack.github.io
[node]: http://nodejs.org
