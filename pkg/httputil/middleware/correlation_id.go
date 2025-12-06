package middleware

import (
	"context"
	"net/http"

	"github.com/apotourlyan/ludus-studii/pkg/httputil/header"
	"github.com/apotourlyan/ludus-studii/pkg/idutil"
	"github.com/apotourlyan/ludus-studii/pkg/typeutil"
)

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get correlation ID from header
		correlationID := r.Header.Get(header.CorrelationID)

		// If not provided, generate one
		if correlationID == "" {
			correlationID = idutil.UUID()
		}

		// Put in context
		ctx := context.WithValue(r.Context(), typeutil.KeyCorrelationId, correlationID)

		// Echo back in response (useful for debugging)
		w.Header().Set(header.CorrelationID, correlationID)

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
