package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	_ "github.com/wasilibs/nottinygc"
	"wasm-filters/dispatch-http-call/http"
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
	proxywasm.LogInfof("Dispatching the HTTP call")
	headers := [][2]string{
		{":method", "GET"},
		{":path", "/api/requestrouting/location/mjorg/portal"},
		{":authority", "ingressgateway"},
		{":scheme", "http"},
	}

	if _, err := proxywasm.DispatchHttpCall(
		"outbound|80||platform-location-service.uipath.svc.cluster.local",
		headers,
		nil,
		nil,
		5000,
		ctx.dispatchCallback,
	); err != nil {
		panic(err)
	}

	return types.ActionPause
}

// dispatchCallback is the callback function called in response to the response arrival from the dispatched request.
func (ctx *httpHeaders) dispatchCallback(numHeaders, bodySize, numTrailers int) {
	proxywasm.LogInfof("executing the callback")

	response, _ := GetGetLocationResponse("mjorg", "", "portal", bodySize)

	responseBody := response.GetBody()
	proxywasm.LogInfof("Location response body in callback: %s", responseBody)
	proxywasm.ReplaceHttpRequestHeader("X-UiPath-AccountId", "testid")

	proxywasm.ResumeHttpRequest()
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

// GetGetLocationResponse returns a response object for GetLocation response, and does the error handling.
func GetGetLocationResponse(org, tenant, serviceType string, bodySize int) (*http.Response, string) {
	response, err := http.NewResponse(bodySize)
	if err != nil {
		proxywasm.LogErrorf("GetLocation() failed to get the response from LS API: %v", err)
		return nil, err.Error()
	}

	responseCode := response.GetStatus()
	responseBody := response.GetBody()

	proxywasm.LogInfof("Location response status: %d", responseCode)
	proxywasm.LogInfof("Location response body: %s", responseBody)

	return response, ""
}

//export sched_yield
func sched_yield() int32 {
	return 0
}
