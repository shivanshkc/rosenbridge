package core

// codeAndReasonFromErr converts an error to a CodeAndReason type.
func codeAndReasonFromErr(err error) *CodeAndReason {
	//nolint:errorlint
	cnr, ok := err.(*CodeAndReason)
	if !ok {
		return &CodeAndReason{Code: CodeUnknown, Reason: err.Error()}
	}
	return cnr
}
