package short

import (
	"crypto/md5"
	"encoding/hex"
)

func URL(url []byte) string {
	hash := md5.Sum(url)

	return hex.EncodeToString(hash[:])[:8]
}
