package tracev2

type Span interface {
	Options() SpanOptions
	SpanData() *SpanData
	AddMetadata(md map[string]string)
	Finish()
	String() string
}
