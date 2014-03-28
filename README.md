# rubyobj

Provides an encoder and a decoder for Ruby ObjectSpace heap dumps.

# Docs

[Godoc](http://godoc.org/github.com/aybabtme/rubyobj)!

# Usage

```
go get github.com/aybabtme/robjspace
```

Then in your Go code (this shows the slow decoder):

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

# Performance

Using the fast `ParallelDecode`:

```bash
$ go run loadall.go parallel --filename ../testdata/huge.json
loading 549MB from 'huge.json'
2489364 heap objects in 4.699698607s
```
