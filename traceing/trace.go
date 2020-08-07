package traceing

import (
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
	"github.com/tal-tech/xtools/traceutil"
)

type TracePlugin struct {
	closed bool
}

func (this *TracePlugin) WrapCall(w *wrap.Wrap) error {
	if this.closed {
		return nil
	}

	span, ctx := traceutil.TraceRPCXExtract(w.GetCtx(), w.MethodTag)
	if span != nil {
		defer span.Finish()
	}
	w.SetCtx(ctx)
	e := w.Next()

	return e

}

func (this *TracePlugin) UpdateConfig(configs config.PluginConfig) {
	this.closed = configs.TRConfigs.Closed
	return
}

func InitTrace(configs config.PluginConfig) *TracePlugin {
	plugin := new(TracePlugin)
	plugin.closed = configs.TRConfigs.Closed
	return plugin
}
