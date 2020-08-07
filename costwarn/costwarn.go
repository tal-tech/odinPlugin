package costwarn

import (
	"time"

	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
)

type CostWarnPlugin struct {
	closed        bool
	costThreshold time.Duration
}

func (this *CostWarnPlugin) WrapCall(w *wrap.Wrap) error {
	if this.closed {
		return nil
	}
	start := time.Now()
	e := w.Next()
	end := time.Now()
	cost := end.Sub(start)
	ctx := w.GetCtx()
	if e != nil {
		logger.Ex(ctx, "CostLog", "call func:%s, cost:%v, remotehost:%v, logid:%v, error:%v", w.MethodTag, cost, ctx.Value("hostname"), ctx.Value("logid"), e)
	} else if cost > this.costThreshold {
		logger.Wx(ctx, "CostLog", "call func:%s, cost:%v>%v, remotehost:%v, logid:%v", w.MethodTag, cost, this.costThreshold, ctx.Value("hostname"), ctx.Value("logid"))
	} else {
		logger.Ix(ctx, "CostLog", "call func:%s, cost:%v<%v, remotehost:%v, logid:%v", w.MethodTag, cost, this.costThreshold, ctx.Value("hostname"), ctx.Value("logid"))
	}
	return e
}

func (this *CostWarnPlugin) UpdateConfig(configs config.PluginConfig) {
	if configs.CWConfigs.CostThreshold <= 0 {
		configs.CWConfigs.CostThreshold = 100
	}
	this.closed = configs.CWConfigs.Closed
	this.costThreshold = time.Duration(configs.CWConfigs.CostThreshold) * time.Millisecond
	return
}

func InitCostWarn(configs config.PluginConfig) *CostWarnPlugin {
	plugin := new(CostWarnPlugin)
	if configs.CWConfigs.CostThreshold <= 0 {
		configs.CWConfigs.CostThreshold = 100
	}
	plugin.closed = configs.CWConfigs.Closed
	plugin.costThreshold = time.Duration(configs.CWConfigs.CostThreshold) * time.Millisecond
	return plugin
}
