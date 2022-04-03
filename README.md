# pcontext

[![GoDoc](https://pkg.go.dev/github.com/northbright/pcontext?status.svg)](https://pkg.go.dev/github.com/northbright/pcontext)

Package pcontext provides a new context which derived from context.Context.
The context creates a channel to send / receive progress data between work goroutine and other goroutines. The channel will be closed automatically when work done.

## Documentation
* [API Reference](http://godoc.org/github.com/northbright/pcontext)

## License
* [MIT License](./LICENSE)
