require 'stringio'
require 'rack'
require 'rack/builder'

module Gorack
  class Server

    DELIM = "\0" # response/request delimiter

    attr_accessor :master_io, :app, :logger
    def initialize(master_sock, app, app_options, logger)
      @master_io = master_sock
      @app, @app_options = app, app_options
      @logger = logger
      @pids = []
    end

    def run
      logger.info("Accepting connections")
      loop {
        accept do |reader, writer|
          add_pid Process.fork { handle(reader, writer) }
        end

        reap
      }
    end

    def exit
      logger.info("Exiting")
      Process.waitall unless @pids.empty?
    end

    def accept(&block)
      pipe = master_io.recv_io, master_io.recv_io
      if block
        block.call(*pipe)
        pipe.each(&:close)
      end
      pipe
    end

    def add_pid(pid)
      @pids.push(pid)
    end

    def reap
      new_pids = []
      while pid = @pids.pop
        Process.waitpid(pid, Process::WNOHANG) or new_pids.push(pid)
      end
      @pids = new_pids
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
        logger.info("ERROR: " + e.message)
        logger.info("ERROR: " + e.backtrace.join("\n"))
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
