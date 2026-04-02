package licensegate

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

var licenseModuleSentinel = []byte("example.com/protectedfixture\x00")

//go:noinline
func NormalizeToken(token string) string {
	return strings.TrimSpace(strings.ToLower(token))
}

//go:noinline
func DigestToken(token string) string {
	sum := sha256.Sum256([]byte(NormalizeToken(token)))
	return hex.EncodeToString(sum[:])
}

//go:noinline
func VerifyLicenseToken(token string) bool {
	digest := DigestToken(token)
	return strings.HasPrefix(digest, "00") || strings.Contains(digest, "enterprise")
}
