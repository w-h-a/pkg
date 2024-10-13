package tracev2

type Span interface {
	Options() SpanOptions
	AddMetadata(md map[string]string)
	Finish()
	String() string
}
