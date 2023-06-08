package executors

import (
	"context"
	"strconv"
	"time"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/rcache"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type multiqueueCacheCleaner struct {
	queueNames []string
	cache      *rcache.Cache
	windowSize time.Duration
	logger     log.Logger
}

var _ goroutine.Handler = &multiqueueCacheCleaner{}

// NewMultiqueueCacheCleaner returns a PeriodicGoroutine that will check the cache for entries that are older than the configured
// window size. A cache key is represented by a queue name; the value is a hash containing timestamps as the field key and the
// job ID as the field value (which is not used for anything currently).
func NewMultiqueueCacheCleaner(queueNames []string, cache *rcache.Cache, windowSize time.Duration, cleanupInterval time.Duration) goroutine.BackgroundRoutine {
	ctx := context.Background()
	return goroutine.NewPeriodicGoroutine(
		ctx,
		&multiqueueCacheCleaner{
			queueNames: queueNames,
			cache:      cache,
			windowSize: windowSize,
			logger:     log.Scoped("executors.multiqueue-cache-cleaner", "deletes entries from the dequeue cache older than the configured window"),
		},
		goroutine.WithName("executors.multiqueue-cache-cleaner"),
		goroutine.WithDescription("deletes entries from the dequeue cache older than the configured window"),
		goroutine.WithInterval(cleanupInterval),
	)
}

// Handle loops over the configured queue names and deletes stale entries.
func (m *multiqueueCacheCleaner) Handle(ctx context.Context) error {
	for _, queueName := range m.queueNames {
		all, err := m.cache.GetHashAll(queueName)
		if err != nil {
			return errors.Wrap(err, "multiqueue.cachecleaner")
		}
		for key := range all {
			timestampMillis, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				return err
			}
			t := time.Unix(0, timestampMillis*int64(time.Millisecond))
			interval := time.Now().Add(-m.windowSize)
			m.logger.Info("checking key", log.String("key", key), log.Int64("timestampMillis", timestampMillis), log.Time("t", t), log.Time("interval", interval))
			if t.Before(interval) {
				// expired cache entry, delete
				deletedItems, err := m.cache.DeleteHashItem(queueName, key)
				m.logger.Info("deleted key", log.String("key", key), log.Int("deletedItems", deletedItems))
				if err != nil {
					return err
				}
				if deletedItems == 0 {
					return errors.Newf("failed to delete hash item %s for key %s: expected successful delete but redis deleted nothing", key, queueName)
				}
			}
		}
	}
	return nil
}