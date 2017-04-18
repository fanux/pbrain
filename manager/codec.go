package manager

//Codec is
type Codec interface {
	Encode(*PluginStrategy) ([]byte, error)
	Decode([]byte) (*PluginStrategy, error)
}

//NewCodec is
func NewCodec(name string) {
	return &JSONCodec{}
}

//JSONCodec is
type JSONCodec struct {
}

//Encode is
func (*JsonCodec) Encode(p *PluginStrategy) ([]byte, error) {
}

//Decode is
func (*JsonCodec) Decode(b []byte) (*PluginStrategy, error) {
}
