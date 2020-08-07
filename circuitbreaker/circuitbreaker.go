package circuitbreaker

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/handy/breaker"
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
)

type CircuitbreakerPlugin struct {
	lock       *sync.Mutex
	breaker    breaker.Breaker
	breakerMap map[string]breaker.Breaker
}

var ErrCircuitBreaker = errors.New("server refused: circuitbreaker is open")

func (this *CircuitbreakerPlugin) WrapCall(w *wrap.Wrap) error {
	if this.breaker == nil {
		return nil
	}
	if !this.breaker.Allow() {
		return ErrCircuitBreaker
	}
	e := w.Next()
	if e != nil {
		this.breaker.Failure(time.Millisecond)
	} else {
		this.breaker.Success(time.Millisecond)
	}
	return e
}

func (this *CircuitbreakerPlugin) UpdateConfig(configs config.PluginConfig) {
	if !configs.CBConfigs.Closed && configs.CBConfigs.FailureRatio > 0 {
		ins := breaker.NewBreaker(configs.CBConfigs.FailureRatio)
		this.breaker = ins
	} else {
		this.breaker = nil
	}
}

func InitCircuitBreaker(configs config.PluginConfig) *CircuitbreakerPlugin {
	plugin := new(CircuitbreakerPlugin)
	if !configs.CBConfigs.Closed && configs.CBConfigs.FailureRatio > 0 {
		ins := breaker.NewBreaker(configs.CBConfigs.FailureRatio)
		plugin.breaker = ins
	}
	return plugin
}
