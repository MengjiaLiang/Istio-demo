package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{
		contextID: contextID,
	}
}

func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogDebug("loading plugin config")

	return types.OnPluginStartStatusOK
}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	rawVal, _, err := proxywasm.GetSharedData("key")

	if err != nil {
		proxywasm.SetSharedData("key", []byte("value"), 0)
	} else {
		proxywasm.ReplaceHttpRequestHeader("X-UiPath-Internal-SharedData", string(rawVal))
	}

	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
