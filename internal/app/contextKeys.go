package app

type IPAddressKey struct {
	Name string
}

type correlationIDKey struct {
	Name string
}

var ContextCorrelationIDKey = &correlationIDKey{"correlationIDKey"}

type requestRefererKey struct {
	Name string
}

var RefererServerKey = &requestRefererKey{"requestRefererKey"}
