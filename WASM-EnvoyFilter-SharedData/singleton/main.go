package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const tickMilliseconds uint32 = 15000

var authHeader string

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
	return &pluginContext{contextID: contextID}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	contextID uint32
	callBack  func(numHeaders, bodySize, numTrailers int)
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{contextID: contextID}
}

type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("set tick period milliseconds: %d", tickMilliseconds)
	// ctx.callBack = func(numHeaders, bodySize, numTrailers int) {
	// 	respHeaders, _ := proxywasm.GetHttpCallResponseHeaders()
	// 	proxywasm.LogInfof("respHeaders: %v", respHeaders)

	// 	for _, headerPairs := range respHeaders {
	// 		if headerPairs[0] == "authorization" {
	// 			authHeader = headerPairs[1]
	// 		}
	// 	}
	// }
	return types.OnPluginStartStatusOK
}

func (ctx *httpContext) OnHttpRequestHeaders(int, bool) types.Action {
	proxywasm.LogInfo("Request received.")

	// headers, err := proxywasm.GetHttpRequestHeaders()
	// if err != nil {
	// 	proxywasm.LogCriticalf("failed to get request headers: %v", err)
	// 	return types.ActionContinue
	// }

	// TODO: implement current routing logic here

	// TODO: extract Org name and Tenant name
	return types.ActionContinue
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnTick() {
	proxywasm.LogDebug("Tick")
}
