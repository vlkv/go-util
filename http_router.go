package util

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"errors"
	"regexp"
	"fmt"
	"strings"
	"net/url"
)


type HttpRouter struct {
	router *httprouter.Router
	routes map[HttpRouteId]*HttpRoute
}

func NewHttpRouter() *HttpRouter {
	result := HttpRouter{}
	result.router = httprouter.New()
	result.routes = map[HttpRouteId]*HttpRoute{}
	return &result
}

type HttpRouteId string

var _ fmt.Stringer = (*HttpRouteId)(nil)

func (this HttpRouteId) String() string {
	return string(this)
}

// TODO: use typedef for paramValues map[string]interface{}
type HttpHandler func(routeId HttpRouteId, w http.ResponseWriter, r *http.Request, paramValues map[string]interface{})

type HttpRoute struct {
	Path string
	Method HttpMethod
	UrlParams []HttpParam
	QueryParams []HttpParam
	FormParams []HttpParam
	Handler HttpHandler
}

func (this *HttpRoute) getAllParams() []HttpParam {
	result := make([]HttpParam, 0, len(this.UrlParams) + len(this.QueryParams) + len(this.FormParams))
	result = append(result, this.UrlParams...)
	result = append(result, this.QueryParams...)
	result = append(result, this.FormParams...)
	return result
}

func (this *HttpRoute) getAllOptionalParams() []HttpParam {
	result := make([]HttpParam, 0)
	allParams := this.getAllParams()
	for i := range allParams {
		p := allParams[i]
		if (p.IsOptional()) {
			result = append(result, p)
		}
	}
	return result
}

func (this *HttpRoute) getAllRequiredParams() []HttpParam {
	result := make([]HttpParam, 0)
	allParams := this.getAllParams()
	for i := range allParams {
		p := allParams[i]
		if (p.IsRequired()) {
			result = append(result, p)
		}
	}
	return result
}

func (this *HttpRoute) parseParamValues(r *http.Request, ps httprouter.Params) map[string]interface{} {
	paramValues := map[string]interface{}{}
	for _, p := range this.UrlParams {
		if p.IsMultiple {
			panic(errors.New("You cannot not use IsMultiple=true for URL param"))
		}
		var val string
		if p.IsRequired() {
			val = ParamByNameReq(&ps, p.Name, this.Path)
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
				val = QueryValueReq(r, p.Name, this.Path)
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
				val = FormValueReq(r, p.Name, this.Path)
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

type HttpMethod int
const (
	HttpMethod_GET HttpMethod = iota
	HttpMethod_POST
)

func NewHttpRoute(path string, method HttpMethod, params []HttpParam, handler HttpHandler) *HttpRoute {
	re := regexp.MustCompile(":[\\w-]+")
	urlParams := re.FindAllString(path, -1)
	for i := range urlParams {
		urlParam := urlParams[i]
		name := urlParam[1:]
		i := FindIndex(len(params), func(i int) bool { return params[i].Name == name })
		if i < 0 {
			panic(errors.New(fmt.Sprintf("Url param %s exists in path, but not declared in params array", name)))
		}
	}

	filterParams := func (t HttpParamType, params []HttpParam) []HttpParam {
		result := make([]HttpParam, 0)
		for i, _ := range params {
			p := params[i]
			if p.Type == t {
				result = append(result, p)
			}
		}
		return result
	}

	result := new(HttpRoute)
	result.Path = path
	result.Method = method
	result.UrlParams = filterParams(HttpParamType_URL, params)
	result.QueryParams = filterParams(HttpParamType_Query, params)
	result.FormParams = filterParams(HttpParamType_Form, params)
	result.Handler = handler
	return result
}

func (this *HttpRouter) DeclareRouteGET(routeId HttpRouteId, path string, handler HttpHandler, params ...HttpParam) {
	route := NewHttpRoute(path, HttpMethod_GET, params, handler)
	this.routes[routeId] = route
}

func (this *HttpRouter) DeclareRoutePOST(routeId HttpRouteId, path string, handler HttpHandler, params ...HttpParam) {
	route := NewHttpRoute(path, HttpMethod_POST, params, handler)
	this.routes[routeId] = route
}

func (this *HttpRouter) BindRoute(routeId HttpRouteId, handler HttpHandler) {
	route, ok := this.routes[routeId]
	if !ok {
		panic(errors.New(fmt.Sprintf("Route %s not found, cannot bind", routeId)))
	}
	if route.Handler != nil {
		panic(errors.New(fmt.Sprintf("Route %s is already bound, cannot rebind", routeId)))
	}
	route.Handler = handler
}

func (this *HttpRouter) addAllDeclaredRoutes() {
	for k, _ := range this.routes {
		route := this.routes[k]
		var methodFunc func (string, httprouter.Handle)
		if route.Method == HttpMethod_GET {
			methodFunc = this.router.GET
		} else if route.Method == HttpMethod_POST {
			methodFunc = this.router.POST
		} else {
			panic(errors.New(fmt.Sprintf("Unexpected method: %v", route.Method)))
		}

		if route.Handler == nil {
			panic(errors.New(fmt.Sprintf("Route %v has unbinded handler, cannot use such route", route.Path)))
		}

		routeId := k // ATTENTION: We need a copy of the outer routeId to put in the closure
		this.addRoute(methodFunc, route.Path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			paramValues := route.parseParamValues(r, ps)
			route.Handler(routeId, w, r, paramValues)
		})
	}
}

func (this *HttpRouter) AddNotFoundRoute(handler http.HandlerFunc) {
	this.router.NotFound = http.HandlerFunc(handler)
}

func (this *HttpRouter) ListenAndServe(addr string) error {
	this.addAllDeclaredRoutes()
	return http.ListenAndServe(addr, this.router)
}

func (this *HttpRouter) addRoute(methodFunc func (string, httprouter.Handle), route string, handler httprouter.Handle) {
	routeNoTrailingSlash := strings.TrimRight(route, "/")
	methodFunc(routeNoTrailingSlash, handler)
	methodFunc(routeNoTrailingSlash + "/", handler)
}



type HttpRequestParams struct {
	URL string
	Method HttpMethod
	Data url.Values
	hasQueryValuesAdded bool
}

func CreateHttpRequestParams(routePath string, routeMethod HttpMethod) HttpRequestParams {
	return HttpRequestParams{ URL: routePath, Method: routeMethod, Data: url.Values{} }
}

func (this *HttpRequestParams) addParamValue(paramType HttpParamType, paramName string, paramValue string) {
	switch paramType {
	case HttpParamType_URL:
		this.URL = strings.Replace(this.URL, ":" + paramName, paramValue, -1)
	case HttpParamType_Query:
		if this.hasQueryValuesAdded {
			this.URL = this.URL + "&" + paramName + "=" + paramValue
		} else {
			this.URL = this.URL + "?" + paramName + "=" + paramValue
			this.hasQueryValuesAdded = true
		}
	case HttpParamType_Form:
		this.Data.Add(paramName, paramValue)
	}
}

func (this *HttpRouter) CreateHttpRequest(routeId HttpRouteId, paramValues map[string]interface{}) HttpRequestParams {
	route, ok := this.routes[routeId]
	if !ok {
		panic(errors.New(fmt.Sprintf("Route %v not found", routeId)))
	}

	result := CreateHttpRequestParams(route.Path, route.Method)

	// Process all required params, panic if some values are missing
	reqParams := route.getAllRequiredParams()
	for i := range reqParams {
		p := reqParams[i]
		if p.IsMultiple {
			panic(errors.New(fmt.Sprintf("Multiple parameter cannot be required, %v, route: %v", p.Name, routeId)))
		}
		value, ok := paramValues[p.Name]
		if !ok {
			panic(errors.New(fmt.Sprintf("Value for required param %v is missing, route: %v", p.Name, routeId)))
		}
		result.addParamValue(p.Type, p.Name, value.(string))
	}

	// Process all optional params, use defaults where needed
	optParams := route.getAllOptionalParams()
	for i := range optParams {
		p := optParams[i]

		if !p.IsMultiple {
			value, ok := paramValues[p.Name]
			if !ok {
				value = p.DefaultValue
			}
			result.addParamValue(p.Type, p.Name, value.(string))
		} else {
			values, ok := paramValues[p.Name]
			if !ok {
				defValues := make([]string, 0)
				if p.DefaultValue != "" {
					defValues = append(defValues, p.DefaultValue) // NOTE: We support only non-multiple DefaultValue
				}
				values = defValues
			}
			stringValues := values.([]string)
			for j := range stringValues {
				value := stringValues[j]
				result.addParamValue(p.Type, p.Name, value)
			}
		}
	}

	return result
}