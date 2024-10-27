package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"slices"

	"github.com/gin-gonic/gin"
)

// func (crhl CryptoHelper) HashSHA256(msg []byte) []byte {

// 	h := sha256.New()

// 	h.Write(msg)

// 	dst := h.Sum(nil)

// 	fmt.Printf("%x", dst)
// 	return dst
// }

type cryptoWriter struct {
	gin.ResponseWriter
	key string
	rb  []byte
}

func (cw *cryptoWriter) WriteHeader(statusCode int) {
	cw.Header().Set("HashSHA256", string(SignSHA256(cw.rb, cw.key)))
	cw.ResponseWriter.WriteHeader(statusCode)
}

// func (cw *cryptoWriter) Close() error {
// 	if err := cw.Close(); err != nil {
// 		return err
// 	}
// 	return cw.Close()
// }

// func SignSHA256Agent(msg []byte, key string) string {
// 	fmt.Println("message", string(msg))
// 	hm := hmac.New(sha256.New, []byte(key))
// 	size, err := hm.Write(msg)
// 	if err != nil {
// 		fmt.Println(err, size)
// 	}

// 	hash := hm.Sum(nil)
// 	// fmt.Println("Hash length:", len(hash), hash)               // Should be 32 bytes
// 	// fmt.Println("Hash (hex):", hex.EncodeToString(hash)) // Should be 64 hex characters
// 	// limiting hash to only include 32 symbols somehow fixes test results WTF?
// 	return hex.EncodeToString(hash)
// }

func SignSHA256(msg []byte, key string) string {
	fmt.Println("message", string(msg), key)
	hm := hmac.New(sha256.New, []byte(key))
	size, err := hm.Write([]byte(msg))
	if err != nil {
		fmt.Println(err, size)
	}

	hash := hm.Sum(nil)
	fmt.Println("Hash length:", len(hash), hash)               // Should be 32 bytes
	fmt.Println("Hash (hex):", hex.EncodeToString(hash)) // Should be 64 hex characters
	return hex.EncodeToString(hash)
}

func CryptoHandler(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// fmt.Println("in middleware for crypto")
		allowedMethods := []string{"GET", "DELETE", "OPTIONS"}
		if slices.Contains(allowedMethods, c.Request.Method) {
			c.Next()
			return
		} else if key != "" {
			hashHeader := c.Request.Header.Get("HashSHA256")

			if hashHeader != "" {
				//Read body
				bodyBytes, err := io.ReadAll(c.Request.Body)
				// fmt.Println("body content", string(bodyBytes))
				if len(bodyBytes) == 0 {
					c.Next()
					return
				}
				if err != nil {
					fmt.Println("had an issue reading request body")
				}
				//Reset body
				signServer := SignSHA256(bodyBytes, key)
				fmt.Println("signServer", signServer, hashHeader)
				if signServer == hashHeader {
					fmt.Println("equal", signServer, hashHeader)
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				//Compare signatures
				if signServer != hashHeader {
					c.AbortWithError(400, errors.New("signatures don't match"))
					return
				}
				//Add header to all responses before sending request
				cw := &cryptoWriter{c.Writer, key, bodyBytes}
				c.Writer = cw

			} else {
				c.Next()
				return
			}
		}

		c.Next()
	}
}
