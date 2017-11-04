package statsd

import (
	"strings"
)

type Metric struct {
	metricType MetricType
	value      interface{}
	count      int
	tags       map[string]string
}

type StandardClient struct {
	base    *BaseClient
	config  *Config
	metrics map[string]Metric
}

func NewStandardClient(config ...*Config) *StandardClient {
	return nil
}

func (s *StandardClient) CloneClient(config ...*Config) *StandardClient {
	return nil
}

func (s *StandardClient) Clone(config ...*Config) Client {
	return s.CloneClient(config...)
}

func (s *StandardClient) Status() error {
	return s.base.Status()
}

func (s *StandardClient) Close() {
}

func (s *StandardClient) GetType(metric string) MetricType {
	if m, ok := s.metrics[metric]; ok {
		return m.metricType
	}

	return TYPE_NONE
}

func (s *StandardClient) GetCount(metric string) int {
	if m, ok := s.metrics[metric]; ok {
		return m.count
	}

	return 0
}

func (s *StandardClient) GetMetric(metric string, metricType MetricType) (interface{}, bool) {
	if m, ok := s.metrics[metric]; ok {
		if m.metricType != metricType {
			return nil, false
		}

		return m.value, true
	}

	switch metricType {
	case TYPE_COUNTER:
		return int64(0), true
	case TYPE_GAUGE:
		return float64(0.0), true
	case TYPE_TIMING:
		return float64(0.0), true
	case TYPE_HISTOGRAM:
		return float64(0.0), true
	}

	return nil, true
}

func (s *StandardClient) GetCounter(metric string) (int64, bool) {
	if value, ok := s.GetMetric(metric, TYPE_COUNTER); ok {
		if v, ok := value.(int64); ok {
			return v, true
		}
	}
	return 0, false
}

func (s *StandardClient) SetMetric(metric string, metricType MetricType, value interface{}) {
	if m, ok := s.metrics[metric]; ok {
		if m.metricType == metricType {
			m.value = value
			m.count++
			s.metrics[metric] = m
		}
	}
	s.metrics[metric] = Metric{value: value, metricType: metricType, count: 1, tags: nil}
}

func (s *StandardClient) SetTags(metric string, tags []string) {
	if m, ok := s.metrics[metric]; ok && m.tags == nil {
		m.tags = make(map[string]string)
		for _, tag := range tags {
			parts := strings.Split(tag, "=")
			key := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}
			m.tags[key] = value
		}
		s.metrics[metric] = m
	}
}

func (s *StandardClient) Add(metric string, amount int, tags ...string) {
	if value, ok := s.GetCounter(metric); ok {
		s.SetMetric(metric, TYPE_COUNTER, value+int64(amount))
		s.SetTags(metric, tags)
	}
}

func (s *StandardClient) Increment(metric string, tags ...string) {
	s.Add(metric, 1, tags...)
}

func (s *StandardClient) Decrement(metric string, tags ...string) {
	s.Add(metric, -1, tags...)
}

func (s *StandardClient) Inc(metric string, tags ...string) {
	s.Add(metric, 1, tags...)
}

func (s *StandardClient) Dec(metric string, tags ...string) {
	s.Add(metric, -1, tags...)
}

func (s *StandardClient) GetValue(metric string, metricType MetricType) (float64, bool) {
	if value, ok := s.GetMetric(metric, metricType); ok {
		if v, ok := value.(float64); ok {
			return v, true
		}
	}
	return 0.0, false
}

// Gauges
func (s *StandardClient) avg(metric string, value float64, metricType MetricType, tags ...string) {
	if v, ok := s.GetValue(metric, metricType); ok {
		// Calculate running average
		count := s.GetCount(metric)
		countp := count + 1
		v *= float64(count) / float64(countp)
		v += value / float64(countp)
		s.SetMetric(metric, metricType, v)
		s.SetTags(metric, tags)
	}
}

func (s *StandardClient) min(metric string, value float64, metricType MetricType, tags ...string) {
	if v, ok := s.GetValue(metric, metricType); ok {
		count := s.GetCount(metric)
		if count == 0 || v > value {
			s.SetMetric(metric, metricType, value)
		}
		s.SetTags(metric, tags)
	}
}

func (s *StandardClient) max(metric string, value float64, metricType MetricType, tags ...string) {
	if v, ok := s.GetValue(metric, metricType); ok {
		count := s.GetCount(metric)
		if count == 0 || v < value {
			s.SetMetric(metric, metricType, value)
		}
		s.SetTags(metric, tags)
	}
}

func (s *StandardClient) Gauge(metric string, value interface{}, tags ...string) {
	s.avg(metric, getValue(value), TYPE_GAUGE, tags...)
}

func (s *StandardClient) GaugeMin(metric string, value interface{}, tags ...string) {
	s.min(metric, getValue(value), TYPE_GAUGE, tags...)
}

func (s *StandardClient) GaugeMax(metric string, value interface{}, tags ...string) {
	s.max(metric, getValue(value), TYPE_GAUGE, tags...)
}

// Timing
func (s *StandardClient) Timing(metric string, interval interface{}, tags ...string) {
	s.avg(metric, getValue(interval), TYPE_TIMING, tags...)
}

func (s *StandardClient) TimingMin(metric string, interval interface{}, tags ...string) {
	s.min(metric, getValue(interval), TYPE_TIMING, tags...)
}

func (s *StandardClient) TimingMax(metric string, interval interface{}, tags ...string) {
	s.max(metric, getValue(interval), TYPE_TIMING, tags...)
}

// Histogram
func (s *StandardClient) Histogram(metric string, interval interface{}, tags ...string) {
	s.avg(metric, getValue(interval), TYPE_HISTOGRAM, tags...)
}

func (s *StandardClient) HistogramMin(metric string, interval interface{}, tags ...string) {
	s.min(metric, getValue(interval), TYPE_HISTOGRAM, tags...)
}

func (s *StandardClient) HistogramMax(metric string, interval interface{}, tags ...string) {
	s.max(metric, getValue(interval), TYPE_HISTOGRAM, tags...)
}

// Sets
func (s *StandardClient) Set(metric string, value interface{}, tags ...string) {
}

// Events [DataDog specific]
func (s *StandardClient) Event(info ...*EventInfo) {
}

// Service Checks [DataDog specific]
func (s *StandardClient) ServiceCheckInfo(info ...*ServiceCheckInfo) {
}
