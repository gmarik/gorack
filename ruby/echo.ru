# vim: set ft=ruby:

run Proc.new {|env|
  IO.copy_stream(env['rack.input'], body = StringIO.new)
  except = [/rack\./]
  env2 = env.delete_if {|k, v| except.any? {|ex| ex === k}}
  [200, {}, [env.inspect, body.string] ]
}
