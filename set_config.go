package parse

import "time"

type Config struct {
	TimeFormat     string
	SliceSeparator string
}

var DefaultTimeFormat = time.RFC3339
var DefaultSliceSeparator = ","

func NewConfig(opts ...Option) *Config {
	c := &Config{
		TimeFormat:     DefaultTimeFormat,
		SliceSeparator: DefaultSliceSeparator,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Option func(c *Config)

func WithTimeFormat(timeFormat string) Option {
	return func(c *Config) {
		c.TimeFormat = timeFormat
	}
}

func WithSliceSeparator(sliceSeparator string) Option {
	return func(c *Config) {
		c.SliceSeparator = sliceSeparator
	}
}
