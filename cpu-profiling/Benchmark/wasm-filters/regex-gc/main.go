package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	_ "github.com/wasilibs/nottinygc"
	"regexp"
	"strings"
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
	path, _ := proxywasm.GetHttpRequestHeader(":path")

	_, _, svc, _ := ExtractStandardParams(path)

	err := proxywasm.ReplaceHttpRequestHeader("serviceType", svc)
	if err != nil {
		proxywasm.LogCritical("failed to set request header: test")
	}

	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

func ExtractStandardParams(path string) (string, string, string, string) {
	var (
		standardUrlPath string
		pathSegment     string
		pathSegmentIdx  int
	)

	pathSegment = ""
	pathSegmentIdx = strings.Index(path, "_/")
	if pathSegmentIdx == -1 {
		// some service put API version after the `{serviceType}_?`, e.g. "/org/tenant/apps_?api1/get"
		pathSegmentIdx = strings.Index(path, "_?")
	}

	if pathSegmentIdx != -1 {
		// Slicing string works as substring in other language. Ref: https://go.dev/ref/spec#Slice_expressions
		pathSegment = path[pathSegmentIdx+1:]
		standardUrlPath = path[0 : pathSegmentIdx+1]
	} else {
		standardUrlPath = path
	}

	tenantRegex := `^\/([\w\-]+[^\W_])\/([\w\-]+[^\W_])\/([a-zA-Z]+)_$`
	re := regexp.MustCompile(tenantRegex)
	results := re.FindStringSubmatch(standardUrlPath)
	if len(results) == 0 {
		// try to match organization standard url pattern
		orgRegex := `^\/([\w\-]+[^\W_])\/([a-zA-Z]+)_$`
		re = regexp.MustCompile(orgRegex)
		results = re.FindStringSubmatch(standardUrlPath)

		if len(results) != 0 {
			return results[1], "", results[2], pathSegment
		}

		return "", "", "", pathSegment
	}

	return results[1], results[2], results[3], pathSegment
}

//export sched_yield
func sched_yield() int32 {
	return 0
}
