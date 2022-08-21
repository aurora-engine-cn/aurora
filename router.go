package aurora

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

/*
	基于字典树的路由器
	路由器规则:
		1.	存储相同路径的处理器，后者覆盖前者
		2.	注册路径必须以 / 开头 (自动处理)
		3.	路径结尾不能以 / 结尾 (自动处理) 此处的要求只是一个框架规范
		4.	路径分隔符只能 单个斜杠 / (多个斜杠会自动处理)
		5.	注册两个形同路径 会被阻止 同样的RESTFul也是如此  /aa/{b} /aa/{b} 被认为是形同的路径（待添加校验，对于RESTFul来说这是个歧义）
	RESTFul:
		1.  使用 {xxx} 标识符来 代替一个任意路径
		2.  /xx{xx}xx/xx  /{xx}xx  /xx{xx}/xx  类似格式均为非法
		3.  有共同父路径(父路径没有RESTFul) 在第一个RESTFul之后就无法存储更多的子路径
				/abc/ccc/{c}
				/abc/ccc/a
				两个路径注册顺序无关都会发生冲突
		4.	RESTFul 解析规则，因为无法保证在同一条路径上 RESTFul的命名规范，比如
			/aaa/bb/{c}/bb/{c}
			这样的路由只能通过 参数自动解析的方式来完成同名的 {c} 部位参数获取

	路由存储规则:
		1.	两个路径有相同前缀，则提取两个路径的公共前缀最为根
		2.	RESTFul 子路径必须完整保存在单个路径节点上，不允许RESTFul产生分裂提取公共前缀，比如
			/abc/{abbc}
			/abc/{abaa}
			使用者可能想表达两个不同的路径，但是两个路径想要表达两个不同的路径这是不允许的，路径会存在歧义，路径为 /abc/1 的时候不知道具体执行哪一个。
			/abc/{ab 部分会产生提 导致TESTFul 不在一个完整的子路径上。
		3.	在同一条 RESTFul 父路径下可以有多条普通子路径
			/aa/abc/{aa}
			/aa/abc/{aa}/bb
			/aa/abc/{aa}/cc
			/aa/abc/{aa}/cc/a
			...
		4. 带有RESTFul作为父路径，RESTFul的命名必须相同 如下
			/aa/abc/{aa}
			/aa/abc/{aa}/bb
			再次注册 /aa/abc/{aaa}   /aa/abc/{aaa}/...  视为非法 {aa} 和 {aaa} 不视为同一条路径
		5. 相同 RESTFul 父路径 下面可以有子路径，规则相同

发生公共根

		1)节点和被添加路径产生公共根，提取公共根后，若公共根未注册，服务处理函数将为nil
		2)若节点恰好是公共根，则设置函数
	6.REST 风格注册
		1)同一个根路径下只能有一个 REST 子路径
		2)REST 作为根路径也只能拥有一个REST 子路径
		3)REST 路径会和其它非REST同级路径发生冲突
	7.注册路径不能以/结尾（bug未修复，/user /user/ 产生 /user 的公共根 使用切割解析路径方式，解析子路径，拼接剩余子路径会存在bug ,注册路径的时候强制无法注册 / 结尾的 url）
*/
const (
	request  = "AuroraRequest"  //go 原生请求
	response = "AuroraResponse" //go 原生响应
	iocs     = "AuroraIoc"      //依赖容器
)

// route Aurora核心路由器
type route struct {
	*Engine     // Aurora 引用
	mx          *sync.Mutex
	catch       map[reflect.Type]catch // 全局错误捕捉处理
	middleware  []Middleware           // 全局中间件
	tree        map[string]*node       // 路由树根节点
	defaultView views                  // 默认视图处理器，初始化采用 Aurora 实现的函数进行渲染
}

// node 路由节点
type node struct {
	Path       string       //当前节点路径
	FullPath   string       //当前处理器全路径
	Count      int          //路径数量
	middleware []Middleware //中间处理函数
	Control    *controller  //服务处理函数
	Child      []*node      //子节点
}

func (r *route) use(middleware ...Middleware) {
	if middleware == nil {
		return
	}
	if r.middleware == nil {
		r.middleware = make([]Middleware, len(middleware))
		for i := range middleware {
			r.middleware[i] = middleware[i]
		}
		return
	}
	for i := range middleware {
		r.middleware = append(r.middleware, middleware[i])
	}
}

func (r *route) Catch(err Error) {
	r.registerErrorCatch(err)
}

// ——————————————————————————————————————————————————————————————————————————路由注册————————————————————————————————————————————————————————————————————————————————————————————
// addRoute 预处理被添加路径
func (r *route) addRoute(method, path string, control Controller, middleware ...Middleware) {
	//非空校验,
	if path == "" || control == nil {
		// 空字符串路径不能注册
		return
	}
	path = urlCheck(path)
	path = urlHead(path)
	path = urlEnd(path)
	err := checkRESTFul(path)
	ErrorMsg(err)
	//校验处理函数的正确性，只能注册函数，不能注册结构体，接口，基本类型等数据
	vt := reflect.TypeOf(control)
	if vt.Kind() != reflect.Func {
		r.Error(method + ":the registered handler is not a function，need a function")
		return
	}

	r.mx.Lock()
	defer r.mx.Unlock()
	//初始化路由树
	if r.tree == nil {
		r.tree = make(map[string]*node)
	}
	if _, ok := r.tree[method]; !ok {
		//初始化 请求类型根
		//初始化根路径,此处是更改 路径注册中的一些bug 而添加，由于 /路径注册的顺序导致了一些意想不到的bug, 特殊情况下 /aa  /a / 等顺序会导致其它两个出现错误
		r.tree[method] = &node{Path: "/"}
	}
	//拿到根路径
	root := r.tree[method]
	r.add(method, root, path, path, control, middleware...) //把路径添加到根路径中中
}

// add 添加路径节点
// method 指定请求类型，root 根路径，Path和fun 被添加的路径和处理函数，path携带路径副本添加过程中不会有任何操作仅用于日志处理
// method: 请求类型(日志相关参数)
// path: 插入的路径(日志相关参数)
func (r *route) add(method string, root *node, Path string, path string, fun Controller, middleware ...Middleware) {
	var l string
	vf := reflect.ValueOf(fun)
	vt := reflect.TypeOf(fun)
	control := &controller{Fun: vf, FunType: vt, Engine: r.Engine}
	control.InitArgs()
	//初始化根,此处的初始化根在Aurora 实例化阶段代替，该段if后期可以暂时忽略，没有初始化的根路由 的第一个节点默认为 ""以此判断初始化
	if root.Path == "" && root.Child == nil {
		root.Path = Path
		root.FullPath = path
		root.Count = strings.Count(path, "/")
		root.Child = nil
		root.Control = control
		root.middleware = middleware

		l = fmt.Sprintf("%-6s  %-10s   %-10s", method, path, getFunName(vf.Interface()))
		r.Debug(l)
		return
	}
	if root.Path == Path { //相同路径可能是分裂或者提取的公共根
		//此处修改，注册同样的路径，选择覆盖前一个
		if root.Control != nil {
			r.Error(method, ": ", path, " already registered")
			os.Exit(-1)
		}
		root.Control = control
		root.middleware = middleware
		root.FullPath = path
		root.Count = strings.Count(path, "/")
		l = fmt.Sprintf("%-6s  %-10s   %-10s", method, path, getFunName(vf.Interface()))
		r.Debug(l)
		return
	}
	//如果当前的节点是 REST API 节点 ，子节点可以添加REST API节点
	//如果当前节点的子节点存在REST API 则不允许添加子节点

	//检擦添加路径 和 当前路径 的关系   Path:添加的路径串 path:当前root的路径（此处path只是和被添加Path区分开，并不是参数中的path）
	//1.Path 长度小于 当前path长度---> (Path 和path 有公共前缀，Path是path的父路径)
	//2.Path 长度大于 当前path长度---> (path 和Path 有公共前缀，path是path的父路径)
	//3.以上两种情况都不满足
	rootPathLength := len(root.Path)    //当前路径长度
	addPathLength := len(Path)          //被添加路径长度
	if rootPathLength < addPathLength { //情况2. 当前节点可能是父节点
		if strings.HasPrefix(Path, root.Path) { //前缀检查
			i := len(root.Path)    //当前root节点路径作为父节点，计算长度用于截取被添加路径的子路径部分
			c := Path[i:]          //截取需要存储的 子路径，被截取的子路径是 待添加路径中截取出来的
			if root.Child != nil { //若当前root有子节点，查看当前被截取需要存储的子节点是否有二级父节点
				for i := 0; i < len(root.Child); i++ {
					/*
						判断前缀 确定当前路径的子节点 是不是和待插入节点 c 具有公共前缀，有公共前缀说明 待插入的 c 是这个子节点的父级
						检查该节点的子节点和和要存储的子路径是否存存在父子关系
						存在父子关系则进入递归
					*/
					if strings.HasPrefix(root.Child[i].Path, c) || strings.HasPrefix(c, root.Child[i].Path) {
						//此处的递归 是将子节点插入当前节点的子路径做检查，所以传递的路径和处理函数是当前正准备添加的函数
						r.add(method, root.Child[i], c, path, fun, middleware...)
						return // / 根路径在后面插入路由 无法走到最下面的 合并api 此处 注释return 解决
					}
				}
			} else {
				//添加子节点
				if root.Child == nil {
					root.Child = make([]*node, 0)
				}
				if len(root.Child) > 0 {
					//如果存储的路径是REST API 检索 当前子节点是否存有路径，存有路径则为冲突
					for i := 0; i < len(root.Child); i++ {
						// strings.HasPrefix(root.Child[i].Path, "{") 判断当前子路径节点是否是 RESTFul
						// strings.HasSuffix(Path, "{") 判断 待加入的子路径是不是 RESTFul
						if !(strings.HasPrefix(root.Child[i].Path, "{") && strings.HasSuffix(Path, "{")) {
							r.Error(method + ":" + path + " RESTFul conflict")
							os.Exit(-1)
						}
					}
				}
				n := &node{
					Path:       c,
					FullPath:   path,
					Count:      strings.Count(path, "/"),
					middleware: middleware,
					Control:    control,
					Child:      nil,
				}
				root.Child = append(root.Child, n)
				l = fmt.Sprintf("%-6s  %-10s   %-10s", method, path, getFunName(vf.Interface()))
				r.Debug(l)
				return
			}
		}
	}
	if rootPathLength > addPathLength { //情况1.当前节点可能作为子节点,被添加的节点作为父节点
		if strings.HasPrefix(root.Path, Path) { //前缀检查
			i := len(Path)     //
			c := root.Path[i:] //需要存储的子路径，c是被分裂出来的子路径(当前节点作为父节点)
			if root.Child != nil {
				for i := 0; i < len(root.Child); i++ {
					/*
						检查该节点的子节点和和要存储的子路径是否存存在父子关系
						存在父子关系则进入递归
					*/
					if strings.HasPrefix(root.Child[i].Path, c) || strings.HasPrefix(c, root.Child[i].Path) {
						//r.add(method, root.Child[i], c, path, fun)
						r.add(method, root.Child[i], c, path, root.Control.Fun.Interface(), root.middleware...) //改  此处的for主要处理  当前路径需要把分裂出来的子路径存储到当前的孩子节点中，传递的被存储的处理器应该是当前的处理器,如出现bug，恢复上面的注释代码
						return
					}
				}
			} else {
				//添加子节点
				if root.Child == nil {
					root.Child = make([]*node, 0)
				}
				if len(root.Child) > 0 {
					//如果存储的路径是REST API 需要检索当前子节点是否存有路径，存有路径则为冲突
					for i := 0; i < len(root.Child); i++ {
						if !(strings.HasPrefix(root.Child[i].Path, "{") && strings.HasSuffix(Path, "{")) {
							r.Error(method + ":" + path + " RESTFul conflict")
							os.Exit(-1)
						}
					}
				}
				tempChild := root.Child       //保存要一起分裂的子节点
				root.Child = make([]*node, 0) //清空当前子节点  root.Child=root.Child[:0]无法清空存在bug ，直接分配保险
				root.Child = append(root.Child,
					&node{
						Path:       c,
						FullPath:   root.FullPath,
						Count:      root.Count,
						Child:      tempChild,
						middleware: root.middleware,
						Control:    root.Control,
					},
				) //封装被分裂的子节点 添加到当前根的子节点中
				root.Path = root.Path[:i] //修改当前节点为添加的路径，（被添加结点刚好是父节点）
				root.FullPath = path
				root.Count = strings.Count(path, "/")
				root.Control = control       //更改当前处理函数
				root.middleware = middleware //更改当前中间件
				//此处的操作，大多是处理以竟被注册好的接口进行分裂，注释此处的目的是对日志的控制，被分裂的路径会被二次打印
				l = fmt.Sprintf("%-6s  %-10s   %-10s", method, path, getFunName(vf.Interface()))
				r.Debug(l)
				return
			}
		}
	}
	//情况3.节点和被添加节点无直接关系 抽取公共前缀最为公共根
	r.merge(method, root, Path, path, fun, middleware...)
	return
}

// merge 检测root节点 和待添加路径 是否有公共根，有则提取合并公共根
// method: 请求类型(日志相关参数)
// path: 插入的路径(日志相关参数)
// root: 根合并相关参数
// Path: 根合并相关参数
// fun: 根合并相关参数
func (r *route) merge(method string, root *node, Path string, path string, fun interface{}, middleware ...Middleware) bool {
	//处理反射
	vf := reflect.ValueOf(fun)
	vt := reflect.TypeOf(fun)
	control := &controller{Fun: vf, FunType: vt, Engine: r.Engine}
	control.InitArgs()
	pub := r.findPublicRoot(method, root.Path, Path, path) //公共路径
	if pub != "" {
		pl := len(pub)
		/*
			此处是提取当前节点公共根以外的操作，若当前公共根是root.Path自身则会取到空字符串 [:] 切片截取的特殊原因
			root.Path[pl:] pl是自生长度，取到最后一个字符串需要pl-1，pl取到的是个空，字符串默认为"",
			root.Path[pl:]取值为""时，说明root.Path本身就是就是公共根，只需要添加另外一个子节点到它的child切片即可
		*/
		ch1 := root.Path[pl:]
		ch2 := Path[pl:]
		if root.Child == nil {
			root.Child = make([]*node, 0)
		}
		if ch1 != "" {
			//ch1 本节点发生分裂 把处理函数也分裂 然后把当前的handler 置空,分裂的子节点也应该按照原有的顺序保留，分裂下去
			chChild := root.Child
			root.Child = make([]*node, 0) //重新分配
			root.Child = append(root.Child,
				&node{
					Path:     ch1,
					FullPath: root.FullPath,

					Count:      root.Count,
					Child:      chChild,
					middleware: root.middleware,
					Control:    root.Control,
				},
			) //把分裂的子节点全部信息添加到公共根中
			root.Control = nil //提取出来的公共根 没有可处理函数
			root.Count = 0
			root.FullPath = ""

		}
		if ch2 != "" {
			flag := true
			if len(root.Child) > 0 {
				for i := 0; i < len(root.Child); i++ {
					//单纯的被添加到此节点的子节点列表中 需要递归检测子节点和被添加节点是否有公共根
					if flag = r.merge(method, root.Child[i], ch2, path, fun, middleware...); flag {
						return true
					}
				}
				// 当前子路径如果没有能与之合并的节点 将选择添加到这个节点的子路径下
				// 检索插入路径REST API冲突。
				for i := 0; i < len(root.Child); i++ {
					if strings.HasPrefix(root.Child[i].Path, "{") || strings.HasPrefix(ch2, "{") {
						r.Error(method + ":" + path + " RESTFul conflict")
						os.Exit(-1)
					}
					if strings.HasPrefix(root.Child[i].Path, "{") && strings.HasPrefix(ch2, "{") {
						r.Error(method + ":" + path + " RESTFul conflict")
						os.Exit(-1)
					}
				}
			}
			n := &node{
				Path:       ch2,
				FullPath:   path,
				Count:      strings.Count(path, "/"),
				Child:      nil,
				middleware: middleware,
				Control:    control,
			}
			root.Child = append(root.Child, n) //作为新的子节点添加到当前的子节点列表中
			l := fmt.Sprintf("%-6s  %-10s   %-10s", method, path, getFunName(vf.Interface()))
			r.Debug(l)
		} else {
			//ch2为空说明 ch2是被添加路径截取的 添加的路径可能是被提出来作为公共根
			if pub == Path {
				root.Control = control
				root.FullPath = path
				root.middleware = middleware
				root.Count = strings.Count(path, "/")
				l := fmt.Sprintf("%-6s  %-10s   %-10s", method, path, getFunName(vf.Interface()))
				r.Debug(l)
			}
		}
		root.Path = pub //覆盖原有值设置公共根
		return true
	}
	return false
}

// FindPublicRoot 查找公共前缀，如无公共前缀则返回 ""
func (r *route) findPublicRoot(method, p1, p2, path string) string {
	l := len(p1)
	if l > len(p2) {
		l = len(p2) //取长度短的
	}
	index := -1
	for i := 0; i <= l && p1[:i] == p2[:i]; i++ { //此处可能发生bug
		index = i
	}
	if index > 0 && index <= l {
		//提取公共根 遇到REST API时候 需要阻止提取  主要修改 /aaa/{name} 和 /aaa/{nme} 情况下会造成提取公共根 /aaa/{n, 会造成破坏RESTFul的单个节点完整性
		s1 := p1[:index]
		// 检擦 最后一个路径分割符 是否存在不完整的 RESTFul
		//找到最 后一个 /
		if lastIndex := strings.LastIndex(s1, "/"); lastIndex != -1 {
			//注册开始对 路径已经进行过校验处理，此处的 url 是标准规范
			temp := s1[lastIndex:]
			lb := strings.Count(temp, "{")
			rb := strings.Count(temp, "}")
			if lb != rb {
				//完整性校验失败 该路径注册会失败,出现完整性 校验失败的 直接 结束程序，后续逻辑无法继续
				r.Error(method, " : ", path, " RESTFul Integrity check failed with conflict")
				os.Exit(-1)
				return ""
			}
		}
		return s1
	}
	return ""
}

// urlRouter 检索指定的path路由
// method 请求类型，path 查询路径，rw，req http生成的请求响应,
// ctx 中间件请求上下文参数
func (r *route) urlRouter(method, path string, rw http.ResponseWriter, req *http.Request, ctx Ctx) (*node, []string, map[string]interface{}, Ctx) {
	if ctx == nil {
		ctx = make(Ctx)
		ctx[request] = req
		ctx[response] = rw
		ctx[iocs] = r.component
		ctx[auroraMaxMultipartMemory] = r.MaxMultipartMemory
	}
	// 全局中间件
	middlewares := r.middleware
	if middlewares != nil {
		for _, middleware := range middlewares {
			if middleware == nil {
				continue
			}
			if f := middleware(ctx); !f {
				return nil, nil, nil, nil
			}
		}
	}
	if r.isStatic(path, rw, req) {
		return nil, nil, nil, nil
	}
	if index := strings.Index(path, "."); index != -1 {
		path = r.fileService
	}
	//查找指定的Method树
	if _, ok := r.tree[method]; !ok {
		http.NotFound(rw, req)
		return nil, nil, nil, nil
	}
	c, u, args := r.bfs(r.tree[method], path)
	if c == nil {
		http.NotFound(rw, req)
		return nil, nil, nil, nil
	}
	return c, u, args, ctx
}

// 路由树查询
func (r *route) bfs(root *node, path string) (*node, []string, map[string]interface{}) {
	var next *element
	reqCount := strings.Count(path, "/")
	q := queue{}
	q.en(root)
star:
	next = q.next()
	for next != nil {
		n := next.value
		if n.Control != nil && reqCount == n.Count {
			if !strings.Contains(n.FullPath, "{") {
				if path == n.FullPath {
					return n, nil, nil
				}
			} else {
				urlArgs, Aargs := analysisRESTFul(n, path)
				if urlArgs == nil {
					goto next
				}
				return n, urlArgs, Aargs
			}
		}
	next:
		child := n.Child
		if child != nil {
			for i := 0; i < len(child); i++ {
				q.en(child[i])
			}
		}
		goto star
	}
	return nil, nil, nil
}

// ServeHTTP 一切的开始
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c, u, args, ctx := engine.router.urlRouter(req.Method, req.URL.Path, rw, req, nil)
	if c == nil {
		return
	}
	engine.handle(c, u, args, rw, req, ctx)

}

// 请求处理
func (engine *Engine) handle(c *node, u []string, args map[string]interface{}, rw http.ResponseWriter, req *http.Request, ctx Ctx) {
	p := engine.proxyPool.Get().(*Proxy)
	p.Rew = rw
	p.Req = req
	p.Ctx = ctx
	p.control = *c.Control
	p.middleware = c.middleware
	p.UrlVariable = u
	p.args = args
	p.view = engine
	p.Engine = engine
	p.start()
	engine.proxyPool.Put(p)
}

// 获取注册接口函数名称
func getFunName(fun interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
}

// 处理注册函数路径的开头,去除多余的 / 开头
func urlHead(url string) string {
	if url == "/" {
		return url
	}
	if url[:1] != "/" {
		url = "/" + url
	}
	ul := len(url)
	index := 0
	for i := 0; i < ul && url[i:i+1] == "/"; i++ {
		index = i
	}
	if index == 0 {
		return url
	}
	return url[index:]
}

// 处理请求接口 后缀
func urlEnd(url string) string {
	if url == "/" {
		return url
	}
	ul := len(url)
	if url[ul-1:] != "/" {
		return url
	}
	index := ul
	for i := ul - 1; i >= 0 && url[i-1:i] == "/"; i-- {
		index = i
	}
	if index == 0 {
		return url
	}
	return url[:index-1]
}

// 检查 请求接口 格式规范
func urlCheck(url string) string {
	re := regexp.MustCompile(`/{2,}`)
	all := re.FindAll([]byte(url), -1)
	if len(all) > 0 {
		for _, r := range all {
			s := string(r)
			url = strings.Replace(url, s, "/", 1)
		}
	}
	return url
}

// checkRESTFul 校验注册路径的RESTFul
func checkRESTFul(url string) error {
	if strings.Contains(url, "{}") { //此处的校验还需要加强，单一判断{}存在其他风险，开发者要么自己不能出现一些其他问题，比如 ...{}ss/.. or  .../a{s}a/.. 等情况 发现时间: 2022.1.5
		return errors.New(url + " RESTFul cannot be empty {}")
	}
	//检查 完整性
	l := strings.Count(url, "{")
	r := strings.Count(url, "}")
	// 通过 括号数量检查是否成对
	if l != r {
		return errors.New(url + " RESTFul integrity check failed")
	}
	// 加上 一个 / 主要是用于校验 /{sss}/aaa/{vv}  {xxx}结尾的辅助
	temp := url + "/"
	re := regexp.MustCompile(`/{[a-z]*[A-Z]*\d*}`)
	re2 := regexp.MustCompile(`{[a-z]*[A-Z]*\d*}/`)
	all := re.FindAll([]byte(temp), -1)
	all2 := re2.FindAll([]byte(temp), -1)
	if len(all) != len(all2) {
		return errors.New(url + " RESTFul integrity check failed")
	}
	return nil
}

// analysisRESTFul 解析路径参数
// n 路由节点
// mapping 前端请求路径
func analysisRESTFul(n *node, mapping string) ([]string, map[string]interface{}) {
	FullPath := n.FullPath
	reg := regexp.MustCompile(`{*[a-z]*[A-Z]*\d*}*`)
	req := reg.FindAll([]byte(mapping), -1)
	res := reg.FindAll([]byte(FullPath), -1)
	if len(req) != len(res) {
		return nil, nil
	}
	urls := make([]string, 0)
	args := make(map[string]interface{})
	for i := 1; i < len(req); i++ {
		rest := string(res[i])
		reqp := string(req[i])
		if !strings.Contains(rest, "{") {
			if rest != reqp {
				return nil, nil
			}
			continue
		}
		urls = append(urls, reqp)
		// 重名RESTFul 参数将被覆盖
		args[rest[1:len(rest)-1]] = reqp
	}
	return urls, args
}

// isStatic 处理静态资源 返回true 表示处理了静态资源
func (r *route) isStatic(path string, rw http.ResponseWriter, req *http.Request) bool {
	mapping := path
	if index := strings.LastIndex(req.URL.Path, "."); index != -1 { //此处判断这个请求可能为静态资源处理
		// 文件服务器校验
		if strings.HasPrefix(path, r.fileService) {
			return false
		}
		t := req.URL.Path[index:] //截取可能的资源类型
		r.resourceHandler(rw, req, mapping, t)
		return true
	}
	return false
}
