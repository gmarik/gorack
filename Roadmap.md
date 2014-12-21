## gorack_server test
- create test to serve response(s)

## Logging
- create separate loggers: INFO, DEBUG1, DEBUG2
- use loggers in both go server and ruby server


## Get rid of Json in GoRack::Server
- use \0 terminated strings
- parse request similarly to rack_response parser in gorack_server.go

## Embed resources
- so there's no dependencies


## Properly reap child processes in GoRack::Server
- so there's no zombies

