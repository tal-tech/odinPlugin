package rpcxplugin

import (
	"github.com/tal-tech/odinPlugin/circuitbreaker"
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/context"
	"github.com/tal-tech/odinPlugin/costwarn"
	"github.com/tal-tech/odinPlugin/cpuidle"
	"github.com/tal-tech/odinPlugin/ratelimit"
	"github.com/tal-tech/odinPlugin/recovery"
	"github.com/tal-tech/odinPlugin/traceing"
	"github.com/tal-tech/odinPlugin/wrap"
)

type Options func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint)

var MiddlewareOptions []Options

func RatelimitOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		rate := ratelimit.InitRateLimit(conf)
		return rate, rate.WrapCall

	}
}

func CpuIdleOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		cpuidle := cpuidle.InitCpuIdle(conf)
		return cpuidle, cpuidle.WrapCall

	}
}

func CircuitBreakerOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		circu := circuitbreaker.InitCircuitBreaker(conf)
		return circu, circu.WrapCall

	}
}

func TraceOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		trace := traceing.InitTrace(conf)
		return trace, trace.WrapCall

	}
}

func CostWarnOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		cw := costwarn.InitCostWarn(conf)
		return cw, cw.WrapCall

	}
}

func ContextOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		ctx := context.InitContext(conf)
		return ctx, ctx.WrapCall

	}
}

func RecoveryOption() Options {
	return func(conf config.PluginConfig) (MiddleWare, wrap.EndPoint) {
		rc := recovery.InitRecovery(conf)
		return rc, rc.WrapCall

	}
}
