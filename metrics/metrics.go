package metrics

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(prefix string) {
	namePrefix = prefix
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
}

type counterCollection struct {
	counters map[string]prometheus.Counter
	mu       sync.RWMutex
}

var (
	cCollection = counterCollection{
		counters: make(map[string]prometheus.Counter),
	}
	namePrefix = ""
)

// RegisterCounter will store a prometheus counter to be used later.
// the parameters of name means the counter's name.
// the parameters of help means the tips for the counter so that it can improve the readibility.
func RegisterCounter(name, help string) {
	cCollection.mu.Lock()
	defer cCollection.mu.Unlock()
	cCollection.counters[name] = promauto.NewCounter(prometheus.CounterOpts{
		Name: namePrefix + "_" + name,
		Help: help,
	})
}

// Incr will increase a feature report to the prometheus
func Incr(name string) {
	cCollection.mu.RLock()
	if v, ok := cCollection.counters[name]; ok {
		v.Inc()
		cCollection.mu.RUnlock()
		return
	}
	cCollection.mu.RUnlock()

	// add the new counter to the collection and increase
	c := promauto.NewCounter(prometheus.CounterOpts{
		Name: namePrefix + "_" + name,
	})
	cCollection.mu.Lock()
	defer cCollection.mu.Unlock()
	cCollection.counters[name] = c
	c.Inc()
}
