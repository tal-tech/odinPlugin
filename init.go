package rpcxplugin

import (
	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/odinPlugin/config"
)

var pluginMap map[string]*RpcxPlugin

func Init(port string, ops ...Options) {
	for _, op := range ops {
		MiddlewareOptions = append(MiddlewareOptions, op)
	}
	config.InitFromFile()
	config.InitConfManager(port)
	return
}

func init() {
	pluginMap = make(map[string]*RpcxPlugin, 0)
	go watchConf()
}

func watchConf() {
	uchan := config.GetUpdateChan()
	mchan := config.GetMetadataChan()
	for {
		select {
		case name := <-uchan:
			update(name)
		case name := <-mchan:
			reregister(name)
		}
	}
}

func addPlugin(name string, p *RpcxPlugin) {
	pluginMap[name] = p
}

func update(name string) {
	for s, p := range pluginMap {
		if s == name {
			p.Update()
		}
	}
	return
}

func reregister(name string) {
	for s, p := range pluginMap {
		if s == name {
			err := p.DoReregister()
			if err != nil {
				logger.E("odinPlugin", "reregister %s,error:%v", name, err)
			}
		}
	}
	return
}
