package aurora

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	auroraQueryCache         = "auroraQueryCache"
	auroraFormCache          = "auroraFormCache"
	auroraMaxMultipartMemory = "auroraMaxMultipartMemory"
)

// Middleware 中间件类型
type Middleware func(Ctx) bool

// Ctx 上下文参数，主要用于在业务之间传递 数据使用
// 上下文参数中获取请求参数需要依赖于传递的参数名称
// Ctx 不是线程安全的，在请求中出现多线程操作需要使用锁来保证安全性
type Ctx map[string]interface{}

func (c Ctx) Clear() {
	for key, _ := range c {
		delete(c, key)
	}
}

// Request 返回元素 Request
func (c Ctx) Request() *http.Request {
	return c[request].(*http.Request)
}

// Response 返回元素 ResponseWriter
func (c Ctx) Response() http.ResponseWriter {
	return c[response].(http.ResponseWriter)
}

// Ref 获取容器中的依赖项
func (c Ctx) Ref(ref string) interface{} {
	if v := c[iocs].(*ioc).get(ref); v == nil {
		return nil
	} else {
		return v.Interface()
	}
}

// Return 设置中断处理，多次调用会覆盖之前设置的值
func (c Ctx) Return(value ...interface{}) {
	values := make([]reflect.Value, 0)
	for _, v := range value {
		values = append(values, reflect.ValueOf(v))
	}
	c["AuroraValues"] = values
}

func (c Ctx) initQueryCache() {
	req := c.Request()
	switch req.Method {
	case http.MethodGet:
		if _, b := c[auroraQueryCache]; !b {
			c[auroraQueryCache] = req.URL.Query()
		}
	case http.MethodPost, http.MethodPut:
		if _, b := c[auroraFormCache]; !b {
			size := c[auroraMaxMultipartMemory].(int64)
			err := req.ParseMultipartForm(size)
			if err != nil {
				return
			}
			c[auroraFormCache] = req.PostForm
		}
	}
}
func (c Ctx) cacheQuery() url.Values {
	c.initQueryCache()
	return c[auroraQueryCache].(url.Values)
}

func (c Ctx) cacheForm() url.Values {
	c.initQueryCache()
	return c[auroraFormCache].(url.Values)
}

func (c Ctx) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}
func (c Ctx) QueryMap(key string) map[string]string {
	dicts, _ := c.GetQueryMap(key)
	return dicts
}

func (c Ctx) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c Ctx) GetQueryArray(key string) ([]string, bool) {
	if values, ok := c.cacheQuery()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c Ctx) GetQueryMap(key string) (map[string]string, bool) {
	return c.get(c.cacheQuery(), key)
}

func (c Ctx) PostFormMap(key string) map[string]string {
	dicts, _ := c.GetPostFormMap(key)
	return dicts
}

func (c Ctx) GetPostFormArray(key string) ([]string, bool) {
	if values := c.cacheForm()[key]; len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c Ctx) GetPostFormMap(key string) (map[string]string, bool) {
	return c.get(c.cacheForm(), key)
}

// FormFile 获取文件
func (c Ctx) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request().MultipartForm == nil {
		if err := c.Request().ParseMultipartForm(c[auroraMaxMultipartMemory].(int64)); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Request().FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (c Ctx) MultipartForm() (*multipart.Form, error) {
	err := c.Request().ParseMultipartForm(c[auroraMaxMultipartMemory].(int64))
	return c.Request().MultipartForm, err
}

// AddHeader 添加一个头
func (c Ctx) AddHeader(name, value string) {
	rew := c.Response()
	if rew.Header().Get(name) == "" {
		rew.Header().Set(name, value)
		return
	}
	rew.Header().Add(name, value)
}

// GetHeader 根据 name 查找一个
func (c Ctx) GetHeader(name string) string {
	return c.Response().Header().Get(name)
}

// DelHeader 删除一个指定的name 的头
func (c Ctx) DelHeader(name string) {
	rew := c.Response()
	if h := rew.Header().Get(name); h != "" {
		rew.Header().Del(name)
	}
}

// SetCookie 设置一个 Cookie
func (c Ctx) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteDefaultMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c Ctx) get(m map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = v[0]
			}
		}
	}
	return dicts, exist
}
