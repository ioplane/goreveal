package features

//go:noinline
func EnterpriseReport(enabled bool) string {
	if enabled {
		return "enterprise-report-enabled"
	}
	return "enterprise-report-disabled"
}

//go:noinline
func AuditFeatureFlag(flag string) string {
	if flag == "enterprise-sync" {
		return "flag:enterprise-sync"
	}
	return "flag:basic"
}
