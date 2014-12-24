
run Proc.new { |env|
  IO.copy_stream(env['rack.input'], sio = StringIO.new)

  headers = {
    'X-This' => 'a messsage',
  }

  [201, headers, [sio.string]]
}
