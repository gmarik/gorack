# Gorack

[gorack] is a [Go] backed fronted webserver for Ruby's [Rack] applications


# Current state

*alpha quality*

# Why

An experiment; inspired by [nodejs]' [nack]

# How To
## Get Up And Running

1. gem install gorack
2. gorack -config ./path/to/config.ru

## Develop

1. git clone http://github.com/gmarik/gorack
2. cd gorack
3. make

Builds gem file

## Testing

1. cd gorack/
2. go test .

Requires [Go] installed. Developed with 1.4 version

## TODO

[ ] improve error handling: broken IPC results in malfunction of the whole thing
[ ] fix weird zombie leaking: running `go test .` leaks zombie processes
    inspect: `\ps -o ppid,state,command|grep 'Z '`
    kill: `\ps -o ppid,state,command|grep 'Z '|cut -f1 -d' '|xargs kill`


[Go]: http://golang.org
[gorack]: http://github.com/gmarik/gorack
[nack]: http://github.com/josh/nack
[Rack]: http://rack.github.io
[nodejs]: http://nodejs.org
