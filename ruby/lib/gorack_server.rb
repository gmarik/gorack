require 'stringio'
require 'rack'
require 'rack/builder'

module Gorack
  class Server

    DELIM = "\0" # response/request delimiter

    def self.run(*args)
      log("Waiting for connections")
      s = new(*args)
      loop {
        s.accept do |reader, writer|
          Process.fork { s.handle(reader, writer) }
        end
      }
    end

    def self.log(msg)
      STDOUT.puts("[PID:#{Process.pid}] #{msg}")
    end

    attr_accessor :config, :app, :file
    attr_accessor :ppid, :server, :heartbeat
    attr_accessor :master_io

    def accept(&block)
      pipe = master_io.recv_io, master_io.recv_io
      if block
        block.call(*pipe)
        pipe.each(&:close)
      end
      pipe
    end

    def initialize(master_sock, config_file, options = {})
      @master_io = master_sock
      @app, @app_options = Rack::Builder.parse_file(config_file)
    end

    def log(msg)
      self.class.log(msg)
    end

    def handle(reader, writer)
      begin
        request_env, body_reader = read_request(reader)

        rack_env = {
          "rack.version" => 1,
          "rack.input" => body_reader,
          "rack.errors" => STDERR,
          "rack.multithread" => false,
          "rack.multiprocess" => true,
          "rack.run_once" => false,
          "rack.url_scheme" => ["yes", "on", "1"].include?(request_env["HTTPS"]) ? "https" : "http"
        }.merge(request_env)

        status, headers, body = app.call(rack_env)

      rescue => e
        log("ERROR: " + e.message)
        log("ERROR: " + e.backtrace.join("\n"))
        status  = 500
        headers = {'Content-Type' => 'text/plain' }
        body    = ["Internal Server Error"]
      end

      write_response(writer, [status, headers, body])
    end

private

    def write_response(writer, resp)
      status, headers, body = *resp

      writer.write("#{status}#{DELIM}")
      writer.write(headers.map {|k, v| "#{k}: #{v}#{DELIM}"}.join(''))
      writer.write(DELIM)
      # TODO: use IO.copy_stream
      body.each(&writer.method(:write))
      writer.close
    end

    def read_request(reader)
      eol = eoh = false
      buf = StringIO.new

      while not eoh do
        break if reader.eof?
        buf.write(char = reader.read(1))
        # TODO: write test for this
        eoh = eol && char == DELIM
        eol = char == DELIM
      end

      lines = buf.string.split(DELIM)
      env = Hash[*lines.flat_map {|l| l.split(": ", 2)}]

      # reader is at body start or EOF
      [env, reader]
    end

  end
end
