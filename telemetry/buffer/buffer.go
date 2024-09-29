package buffer

var (
	defaultSize = 1024
)

type Buffer interface {
	Options() BufferOptions
	Put(v interface{})
	Get(n int) []*Entry
	String() string
}
