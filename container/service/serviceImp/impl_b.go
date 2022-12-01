package serviceImp

import "gitee.com/aurora-engine/aurora/container/service"

type Bbb struct {
	service.A `impl:"gitee.com/aurora-engine/aurora/container/service/serviceImp-*serviceImp.Aaa"`
}
