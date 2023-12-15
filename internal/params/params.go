package params

import (
	"flag"
	"sync"
)

type Params struct {
	Version bool
}

var params *Params
var once sync.Once

func GetParams() *Params {
	once.Do(func() {
		params = &Params{}
		params.init()
	})
	return params
}

func (params *Params) init() {
	// 解析
	flag.BoolVar(&params.Version, "v", false, "打印版本，然后直接退出。可以在 template 存活时候，再次用 -v 参数启动 template。"+" (shorthand)")
	flag.BoolVar(&params.Version, "version", false, "打印版本，然后直接退出。可以在 template 存活时候，再次用 -v 参数启动 template。")
	flag.Parse()
}
