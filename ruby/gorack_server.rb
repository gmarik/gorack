require 'fcntl'
require 'socket'
require 'stringio'
require 'rack'
require 'rack/builder'
require 'json'

module Gorack
  class Server
    def self.run(*args)
      log("Waiting for connections")
      s = new(*args)
      loop { s.handle }
    end

    def self.log(msg)
      STDOUT.puts("[Master] #{msg}")
    end

    attr_accessor :config, :app, :file
    attr_accessor :ppid, :server, :heartbeat
    attr_accessor :master_io

    def initialize(config, options = {})
      self.config = config
      @master_io = UNIXSocket.for_fd(3)
      @app = load_config
    end

    def log(msg)
      self.class.log(msg)
    end

    def load_config
      cfgfile = File.read(config)
      eval("Rack::Builder.new {( #{cfgfile}\n )}.to_app", TOPLEVEL_BINDING, config)
    end

    def handle
      reader, writer = master_io.recv_io, master_io.recv_io

      log("Connection received")

      status  = 500
      headers = { 'Content-Type' => 'text/html' }
      body    = ["Internal Server Error"]

      IO.copy_stream(reader, req = StringIO.new)

      env = ::JSON.parse(req.string)

      env = {
        "rack.version" => 1,
        "rack.input" => StringIO.new,
        "rack.errors" => $stderr,
        "rack.multithread" => false,
        "rack.multiprocess" => true,
        "rack.run_once" => false,
        "rack.url_scheme" => ["yes", "on", "1"].include?(env["HTTPS"]) ? "https" : "http"
      }.merge(env)

      status, headers, body = app.call(env)

      writer.write("#{status}\n")
      writer.write(headers.map {|k, v| "#{k}: #{v}"}.join("\n"))
      writer.write("\n\n")
      # TODO:
      # IO.copy_stream(body, @writer)
      writer.write(body.join)
      writer.close
      # puts 'Done'
    end
  end
end
