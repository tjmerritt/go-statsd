
package statsd

import "time"

type BaseClient struct {
        config *Config
}

func NewBaseClient(config *Config) *BaseClient {
        return nil
}

func (s *BaseClient) Status() error {
        return nil
}

func mergeConfig(configs []*Config) *Config {
        cfg := &Config{
                Host: "localhost:8125",
                Interval: 60 * time.Second,
                Source: "",
                Metric: "",
                Shared: true,
                TagFormat: TAGS_NONE,
                SampleRate: time.Second,
        }
        return cfg
}

func NewClient(config... *Config) Client {
        cfg := mergeConfig(config)

        if cfg.Shared {
                return NewSharedClient(cfg)
        } else {
                return NewStandardClient(cfg)
        }
}

func getValue(value interface{}) float64 {
        switch value.(type) {
        case int8:
                return float64(value.(int8))
        case uint8:
                return float64(value.(uint8))
        case int16:
                return float64(value.(int16))
        case uint16:
                return float64(value.(uint16))
        case int32:
                return float64(value.(int32))
        case uint32:
                return float64(value.(uint32))
        case int64:
                return float64(value.(int64))
        case uint64:
                return float64(value.(uint64))
        case int:
                return float64(value.(int))
        case uint:
                return float64(value.(uint))
        case float32:
                return float64(value.(float32))
        case float64:
                return value.(float64)
        case time.Duration:
                return float64(value.(time.Duration).Nanoseconds()) / 1000000000.0
        default:
                return 0.0
        }
}
