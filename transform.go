package wire

type Transformer interface {
	Transform(msg interface{}) (interface{}, error)
}
