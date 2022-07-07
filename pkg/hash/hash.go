package hash

import (
	"crypto/md5"
	"fmt"
)

func CalculateChecksum(data []int32) string {
	hash := md5.New()
	for _, v := range data {
		hash.Write([]byte(fmt.Sprintf("%d", v)))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
