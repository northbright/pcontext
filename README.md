# pcontext

Package pcontext provides a new context which derived from context.Context.
The context creates a channel to send / receive progress data between work goroutine and other goroutines. The channel will be closed automatically when work done.
