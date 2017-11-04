
package statsd

import (
        "sync"
)

type SharedClient struct {
        lock sync.Mutex
        base *StandardClient
}

func NewSharedClient(config... *Config) *SharedClient {
        return nil
}

func (s *SharedClient) CloneClient(config... *Config) *SharedClient {
        return nil
}

func (s *SharedClient) Clone(config... *Config) Client {
        return s.CloneClient(config...)
}

func (s *SharedClient) Status() error {
        return s.base.Status()
}

func (s *SharedClient) Close() {
}

func (s *SharedClient) GetType(metric string) MetricType {
        s.lock.Lock()
        defer s.lock.Unlock()
        return s.base.GetType(metric)
}

func (s *SharedClient) GetCount(metric string) int {
        s.lock.Lock()
        defer s.lock.Unlock()
        return s.base.GetCount(metric)
}

func (s *SharedClient) GetMetric(metric string, metricType MetricType) (interface{}, bool) {
        s.lock.Lock()
        defer s.lock.Unlock()
        return s.base.GetMetric(metric, metricType)
}

func (s *SharedClient) GetCounter(metric string) (int64,  bool) {
        s.lock.Lock()
        defer s.lock.Unlock()
        return s.base.GetCounter(metric)
}

func (s *SharedClient) SetMetric(metric string, metricType MetricType, value interface{}) {
        s.lock.Lock()
        defer s.lock.Unlock()
        s.base.SetMetric(metric, metricType, value)
}

func (s *SharedClient) SetTags(metric string, tags []string) {
        s.lock.Lock()
        defer s.lock.Unlock()
        s.base.SetTags(metric, tags)
}

func (s *SharedClient) Add(metric string, amount int, tags... string) {
        s.lock.Lock()
        defer s.lock.Unlock()
        s.base.Add(metric, amount, tags...)
}

func (s *SharedClient) Increment(metric string, tags... string) {
        s.Add(metric, 1, tags...)
}

func (s *SharedClient) Decrement(metric string, tags... string) {
        s.Add(metric, -1, tags...)
}

func (s *SharedClient) Inc(metric string, tags... string) {
        s.Add(metric, 1, tags...)
}

func (s *SharedClient) Dec(metric string, tags... string) {
        s.Add(metric, -1, tags...)
}

func (s *SharedClient) GetValue(metric string, metricType MetricType) (float64, bool) {
        s.lock.Lock()
        defer s.lock.Unlock()
        return s.base.GetValue(metric, metricType)
}

// Gauges
func (s *SharedClient) avg(metric string, value float64, metricType MetricType, tags... string) {
        s.lock.Lock()
        defer s.lock.Unlock()
        if v, ok := s.base.GetValue(metric, metricType); ok {
                // Calculate running average
                count := s.base.GetCount(metric)
                countp := count + 1
                v *= float64(count) / float64(countp)
                v += value / float64(countp)
                s.base.SetMetric(metric, metricType, v)
                s.base.SetTags(metric, tags)
        }
}

func (s *SharedClient) min(metric string, value float64, metricType MetricType, tags... string) {
        s.lock.Lock()
        defer s.lock.Unlock()
        if v, ok := s.base.GetValue(metric, metricType); ok {
                count := s.base.GetCount(metric)
                if count == 0 || v > value {
                        s.base.SetMetric(metric, metricType, value)
                }
                s.base.SetTags(metric, tags)
        }
}

func (s *SharedClient) max(metric string, value float64, metricType MetricType, tags... string) {
        s.lock.Lock()
        defer s.lock.Unlock()
        if v, ok := s.base.GetValue(metric, metricType); ok {
                count := s.base.GetCount(metric)
                if count == 0 || v < value {
                        s.base.SetMetric(metric, metricType, value)
                }
                s.base.SetTags(metric, tags)
        }
}

func (s *SharedClient) Gauge(metric string, value interface{}, tags... string) {
        s.avg(metric, getValue(value), TYPE_GAUGE, tags...)
}

func (s *SharedClient) GaugeMin(metric string, value interface{}, tags... string) {
        s.min(metric, getValue(value), TYPE_GAUGE, tags...)
}

func (s *SharedClient) GaugeMax(metric string, value interface{}, tags... string) {
        s.max(metric, getValue(value), TYPE_GAUGE, tags...)
}

// Timing
func (s *SharedClient) Timing(metric string, interval interface{}, tags... string) {
        s.avg(metric, getValue(interval), TYPE_TIMING, tags...)
}

func (s *SharedClient) TimingMin(metric string, interval interface{}, tags... string) {
        s.min(metric, getValue(interval), TYPE_TIMING, tags...)
}

func (s *SharedClient) TimingMax(metric string, interval interface{}, tags... string) {
        s.max(metric, getValue(interval), TYPE_TIMING, tags...)
}

// Histogram
func (s *SharedClient) Histogram(metric string, interval interface{}, tags... string) {
        s.avg(metric, getValue(interval), TYPE_HISTOGRAM, tags...)
}

func (s *SharedClient) HistogramMin(metric string, interval interface{}, tags... string) {
        s.min(metric, getValue(interval), TYPE_HISTOGRAM, tags...)
}

func (s *SharedClient) HistogramMax(metric string, interval interface{}, tags... string) {
        s.max(metric, getValue(interval), TYPE_HISTOGRAM, tags...)
}

// Sets
func (s *SharedClient) Set(metric string, value interface{}, tags... string) {
}

// Events [DataDog specific]
func (s *SharedClient) Event(info... *EventInfo) {
}

// Service Checks [DataDog specific]
func (s *SharedClient) ServiceCheckInfo(info... *ServiceCheckInfo) {
}
