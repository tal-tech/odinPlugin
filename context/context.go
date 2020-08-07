package context

import (
	"context"
	"strconv"
	"time"

	"github.com/smallnest/rpcx/share"
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/loggerX/logtrace"
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
)

type ContextPlugin struct {
	closed bool
}

func (this *ContextPlugin) WrapCall(w *wrap.Wrap) error {
	if this.closed {
		return nil
	}
	ctx := w.GetCtx()
	reqMeta, ok := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	if ok {
		for k, v := range reqMeta {
			ctx = context.WithValue(ctx, k, v)
		}
	}
	if ctx.Value("logid") == nil {
		var logidStr string
		logidStr = strconv.FormatInt(logger.Id(), 10)
		ctx = context.WithValue(ctx, "logid", logidStr)
	}
	ctx = logtrace.ExtractTraceNodeToXexContext(ctx)
	ctx = context.WithValue(ctx, "start", time.Now())
	w.SetCtx(ctx)
	return nil

}

func (this *ContextPlugin) UpdateConfig(configs config.PluginConfig) {
	this.closed = configs.CXConfigs.Closed
	return
}

func InitContext(configs config.PluginConfig) *ContextPlugin {
	plugin := new(ContextPlugin)
	plugin.closed = configs.CXConfigs.Closed
	return plugin
}
