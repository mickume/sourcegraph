package documents

import (
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/documents/internal/store"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/trace"
	"github.com/sourcegraph/sourcegraph/lib/log"
)

var (
	svc     *Service
	svcOnce sync.Once
)

// GetService creates or returns an already-initialized documents service. If the service is
// new, it will use the given database handle.
func GetService(db database.DB) *Service {
	svcOnce.Do(func() {
		observationContext := &observation.Context{
			Logger:     log.Scoped("documents.service", "codeintel documents service"),
			Tracer:     &trace.Tracer{Tracer: opentracing.GlobalTracer()},
			Registerer: prometheus.DefaultRegisterer,
		}

		svc = newService(
			store.GetStore(db),
			observationContext,
		)
	})

	return svc
}

// TestService creates a fresh documents service with the given database handle.
func TestService(db database.DB) *Service {
	return newService(
		store.GetStore(db),
		&observation.TestContext,
	)
}
