require 'fcntl'
require 'socket'
require 'stringio'
require 'rack'
require 'rack/builder'

module Gorack
  class Rack
    def self.run(*args)
      s = new(*args)
      # s.handle
    end

    attr_accessor :config, :app, :file
    attr_accessor :ppid, :server, :heartbeat
    attr_accessor :io

    def initialize(config, options = {})
      self.config = config
      self.file   = options[:file]
      self.ppid   = Process.ppid

      @io = IO.open(options[:io].to_i)
      puts "Reading from #{io.fileno}"

      puts io.gets

      # at_exit { close }

      trap('TERM') { exit }
      trap('INT')  { exit }
      trap('QUIT') { close }

    end

    def load_config
      cfgfile = File.read(config)
      eval("Rack::Builder.new {( #{cfgfile}\n )}.to_app", TOPLEVEL_BINDING, config)
    end

    def load_json
      load_json!
    rescue LoadError
    end

    def load_json!
      require 'json' unless defined? ::JSON
    end

    def handle
      status  = 500
      headers = { 'Content-Type' => 'text/html' }
      body    = ["Internal Server Error"]

      env, input = nil, StringIO.new
      input.set_encoding('ASCII-8BIT') if input.respond_to?(:set_encoding)

      env = ::JSON.parse(io.read)

      input.rewind

      env = {
        "rack.version" => Rack::VERSION,
        "rack.input" => input,
        "rack.errors" => $stderr,
        "rack.multithread" => false,
        "rack.multiprocess" => true,
        "rack.run_once" => false,
        "rack.url_scheme" => ["yes", "on", "1"].include?(env["HTTPS"]) ? "https" : "http"
      }.merge(env)

      status, headers, body = app.call(env)
    end
  end
end
