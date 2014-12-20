#!/usr/bin/env ruby
require 'socket'

def log(msg)
	puts "[RUBY] #{msg}"
end

# passed from parent process
sock = UNIXSocket.for_fd(3)

log "receiving socket"
remote_r = sock.recv_io

log "creating proxy pipe"
local_r, local_w = IO.pipe

log "sending the pipe"
sock.send_io(local_r)

log "copying stream"
IO.copy_stream(remote_r, local_w)
