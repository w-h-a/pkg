package broker

type Message struct {
	Header map[string]string
	Body   []byte
}
