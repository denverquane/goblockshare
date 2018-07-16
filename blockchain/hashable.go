package blockchain

type Hashable interface {
	GetHash() []byte
}
