package cors

import (
	"gitee.com/aurora-engine/aurora"
	"net/http"
	"strings"
)

const (
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

func New() *Cors {
	c := &Cors{}
	c.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"}
	c.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	c.AllowCredentials = []string{"false"}
	c.ExposeHeaders = []string{"Content-Length", " Access-Control-Allow-Origin", " Access-Control-Allow-Headers", " Cache-Control", " Content-Language", " Content-Type"}
	c.Method = map[string][]string{
		http.MethodOptions: {"/*"},
	}
	return c
}

type Cors struct {
	Method           map[string][]string //运行跨域的 api
	Host             map[string][]string //允许跨域的主机
	AllowOrigin      []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials []string
	ExposeHeaders    []string
}

// Domain 添加跨域路径
func (c *Cors) Domain(method string, url ...string) {
	if c.Method == nil {
		c.Method = make(map[string][]string)
	}
	if _, b := c.Method[method]; !b {
		c.Method[method] = make([]string, 0)
	}
	c.Method[method] = append(c.Method[method], url...)
}

// Origin 配置可跨域主机
func (c *Cors) Origin(host string, method ...string) {
	if c.Host == nil {
		c.Host = make(map[string][]string)
	}
	c.Host[host] = method
}

func (c *Cors) Cors() aurora.Middleware {
	return func(ctx aurora.Ctx) bool {
		return c.requestCheck(ctx.Request(), ctx.Response())
	}
}

// requestCheck 检查请求是否在可跨域列表中
func (c *Cors) requestCheck(r *http.Request, w http.ResponseWriter) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		if !c.checkRequest(r, w) {
			return false
		}
		if _, b := c.Method[r.Method]; b {
			for _, p := range c.Method[r.Method] {
				if strings.HasSuffix(p, "*") || p == r.RequestURI {
					if strings.HasPrefix(r.RequestURI, p[:len(p)-1]) || p == r.RequestURI {
						// 路径匹配成功 写入跨域头
						w.Header().Set(AccessControlAllowOrigin, origin)
						w.Header().Set(AccessControlAllowMethods, r.Method)
						w.Header().Set(AccessControlAllowHeaders, strings.Join(c.AllowHeaders, ","))
						w.Header().Set(AccessControlExposeHeaders, strings.Join(c.ExposeHeaders, ","))
						w.Header().Set(AccessControlAllowCredentials, strings.Join(c.AllowCredentials, ","))
						return true
					}
				}
			}
			w.WriteHeader(http.StatusNoContent)
			return false
		} else {
			w.WriteHeader(http.StatusNoContent)
			return false
		}
	}
	return true
}

// checkRequest 跨域预检测处理
func (c *Cors) checkRequest(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get("Origin") != "" {
		if r.Method == http.MethodOptions {
			if c.Host == nil {
				w.Header().Set(AccessControlAllowOrigin, "*")
				w.Header().Set(AccessControlAllowMethods, strings.Join(c.AllowMethods, ","))
			} else {
				if m, b := c.Host[r.Host]; b {
					w.Header().Set(AccessControlAllowOrigin, r.Host)
					w.Header().Set(AccessControlAllowMethods, strings.Join(m, ","))
				} else {
					w.WriteHeader(http.StatusNoContent)
					return false
				}
			}
			if c.ExposeHeaders != nil {
				w.Header().Set(AccessControlAllowHeaders, strings.Join(c.AllowHeaders, ","))
			} else {
				w.Header().Set(AccessControlAllowHeaders, "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			}
			if c.ExposeHeaders != nil {
				w.Header().Set(AccessControlExposeHeaders, strings.Join(c.ExposeHeaders, ","))
			} else {
				w.Header().Set(AccessControlExposeHeaders, "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			}
			if c.AllowCredentials != nil {
				w.Header().Set(AccessControlAllowCredentials, strings.Join(c.AllowCredentials, ","))
			} else {
				w.Header().Set(AccessControlAllowCredentials, "false")
			}
			w.WriteHeader(http.StatusNoContent)
			return false
		}
	}

	return true
}
