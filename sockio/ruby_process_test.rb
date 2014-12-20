#!/usr/bin/env ruby
require 'socket'

def log(msg)
	puts "[RUBY] #{msg}"
end

# passed from parent process
sock = UNIXSocket.for_fd(3)

log "receiving socket"
r = sock.recv_io

log "creating proxy pipe"
ior, iow = IO.pipe

log "sending the pipe"
sock.send_io(ior)

log "copying stream"
IO.copy_stream(r, iow)
