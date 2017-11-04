package statsd

import "time"

type TaggingFormat int

const (
	TAGS_NONE TaggingFormat = iota
	TAGS_LIBRATO
	TAGS_DATADOG
)

type MetricType int

const (
	TYPE_NONE MetricType = iota
	TYPE_COUNTER
	TYPE_GAUGE
	TYPE_TIMING
	TYPE_HISTOGRAM
	TYPE_SET
	TYPE_DOGEVENT
	TYPE_DOGSERVICECHECK
)

type Config struct {
	Host       string
	Interval   time.Duration
	Source     string
	Metric     string
	Shared     bool
	TagFormat  TaggingFormat
	Tags       map[string]string
	SampleRate time.Duration
	BufferSize int
	baseClient *BaseClient
}

type EventInfo struct {
	Title          string
	Text           string
	Timestamp      time.Time
	Hostname       string
	AggregationKey string
	Priority       string
	SourceType     string
	AlertType      string
	Tags           map[string]string
}

type ServiceCheckStatus int

const (
	STATUS_DEFAULT  ServiceCheckStatus = -1
	STATUS_OK       ServiceCheckStatus = 0 // Sucks that this is zero
	STATUS_WARNING  ServiceCheckStatus = 1
	STATUS_CRITICAL ServiceCheckStatus = 2
	STATUS_UNKNOWN  ServiceCheckStatus = 3
)

type ServiceCheckInfo struct {
	Name      string
	Status    ServiceCheckStatus
	Timestamp time.Time
	Hostname  string
	Message   string
	Tags      map[string]string
}

type Client interface {
	// Housekeeping
	Status() error
	Clone(config ...*Config) Client
	Close()

	// Counters
	Add(metric string, amount int, tags ...string)
	Increment(metric string, tags ...string)
	Decrement(metric string, tags ...string)
	Inc(metric string, tags ...string)
	Dec(metric string, tags ...string)

	// Gauges
	Gauge(metric string, value interface{}, tags ...string)
	GaugeMin(metric string, value interface{}, tags ...string)
	GaugeMax(metric string, value interface{}, tags ...string)

	// Timing
	Timing(metric string, interval interface{}, tags ...string)
	TimingMin(metric string, interval interface{}, tags ...string)
	TimingMax(metric string, interval interface{}, tags ...string)

	// Histogram
	Histogram(metric string, interval interface{}, tags ...string)
	HistogramMin(metric string, interval interface{}, tags ...string)
	HistogramMax(metric string, interval interface{}, tags ...string)

	// Sets
	Set(metric string, value interface{}, tags ...string)

	// Events [DataDog specific]
	Event(info ...*EventInfo)

	// Service Checks [DataDog specific]
	ServiceCheckInfo(info ...*ServiceCheckInfo)
}

// Wire formats:
//   Standard: metric:value|type@rate
//   Librato: metric--source#tag1=value1,tag2=value2:value|type@rate
//   DataDog: metric|type|@rate|#tag1:value1,tags2:value2
