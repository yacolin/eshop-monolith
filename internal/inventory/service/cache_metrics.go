package service

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var registerOnce sync.Once

type cacheMetrics struct {
	l1Hit             prometheus.Counter
	l1Miss            prometheus.Counter
	l2Hit             prometheus.Counter
	l2Miss            prometheus.Counter
	dbFallback        prometheus.Counter
	bloomReject       prometheus.Counter
	evictionTotal     prometheus.Counter
	singleflightDedup prometheus.Counter
	lookupDuration    *prometheus.HistogramVec
}

var defaultMetrics *cacheMetrics

func getCacheMetrics() *cacheMetrics {
	if defaultMetrics == nil {
		defaultMetrics = newCacheMetrics()
		registerOnce.Do(func() {
			defaultMetrics.register()
		})
	}
	return defaultMetrics
}

func newCacheMetrics() *cacheMetrics {
	return &cacheMetrics{
		l1Hit: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_l1_hit_total",
			Help: "Total number of L1 (local LRU) cache hits",
		}),
		l1Miss: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_l1_miss_total",
			Help: "Total number of L1 (local LRU) cache misses",
		}),
		l2Hit: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_l2_hit_total",
			Help: "Total number of L2 (Redis) cache hits",
		}),
		l2Miss: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_l2_miss_total",
			Help: "Total number of L2 (Redis) cache misses",
		}),
		dbFallback: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_db_fallback_total",
			Help: "Total number of DB fallbacks when Redis is unavailable",
		}),
		bloomReject: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_bloom_reject_total",
			Help: "Total number of requests rejected by bloom filter",
		}),
		evictionTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_eviction_total",
			Help: "Total number of cache evictions (write invalidations)",
		}),
		singleflightDedup: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_singleflight_dedup_total",
			Help: "Total number of duplicate requests deduplicated by singleflight",
		}),
		lookupDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cache_lookup_duration_seconds",
			Help:    "Latency of cache lookup by layer",
			Buckets: prometheus.DefBuckets,
		}, []string{"layer"}),
	}
}

func (c *cacheMetrics) register() {
	prometheus.MustRegister(
		c.l1Hit, c.l1Miss, c.l2Hit, c.l2Miss, c.dbFallback,
		c.bloomReject, c.evictionTotal, c.singleflightDedup,
		c.lookupDuration,
	)
}

func (c *cacheMetrics) incL1Hit()      { c.l1Hit.Inc() }
func (c *cacheMetrics) incL1Miss()     { c.l1Miss.Inc() }
func (c *cacheMetrics) incL2Hit()      { c.l2Hit.Inc() }
func (c *cacheMetrics) incL2Miss()     { c.l2Miss.Inc() }
func (c *cacheMetrics) incDBFallback() { c.dbFallback.Inc() }
func (c *cacheMetrics) incBloomReject() { c.bloomReject.Inc() }
func (c *cacheMetrics) incEviction()   { c.evictionTotal.Inc() }
func (c *cacheMetrics) incSingleflightDedup() { c.singleflightDedup.Inc() }

func (c *cacheMetrics) observeDuration(layer string, duration float64) {
	c.lookupDuration.WithLabelValues(layer).Observe(duration)
}