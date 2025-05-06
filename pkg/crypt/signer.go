package crypt

type Signer interface {
	Sign([]byte) ([]byte, error)
	Verify([]byte, []byte) bool
}
