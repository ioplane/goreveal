// Package recoveryerr defines core-owned, machine-readable recovery outcomes.
package recoveryerr

import (
	"errors"
	"fmt"
)

// Kind distinguishes expected absence from a recognized unsupported surface.
type Kind string

const (
	KindUnavailable Kind = "unavailable"
	KindUnsupported Kind = "unsupported"
)

// Code is a stable machine-readable recovery outcome code.
type Code string

const (
	CodeBuildInfoNotFound                 Code = "build_info_not_found"
	CodeRuntimeUnsupportedContainer       Code = "runtime_unsupported_container"
	CodeRuntimeMetadataNotFound           Code = "runtime_metadata_not_found"
	CodePclntabUnsupportedContainer       Code = "pclntab_unsupported_container"
	CodePclntabNotFound                   Code = "pclntab_not_found"
	CodeDWARFUnsupportedContainer         Code = "dwarf_unsupported_container"
	CodeDWARFNotFound                     Code = "dwarf_not_found"
	CodeStringRegionsUnsupportedContainer Code = "string_regions_unsupported_container"
	CodeStringRegionsNotFound             Code = "string_regions_not_found"
	CodeStringCandidatesNotFound          Code = "string_candidates_not_found"
	CodeSourceTreeNotFound                Code = "source_tree_not_found"
)

var (
	// ErrUnavailable matches typed errors for artifacts that truthfully lack evidence.
	ErrUnavailable = errors.New("recovery evidence unavailable")
	// ErrUnsupported matches typed errors for identified surfaces outside the implemented matrix.
	ErrUnsupported = errors.New("recovery surface unsupported")
)

// Error carries a stable recovery kind and code while preserving its cause.
type Error struct {
	Kind    Kind
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string {
	message := e.Message
	if message == "" {
		message = string(e.Code)
	}
	if e.Cause == nil {
		return message
	}

	return fmt.Sprintf("%s: %v", message, e.Cause)
}

// Unwrap preserves the original reader error for errors.Is/errors.As.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is supports both package sentinels and typed kind/code matching.
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	if target == ErrUnavailable {
		return e.Kind == KindUnavailable
	}
	if target == ErrUnsupported {
		return e.Kind == KindUnsupported
	}

	typedTarget, ok := target.(*Error)
	if !ok || typedTarget.Kind != e.Kind {
		return false
	}

	return typedTarget.Code == "" || typedTarget.Code == e.Code
}

// NewUnavailable reports a proven absence of recoverable evidence.
func NewUnavailable(code Code, message string, cause error) error {
	return &Error{Kind: KindUnavailable, Code: code, Message: message, Cause: cause}
}

// NewUnsupported reports an identified format/version/architecture outside the implemented matrix.
func NewUnsupported(code Code, message string, cause error) error {
	return &Error{Kind: KindUnsupported, Code: code, Message: message, Cause: cause}
}
