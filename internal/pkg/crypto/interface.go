package crypto

type ICryptoHelper interface {
	SignSHA256(message []byte) []byte
}
