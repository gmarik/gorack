headers = {
  'X-This' => 'a messsage',
  'Content-Length' => 5,
}

run lambda {|env| [200, headers, ["hello"]]}
