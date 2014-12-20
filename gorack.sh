#!/usr/bin/env ruby
# vim: set ft=ruby:

begin
  require 'rubygems'
  require 'bundler'
  require 'bundler/setup'
rescue => e
  $stderr.puts e.message
end

require_relative './gorack_server'

Gorack::Server.run(ARGV[0])



