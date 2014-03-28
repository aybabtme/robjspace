# rubyobj

Provides an encoder and a decoder for Ruby ObjectSpace heap dumps.

# Docs

[Godoc](http://godoc.org/github.com/aybabtme/rubyobj)!

# Usage

```
go get github.com/aybabtme/robjspace
```

Then in your Go code:

```go
r := yourFavoriteReader() // say a file, or stdin

rubyObj := rubyobj.RubyObject{}
var err error

for dec := rubyobj.NewDecoder(r); err == nil; err = dec.Decode(&rubyObj) {
  fmt.Printf("%v\n", &rubyObj)
}

if err != io.EOF {
  perror(err)
}
```
