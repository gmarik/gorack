#!/usr/bin/env ruby
# vim: set ft=ruby:

require_relative './gorack'

# puts ARGV.join(" ")

begin
Gorack::Rack.run(ARGV[0])
rescue => e
  puts "Error: #{e.message}\n #{e.backtrace[0]}"
  exit(1)
end



