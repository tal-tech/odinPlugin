package ratelimit

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
	"github.com/tal-tech/xtools/rateutil"
)

type RatelimitPlugin struct {
	lock       *sync.Mutex
	limiterMap map[string]*limiter
}

type limiter struct {
	maxdelayduration time.Duration
	instacnce        rateutil.RateInstance
}

var ErrRateLimit = errors.New("server refused: due to ratelimit")

func (this *RatelimitPlugin) WrapCall(w *wrap.Wrap) error {
	this.lock.Lock()
	tempMap := this.limiterMap
	this.lock.Unlock()
	key := w.MethodTag
	if limiter, ok := tempMap[key]; ok {
		if !limiter.check() {
			return ErrRateLimit
		}
	}
	return nil
}

func (this *RatelimitPlugin) UpdateConfig(configs config.PluginConfig) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.limiterMap = make(map[string]*limiter)
	for _, config := range configs.RLConfigs {
		if config.Closed {
			delete(this.limiterMap, config.Path)
			continue
		}
		li := new(limiter)
		li.maxdelayduration = time.Duration(config.MaxDelayTime) * time.Millisecond
		li.instacnce, _ = rateutil.RateInstanceFactory(rateutil.TokenBucket, config.Path+strconv.FormatInt(time.Now().Unix(), 10), config.Limit, config.Burst)
		this.limiterMap[config.Path] = li
	}
	return
}

func (l *limiter) check() bool {
	if l.maxdelayduration > time.Duration(0) {
		return l.instacnce.Allow(l.maxdelayduration, false)
	} else {
		return l.instacnce.TryAllow()
	}
}

func InitRateLimit(configs config.PluginConfig) *RatelimitPlugin {
	plugin := new(RatelimitPlugin)
	plugin.lock = new(sync.Mutex)
	plugin.limiterMap = make(map[string]*limiter)
	for _, config := range configs.RLConfigs {
		if config.Closed {
			continue
		}
		li := new(limiter)
		li.maxdelayduration = time.Duration(config.MaxDelayTime) * time.Millisecond
		li.instacnce, _ = rateutil.RateInstanceFactory(rateutil.TokenBucket, config.Path, config.Limit, config.Burst)
		plugin.limiterMap[config.Path] = li
	}
	return plugin
}
