package app

import "net/http"

// MethodFilter.
type MethodFilter struct {
	Method   string
	Handler  http.Handler
	Fallback http.Handler
}

// Create new MethodFilter.
func NewMethodFilter(method string, handler http.Handler, fallback http.Handler) *MethodFilter {
	var filter = new(MethodFilter)
	filter.Method = method
	filter.Handler = handler
	filter.Fallback = fallback
	return filter
}

func (self *MethodFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	if r.Method == self.Method {
		handler = self.Handler
	} else {
		handler = self.Fallback
	}
	handler.ServeHTTP(w, r)
}
