package message

type Serializable interface {
	Serialize() ([]byte, error)
}

type Desirializable interface {
	Desirialize([]byte) error
}
