package http

import (
	"github.com/cespare/xxhash/v2"
	"net/url"
	"strings"
	"crypto/rand"
	"fmt"
	"time"
	"os"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeBase62(num uint64) string {
	if num == 0 {
		return "0"
	}

	encoded := ""
	base := uint64(len(base62Chars))

	for num > 0 {
		remainder := num % base
		encoded = string(base62Chars[remainder]) + encoded
		num = num / base
	}

	return encoded
}

func ShortHash(url string) string {
	hash := xxhash.Sum64String(url)
	return EncodeBase62(hash)
}

func Normalize(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Host = strings.ToLower(u.Host)
	return u.String()
}


func generateRandomSuffix(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano()%1000)
	}
	
	for i := range b {
		b[i] = base62Chars[int(b[i])%len(base62Chars)]
	}
	return string(b)
}

func getShortURLBase() string {
	if base := os.Getenv("SHORT_URL_BASE"); base != "" {
		return base
	}
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	if host == "" || host == "0.0.0.0" {
		host = "localhost"
	}
	if port == "" {
		port = "8080"
	}
	return "http://" + host + ":" + port
}