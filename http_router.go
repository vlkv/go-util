package util

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"errors"
	"regexp"
	"fmt"
	"strings"
)


type HttpRouter struct {
	router *httprouter.Router
	routes map[string]HttpRoute
}

func NewHttpRouter() *HttpRouter {
	result := HttpRouter{}
	result.router = httprouter.New()
	result.routes = map[string]HttpRoute{}
	return &result
}

type HttpHandler func(w http.ResponseWriter, r *http.Request, paramValues map[string]interface{})

type HttpRoute struct {
	Path string
	UrlParams []HttpParam
	QueryParams []HttpParam
	FormParams []HttpParam
	Handler HttpHandler
}

func (this *HttpRoute) getParamValues(r *http.Request, ps httprouter.Params) map[string]interface{} {
	paramValues := map[string]interface{}{}
	for _, p := range this.UrlParams {
		if p.IsMultiple {
			panic(errors.New("You should not use both IsMultiple=true for URL param"))
		}
		var val string
		if p.IsRequired() {
			val = ParamByNameReq(&ps, p.Name, "")
		} else {
			val = ParamByNameOpt(&ps, p.Name, p.DefaultValue)
		}
		paramValues[p.Name] = val
	}
	for _, p := range this.QueryParams {
		if p.IsMultiple {
			if p.IsRequired() {
				panic(errors.New("You should use IsMultiple=true only with ForceOptional=true"))
			}
			vals := r.URL.Query()[p.Name]
			paramValues[p.Name] = vals
		} else {
			var val string
			if p.IsRequired() {
				val = QueryValueReq(r, p.Name, "")
			} else {
				val = QueryValueOpt(r, p.Name, p.DefaultValue)
			}
			paramValues[p.Name] = val
		}
	}
	for _, p := range this.FormParams {
		if p.IsMultiple {
			if p.IsRequired() {
				panic(errors.New("You should not use both IsMultiple=true and IsRequired=true"))
			}
			vals := r.URL.Query()[p.Name]
			paramValues[p.Name] = vals
		} else {
			var val string
			if p.IsRequired() {
				val = FormValueReq(r, p.Name, "")
			} else {
				val = FormValueOpt(r, p.Name, p.DefaultValue)
			}
			paramValues[p.Name] = val
		}
	}
	return paramValues
}


type HttpParam struct {
	Type HttpParamType
	Name string
	DefaultValue string // Has sense only when IsMultiple==false
	ForceOptional bool
	IsMultiple bool
}

func (this *HttpParam) IsRequired() bool {
	return this.DefaultValue == "" && !this.ForceOptional
}

func (this *HttpParam) IsOptional() bool {
	return !this.IsRequired()
}

type HttpParamType int
const (
	HttpParamType_URL HttpParamType = iota
	HttpParamType_Query
	HttpParamType_Form
)

func CreateHttpRoute(path string, params []HttpParam, handler HttpHandler) HttpRoute {
	re := regexp.MustCompile(":[\\w-]+")
	urlParams := re.FindAllString(path, -1)
	for _, urlParam := range urlParams {
		name := urlParam[1:]
		i := FindIndex(len(params), func(i int) bool { return params[i].Name == name })
		if i < 0 {
			panic(errors.New(fmt.Sprintf("Url param %s exists in path, but not declared in params array", name)))
		}
	}

	filterParams := func (t HttpParamType, params []HttpParam) []HttpParam {
		result := make([]HttpParam, 0)
		for _, p := range params {
			if p.Type == t {
				result = append(result, p)
			}
		}
		return result
	}

	return HttpRoute{
		Path: path,
		UrlParams: filterParams(HttpParamType_URL, params),
		QueryParams: filterParams(HttpParamType_Query, params),
		FormParams: filterParams(HttpParamType_Form, params),
		Handler: handler,
	}
}

func (this *HttpRouter) AddRouteGET(routeId, path string, handler HttpHandler, params ...HttpParam) {
	route := CreateHttpRoute(path, params, handler)
	this.routes[routeId] = route
	this.addRoute(this.router.GET, path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		paramValues := route.getParamValues(r, ps)
		handler(w, r, paramValues)
	})
}

func (this *HttpRouter) AddRoutePOST(routeId, path string, handler HttpHandler, params ...HttpParam, ) {
	route := CreateHttpRoute(path, params, handler)
	this.routes[routeId] = route
	this.addRoute(this.router.POST, path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		paramValues := route.getParamValues(r, ps)
		handler(w, r, paramValues)
	})
}

func (this *HttpRouter) AddNotFoundRoute(handler http.HandlerFunc) {
	this.router.NotFound = http.HandlerFunc(handler)
}

func (this *HttpRouter) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, this.router)
}

func (this *HttpRouter) addRoute(routerMethod func (string, httprouter.Handle), route string, handler httprouter.Handle) {
	routeNoTrailingSlash := strings.TrimRight(route, "/")
	routerMethod(routeNoTrailingSlash, handler)
	routerMethod(routeNoTrailingSlash + "/", handler)
}
