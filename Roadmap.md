# Roadmap

## gorack_server test
[x] create test to serve response(s)
[x] submit request body as well
[ ] submit request headers (go) and parse them properly(rack)

## Logging
[x] use loggers in both go server and ruby server
- create separate loggers: INFO, DEBUG1, DEBUG2

## Get rid of Json in GoRack::Server
[x] get rid of Json serialisation
[x] parse request similarly to rack_response parser in gorack_server.go
[x] use \0 terminated strings

## Embed resources
- so there's no dependencies


## Properly reap child processes in GoRack::Server
- so there's no zombies

