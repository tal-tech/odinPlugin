package cpuidle

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
	"github.com/toolkits/nux"
)

const historyCount int = 2

type CpuIdlePlugin struct {
	lock            *sync.RWMutex
	lowthreshold    float64
	highthreshold   float64
	interval        int
	weight          int
	idle            float64
	procStatHistory [historyCount]*nux.ProcStat
	stopchan        chan struct{}
	closed          bool
	chanuse         bool
}

var ErrCpuIdleLow = errors.New("server refused: cpuidle is low")

func (this *CpuIdlePlugin) WrapCall(w *wrap.Wrap) error {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.closed || this.weight >= this.interval || this.idle == 0 {
		return nil
	}
	if rand.Intn(this.interval) > this.weight {
		return ErrCpuIdleLow
	}
	return nil
}

func (this *CpuIdlePlugin) UpdateConfig(configs config.PluginConfig) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.closed = configs.CIConfigs.Closed
	this.lowthreshold = configs.CIConfigs.LowThreshold
	this.highthreshold = configs.CIConfigs.HighThreshold
	this.interval = cast.ToInt(this.highthreshold - this.lowthreshold)
	this.weight = cast.ToInt(this.idle - this.lowthreshold)
	if this.chanuse {
		this.stopchan <- struct{}{}
	}
	this.chanuse = false
	if !configs.CIConfigs.Closed && configs.CIConfigs.LowThreshold > 0 && configs.CIConfigs.LowThreshold < configs.CIConfigs.HighThreshold {
		go this.updateCpuStat(configs.CIConfigs.RefreshTime)
		this.chanuse = true
	}
}

func (this *CpuIdlePlugin) updateCpuStat(refreshTime int64) {
	var tickerTime time.Duration
	if refreshTime == 0 {
		tickerTime = time.Second * 15
	} else {
		tickerTime = time.Duration(refreshTime) * time.Second
	}
	t := time.NewTicker(tickerTime)
	for {
		select {
		case <-t.C:
			ps, err := nux.CurrentProcStat()
			if err != nil {
				continue
			}

			this.lock.Lock()
			for i := historyCount - 1; i > 0; i-- {
				this.procStatHistory[i] = this.procStatHistory[i-1]
			}

			this.procStatHistory[0] = ps
			if this.procStatHistory[1] != nil {
				dt := this.procStatHistory[0].Cpu.Total - this.procStatHistory[1].Cpu.Total
				if dt != 0 {
					invQuotient := 100.00 / float64(dt)
					this.idle = float64(this.procStatHistory[0].Cpu.Idle-this.procStatHistory[1].Cpu.Idle) * invQuotient
					this.weight = cast.ToInt(this.idle - this.lowthreshold)
				}
			}
			this.lock.Unlock()
		case <-this.stopchan:
			return
		}

	}
}

func InitCpuIdle(configs config.PluginConfig) *CpuIdlePlugin {
	plugin := new(CpuIdlePlugin)
	plugin.lock = new(sync.RWMutex)
	plugin.stopchan = make(chan struct{}, 0)
	plugin.closed = configs.CIConfigs.Closed
	plugin.lowthreshold = configs.CIConfigs.LowThreshold
	plugin.highthreshold = configs.CIConfigs.HighThreshold
	plugin.interval = cast.ToInt(plugin.highthreshold - plugin.lowthreshold)
	plugin.weight = plugin.interval
	if !configs.CIConfigs.Closed && configs.CIConfigs.LowThreshold > 0 && configs.CIConfigs.LowThreshold < configs.CIConfigs.HighThreshold {
		go plugin.updateCpuStat(configs.CIConfigs.RefreshTime)
		plugin.chanuse = true
	}
	return plugin
}
