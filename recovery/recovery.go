package recovery

import (
	"fmt"
	"runtime/debug"

	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
)

type RecoveryPlugin struct {
	closed bool
}

func (this *RecoveryPlugin) WrapCall(w *wrap.Wrap) (err error) {
	if this.closed {
		return nil
	}
	ctx := w.GetCtx()
	defer func() {
		if r := recover(); r != nil {
			logger.Ex(ctx, "[Recovery]", "call func :%s, err: %v, stacks: %s", w.MethodTag, r, string(debug.Stack()))
			err = fmt.Errorf("[service internal error]: %v", r)
		}
	}()
	err = w.Next()
	return
}

func (this *RecoveryPlugin) UpdateConfig(configs config.PluginConfig) {
	this.closed = configs.RCConfigs.Closed
	return
}

func InitRecovery(configs config.PluginConfig) *RecoveryPlugin {
	plugin := new(RecoveryPlugin)
	plugin.closed = configs.RCConfigs.Closed
	return plugin
}
