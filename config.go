package aurora

import (
	"fmt"
	"gitee.com/aurora-engine/aurora/cnf"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

const (
	yml  = "application.yml"
	yaml = "application.yaml"
)


// viperConfig 配置并加载 application.yml 配置文件
func (engine *Engine) viperConfig() {
	var ConfigPath string
	var err error
	if engine.configpath == "" {
		// 扫描配置文件
		filepath.WalkDir(engine.projectRoot, func(p string, d fs.DirEntry, err error) error {
			//找到配置及文件,基于根目录优先加载最外层的application.yml
			if !d.IsDir() && (strings.HasSuffix(p, yml) || (strings.HasSuffix(p, yaml))) && ConfigPath == "" {
				//修复 项目加载配置覆盖，检索项目配置文件，避免内层同名配置文件覆盖外层，这个情况可能发生在 开发者把两个go mod 项目嵌套在一起，导致配置被覆盖
				//此处校验，根据检索的更路径，只加载最外层的配置文件
				ConfigPath = p
			}
			return nil
		})
	} else {
		ConfigPath = engine.configpath
	}
	if ConfigPath == "" {
		engine.config = &cnf.ConfigCenter{viper.New(), &sync.RWMutex{}}
		return
	}
	if engine.config == nil {
		// 用户没有提供 配置项 则创建默认的配置处理
		cnf := &cnf.ConfigCenter{
			viper.New(),
			&sync.RWMutex{},
		}
		cnf.SetConfigFile(ConfigPath)
		err = cnf.ReadInConfig()
		ErrorMsg(err)
		engine.config = cnf
	}
	// 加载基础配置
	if engine.config != nil {                      //是否加载配置文件 覆盖配置项
		engine.Info("the configuration file is loaded successfully.")
		// 读取web服务端口号配置
		port := engine.config.GetString("aurora.server.port")
		if port != "" {
			engine.port = port
		}
		// 读取静态资源配置路径
		engine.resource = "/"
		p := engine.config.GetString("aurora.resource")
		// 构建路径拼接，此处在路径前后加上斜杠 用于静态资源的路径凭借方便
		if p != "" {
			if p[:1] != "/" {
				p = "/" + p
			}
			if p[len(p)-1:] != "/" {
				p = p + "/"
			}
			engine.resource = p
		}
		// 读取文件服务配置
		p = engine.config.GetString("aurora.server.file")
		engine.fileService = p
		engine.Info(fmt.Sprintf("server static resource root directory:%1s", engine.resource))
		// 读取服务名称
		name := engine.config.GetString("aurora.application.name")
		if name != "" {
			engine.name = name
			engine.Info("the service name is " + engine.name)
		}
	}
}

// GetConfig 获取 Aurora 配置实例 对配置文件内容的读取都是协程安全的
func (engine *Engine) GetConfig() cnf.Config {
	return engine.config
}
