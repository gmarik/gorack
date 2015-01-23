#!/usr/bin/env gem build
# encoding: utf-8

Gem::Specification.new do |s|
  s.name = "gorack"
  s.version = "0.0.1"
  s.authors = ["http://github.com/gmarik"]
  s.homepage = "http://github.com/gmarik/example"
  s.summary = "Golang HTTP frontend for Ruby's Rack Applications"
  s.description = "#{s.summary}. "
  # s.cert_chain = nil
  s.email = "Z21hcmlrQGdtYWlsLmNvbQ==\n".unpack('m').first
  s.has_rdoc = false

  # files
  s.files = `find . -type f`.split("\n")

  s.executables = Dir["bin/*"].map(&File.method(:basename))
  # s.default_executable = "example"
  s.require_paths = ["lib"]

  # Ruby version
  s.required_ruby_version = ::Gem::Requirement.new("~> 2.0")

  # dependencies
  s.add_dependency('rack')
  # development dependencies (add_development_dependency)
end
