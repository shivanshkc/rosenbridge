package miscutils

// StringSliceContains returns true if "value" is in "slice".
func StringSliceContains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}
