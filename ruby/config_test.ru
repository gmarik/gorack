
run Proc.new { |env|
  body = env['rack.input']
  # STDERR.puts body.read(100)
  IO.copy_stream(body, sio = StringIO.new)

  headers = {
    'X-This' => 'a messsage',
  }

  [201, headers, [sio.string]]
}
