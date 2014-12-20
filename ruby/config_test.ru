resp = "hellozzzz"

headers = {
  'X-This' => 'a messsage',
  'Content-Length' => resp.size,
}


run lambda {|env| [201, headers, [resp]]}
