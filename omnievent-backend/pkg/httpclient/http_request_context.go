package httpclient

const (
	logHandleKey = "log_handler"
)

// HttpResponseLogHandlerFunc represents the http response log handler function
type HttpResponseLogHandlerFunc func([]byte)

// CustomHttpResponseLog returns a context with http response log handler
func CustomHttpResponseLog(c any, responseLogHandler HttpResponseLogHandlerFunc) any {
	return &httpRequestContext{
		context:   c,
		logHandler: responseLogHandler,
	}
}

// httpRequestContext represents the context for http request
type httpRequestContext struct {
	context   any
	logHandler HttpResponseLogHandlerFunc
}

// Value returns the value associated with key
func (c *httpRequestContext) Value(key any) any {
	if key == logHandleKey {
		return c.logHandler
	}

	return nil
}