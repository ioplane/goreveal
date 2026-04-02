package main

import (
	"fmt"
	"os"

	"example.com/protectedfixture/internal/features"
	"example.com/protectedfixture/internal/licensegate"
)

var protectedSentinel = []byte("enterprise-gated-goreveal-sample\x00")

//go:noinline
func readLicenseToken(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	return "trial-token"
}

//go:noinline
func auditFeatureGate(token string) string {
	if licensegate.VerifyLicenseToken(token) {
		return features.AuditFeatureFlag("enterprise-sync")
	}
	return features.AuditFeatureFlag("basic")
}

//go:noinline
func runEnterpriseReport(token string) string {
	enabled := licensegate.VerifyLicenseToken(token)
	return features.EnterpriseReport(enabled)
}

func main() {
	token := readLicenseToken(os.Args)
	fmt.Println(
		auditFeatureGate(token),
		runEnterpriseReport(token),
		licensegate.DigestToken(token)[:8],
		len(protectedSentinel),
	)
}
