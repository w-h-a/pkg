package buffer

var (
	DefaultSize = 1000
)

type Buffer interface {
	Options() BufferOptions
	Put(v interface{})
	Get(n int) []*Entry
	String() string
}
