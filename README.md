[![Go Reference](https://pkg.go.dev/badge/github.com/eientei/cookiejarx.svg)](https://pkg.go.dev/github.com/eientei/cookiejarx)

This package provides public version of `net/http/cookiejar` package from go 1.17.2 with all helper methods and
storage interface publicly exported.  

It aims to provide both pluggable [storage](https://pkg.go.dev/github.com/eientei/cookiejarx#Storage) for cookiejar 
implementations using this [cookie jar](https://pkg.go.dev/github.com/eientei/cookiejarx#Jar) and all required helpers 
to implement own `http.CookieJar` from scratch according to RFC 6265, including punycode hostname normalization.

By default, it uses same in-memory storage as original package, however when provided explicitly, it is possible to
leverage additional methods of in-memory storage to implement simple external saving/loading and clearing of in-memory
storage;

see
- [InMemoryStorage.EntriesDump](https://pkg.go.dev/github.com/eientei/cookiejarx#InMemoryStorage.EntriesDump)
- [InMemoryStorage.EntriesRestore](https://pkg.go.dev/github.com/eientei/cookiejarx#InMemoryStorage.EntriesRestore)
- [InMemoryStorage.EntriesClear](https://pkg.go.dev/github.com/eientei/cookiejarx#InMemoryStorage.EntriesClear)
