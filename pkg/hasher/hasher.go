package hasher

import (
	"crypto/sha1"
	"fmt"
)

type HasherWithSalt struct {
	salt []byte
}

func NewHasherWithSalt(salt []byte) *HasherWithSalt {
	return &HasherWithSalt{salt: salt}
}

func (h *HasherWithSalt) GenerateHash(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))

	return fmt.Sprintf("%x", hash.Sum(h.salt))
}
