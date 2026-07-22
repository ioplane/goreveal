package recoveryerr

import (
	"errors"
	"io/fs"
	"testing"
)

func TestUnavailableSupportsIsAsAndWrappedCause(t *testing.T) {
	t.Parallel()

	err := NewUnavailable(CodeBuildInfoNotFound, "build info is absent", fs.ErrNotExist)
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("errors.Is(%v, ErrUnavailable) = false", err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("errors.Is(%v, fs.ErrNotExist) = false", err)
	}

	var recoveryError *Error
	if !errors.As(err, &recoveryError) {
		t.Fatalf("errors.As(%v, *Error) = false", err)
	}
	if recoveryError.Kind != KindUnavailable {
		t.Fatalf("Error.Kind = %q, want %q", recoveryError.Kind, KindUnavailable)
	}
	if recoveryError.Code != CodeBuildInfoNotFound {
		t.Fatalf("Error.Code = %q, want %q", recoveryError.Code, CodeBuildInfoNotFound)
	}
	if recoveryError.Message != "build info is absent" {
		t.Fatalf("Error.Message = %q", recoveryError.Message)
	}
	if !errors.Is(recoveryError.Cause, fs.ErrNotExist) {
		t.Fatalf("Error.Cause = %v, want fs.ErrNotExist", recoveryError.Cause)
	}
}

func TestUnsupportedSupportsKindAndCodeMatching(t *testing.T) {
	t.Parallel()

	err := NewUnsupported(CodeRuntimeUnsupportedContainer, "Mach-O runtime metadata is unsupported", nil)
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("errors.Is(%v, ErrUnsupported) = false", err)
	}
	if !errors.Is(err, &Error{Kind: KindUnsupported}) {
		t.Fatalf("errors.Is(%v, unsupported Error) = false", err)
	}
	if !errors.Is(err, &Error{Kind: KindUnsupported, Code: CodeRuntimeUnsupportedContainer}) {
		t.Fatalf("errors.Is(%v, unsupported runtime Error) = false", err)
	}
	if errors.Is(err, &Error{Kind: KindUnavailable}) {
		t.Fatalf("errors.Is(%v, unavailable Error) = true", err)
	}
}
