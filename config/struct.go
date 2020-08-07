package config

type PluginConfig struct {
	RLConfigs []RateLimitConfig    `json:"ratelimit"`
	CBConfigs CircuitBreakerConfig `json:"circuitbreaker"`
	MSConfigs MetricsConfig        `json:"metrics"`
	TRConfigs TraceConfig          `json:"trace"`
	CIConfigs CpuIdleConfig        `json:"cpuidle"`
	CWConfigs CostWarnConfig       `json:"costwarn"`
	CXConfigs ContextConfig        `json:"context"`
	RCConfigs RecoveryConfig       `json:"recovery"`
	MDConfigs map[string]string    `json:"metadata"`
}

type RateLimitConfig struct {
	Path         string `json:"path"`
	Limit        int    `json:"limit"`
	Burst        int    `json:"burst"`
	Closed       bool   `json:"closed"`
	MaxDelayTime int64  `json:"maxDelayTime"`
}

type CircuitBreakerConfig struct {
	FailureRatio float64 `json:"failureRatio"`
	Closed       bool    `json:"closed"`
}

type MetricsConfig struct {
	Closed bool   `json:"closed"`
	Prefix string `json:"prefix"`
}

type ContextConfig struct {
	Closed bool `json:"closed"`
}

type TraceConfig struct {
	Closed bool `json:"closed"`
}

type CpuIdleConfig struct {
	LowThreshold  float64 `json:"lowThreshold"`
	HighThreshold float64 `json:"highThreshold"`
	RefreshTime   int64   `json:"refreshTime"`
	Closed        bool    `json:"closed"`
}

type CostWarnConfig struct {
	CostThreshold int64 `json:"costThreshold"`
	Closed        bool  `json:"closed"`
}

type RecoveryConfig struct {
	Closed bool `json:"closed"`
}
