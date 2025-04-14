package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// Requests older than this will be rejected
	maxRequestAge = 5 * time.Minute
	// How long to keep nonces in memory to prevent replay attacks
	nonceExpiration = 10 * time.Minute
)

type SecurityMiddleware struct {
	secretKey []byte
	nonces    sync.Map
}

func NewSecurityMiddleware(secretKey string) *SecurityMiddleware {
	return &SecurityMiddleware{
		secretKey: []byte(secretKey),
	}
}

// cleanupNonces periodically removes expired nonces
func (s *SecurityMiddleware) CleanupNonces() {
	ticker := time.NewTicker(nonceExpiration)
	go func() {
		for range ticker.C {
			now := time.Now()
			s.nonces.Range(func(key, value interface{}) bool {
				if timestamp, ok := value.(time.Time); ok {
					if now.Sub(timestamp) > nonceExpiration {
						s.nonces.Delete(key)
					}
				}
				return true
			})
		}
	}()
}

// ValidateRequest ensures requests are properly signed and not replayed
func (s *SecurityMiddleware) ValidateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		signature := c.GetHeader("X-Signature")

		if timestamp == "" || nonce == "" || signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing security headers"})
			c.Abort()
			return
		}

		// Validate timestamp is recent
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
			c.Abort()
			return
		}

		requestTime := time.Unix(ts, 0)
		if time.Since(requestTime) > maxRequestAge {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request expired"})
			c.Abort()
			return
		}

		// Check for replay attacks using nonce
		if _, exists := s.nonces.LoadOrStore(nonce, time.Now()); exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate request"})
			c.Abort()
			return
		}

		// Validate HMAC signature
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		// Important: Set the body back so it can be read again by later handlers
		c.Request.Body = NewBodyReader(body)

		// Calculate expected signature
		message := fmt.Sprintf("%s:%s:%s", timestamp, nonce, string(body))
		expectedSignature := s.calculateHMAC(message)

		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *SecurityMiddleware) calculateHMAC(message string) string {
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// NewBodyReader creates a new io.ReadCloser from a byte slice
func NewBodyReader(body []byte) *BodyReader {
	return &BodyReader{
		Reader: body,
		pos:    0,
	}
}

type BodyReader struct {
	Reader []byte
	pos    int
}

func (b *BodyReader) Read(p []byte) (n int, err error) {
	if b.pos >= len(b.Reader) {
		return 0, http.ErrBodyReadAfterClose
	}
	n = copy(p, b.Reader[b.pos:])
	b.pos += n
	return n, nil
}

func (b *BodyReader) Close() error {
	return nil
}
