package rpcxplugin

import (
	"context"
	"errors"
	"sync"

	"github.com/tal-tech/odinPlugin/config"
	"github.com/tal-tech/odinPlugin/wrap"
)

type MiddleWare interface {
	WrapCall(*wrap.Wrap) error
	UpdateConfig(config.PluginConfig)
}

type RpcxPlugin struct {
	EndPoints   []wrap.EndPoint
	MiddleWares []MiddleWare
	ServiceName string
	Reregister  func() error
}

func NewRpcxPlugin(sn string) *RpcxPlugin {
	ins := new(RpcxPlugin)
	ins.ServiceName = sn
	conf := config.Conf[sn]
	conf.MSConfigs.Prefix = sn
	config.Conf[sn] = conf
	tempMiddleWares := make([]MiddleWare, 0)
	tempEndPoints := make([]wrap.EndPoint, 0)
	for _, op := range MiddlewareOptions {
		middleware, endpoint := op(conf)
		tempMiddleWares = append(tempMiddleWares, middleware)
		tempEndPoints = append(tempEndPoints, endpoint)
	}
	ins.MiddleWares = tempMiddleWares
	ins.EndPoints = tempEndPoints
	addPlugin(sn, ins)
	return ins
}

func (this *RpcxPlugin) Wrapcall(ctx context.Context, tag string, end wrap.EndPoint) (err error) {
	wrap := wrapPool.Get().(*wrap.Wrap)
	wrap.Methods = this.EndPoints
	wrap.Fn = end
	wrap.MethodTag = tag
	wrap.SetCtx(ctx)
	err = wrap.Next()
	wrap.Reset()
	wrapPool.Put(wrap)
	return
}

func (this *RpcxPlugin) Update() (err error) {
	conf, ok := config.Conf[this.ServiceName]
	if !ok {
		return errors.New("get conf failed")
	}
	for _, mid := range this.MiddleWares {
		mid.UpdateConfig(conf)
	}

	return nil
}

func (this *RpcxPlugin) SetReregister(fn func() error) {
	this.Reregister = fn
}

func (this *RpcxPlugin) DoReregister() (err error) {
	if this.Reregister != nil {
		err = this.Reregister()
	} else {
		return errors.New("reregister func is nil")
	}
	return err
}

func (this *RpcxPlugin) GetMetadata() (ret string) {
	conf, ok := config.Conf[this.ServiceName]
	if !ok {
		return
	}
	md := conf.MDConfigs
	for k, v := range md {
		ret += (k + "=" + v + "&")
	}
	return
}

var wrapPool = sync.Pool{
	New: func() interface{} {
		return new(wrap.Wrap)
	},
}
