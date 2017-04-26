package manager

//Codec is
type Codec interface {
	Encode(*PluginStrategy) ([]byte, error)
	Decode([]byte) (*PluginStrategy, error)
}

//NewCodec is
func NewCodec(name string) Codec {
	return &JSONCodec{}
}

//JSONCodec is
type JSONCodec struct {
}

//Encode is
func (j *JSONCodec) Encode(p *PluginStrategy) ([]byte, error) {
	return []byte{}, nil
}

//Decode is
func (j *JSONCodec) Decode(b []byte) (*PluginStrategy, error) {
	return nil, nil
}
