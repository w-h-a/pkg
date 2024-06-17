package runtime

type Service struct {
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	Port      int               `json:"port"`
	Version   string            `json:"version"`
	Address   string            `json:"address"`
	Metadata  map[string]string `json:"metadata"`
}
