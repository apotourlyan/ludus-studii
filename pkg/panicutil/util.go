package panicutil

import (
	"fmt"

	"github.com/apotourlyan/ludus-studii/pkg/stringutil"
	"github.com/apotourlyan/ludus-studii/pkg/typeutil"
)

func RequireNotNil(pval any, pname string) {
	if pval == nil {
		panic(fmt.Sprintf("%q must not be nil", pname))
	}
}

func RequireNotEmptyOrWhitespace(pval string, pname string) {
	if len(pval) == 0 {
		panic(fmt.Sprintf("%q must not be empty", pname))
	}

	if stringutil.IsWhitespace(pval) {
		panic(fmt.Sprintf("%q must not be whitespace", pname))
	}
}

func RequireNonNegative[T typeutil.Numeric](pval T, pname string) {
	if pval < 0 {
		panic(fmt.Sprintf("%q must be >= 0, got %v", pname, pval))
	}
}

func RequirePositive[T typeutil.Numeric](pval T, pname string) {
	if pval <= 0 {
		panic(fmt.Sprintf("%q must be > 0, got %v", pname, pval))
	}
}

func RequireEqualTo[T typeutil.Numeric](pval T, limit T, pname string) {
	if pval != limit {
		panic(fmt.Sprintf("%q must be == %v, got %v", pname, limit, pval))
	}
}

func RequireLessThan[T typeutil.Numeric](pval T, limit T, pname string) {
	if pval >= limit {
		panic(fmt.Sprintf("%q must be < %v, got %v", pname, limit, pval))
	}
}

func RequireLessThanOrEqualTo[T typeutil.Numeric](pval T, limit T, pname string) {
	if pval > limit {
		panic(fmt.Sprintf("%q must be <= %v, got %v", pname, limit, pval))
	}
}

func RequireGreaterThan[T typeutil.Numeric](pval T, limit T, pname string) {
	if pval <= limit {
		panic(fmt.Sprintf("%q must be > %v, got %v", pname, limit, pval))
	}
}

func RequireGreaterThanOrEqualTo[T typeutil.Numeric](pval T, limit T, pname string) {
	if pval < limit {
		panic(fmt.Sprintf("%q must be >= %v, got %v", pname, limit, pval))
	}
}
