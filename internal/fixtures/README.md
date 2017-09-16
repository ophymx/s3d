Adding/Removing XML Fixutres
============================

XML file assets used in tests are referenced via [go-bindata](https://github.com/jteeuwen/go-bindata).
After adding or removing a file, update `xml.go` with `go generate`.

```
go get -u github.com/jteeuwen/go-bindata/...
go generate
```
