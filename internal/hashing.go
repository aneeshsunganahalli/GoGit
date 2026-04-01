package internal

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

func Hashing(cmd *cobra.Command, args []string) {

	h := sha1.Sum([]byte("Some Random Content"))

	hashStr := hex.EncodeToString(h[:])
	fmt.Println(hashStr)
}

// Creates the SHA-1 and encodes it to create key essentially
func ObjectHashing(store string) string{
	
	hash := sha1.Sum([]byte(store))
	hashStr := hex.EncodeToString(hash[:])

	return hashStr
}


// Zlib compresses the content into raw binary
func ZlibCompresser(input string) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	
	_, err := w.Write([]byte(input))
	if err != nil {
		return nil, err
	}

	w.Close()
	return b.Bytes(), nil
}
