package pow

import (
	"crypto/sha1"
	"fmt"
)

const zeroByte = 48

type HashCash struct {
	Ver      int
	Bits     int
	Date     int64
	Resource string
	Rand     string
	Counter  int
}

func (h HashCash) PrepareToSend() string {
	return fmt.Sprintf("%d:%d:%d:%s:%s:%d", h.Ver, h.Bits, h.Date, h.Resource, h.Rand, h.Counter)
}

func CheckHash(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}

	return true
}

func sha1Hash(toBeHashed string) string {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(toBeHashed))
	byteSlice := sha1Hash.Sum(nil)
	return fmt.Sprintf("%x", byteSlice)
}
func (h HashCash) ComputeHashCash(maxNumberOfIteration int) (HashCash, error) {
	for h.Counter <= maxNumberOfIteration || maxNumberOfIteration <= 0 {
		header := h.PrepareToSend()
		hash := sha1Hash(header)
		if CheckHash(hash, h.Ver) {
			return h, nil
		}
		h.Counter++

	}

	return h, fmt.Errorf("maximum iterations reached")

}
