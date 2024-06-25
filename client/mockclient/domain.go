package mockclient

type Response struct {
	Response interface{} `json:"response,omitempty"`
	Err      error       `json:"err,omitempty"`
}
