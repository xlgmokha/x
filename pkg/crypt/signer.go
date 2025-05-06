package crypt

type Signer interface {
	Sign([]byte) ([]byte, error)
}
