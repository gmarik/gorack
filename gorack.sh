#!/usr/bin/env ruby
# vim: set ft=ruby:

require_relative './gorack'

puts ARGV.join(" ")

begin
Gorack::Rack.run(ARGV[0], {reader: ARGV[1], writer: ARGV[2]})
rescue => e
  puts "Error: #{e.message}\n #{e.backtrace[0]}"
  exit(1)
end



