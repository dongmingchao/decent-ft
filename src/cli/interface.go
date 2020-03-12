package cli

import (
	"bytes"
	resourcePool "github.com/dongmingchao/decent-ft/src/resource-pool"
)

type Cli interface {
	// 获取文件内容
	cat(string) (bytes.Buffer, error)
	// 列出文件路径集合
	list() []resourcePool.GFile
	// 列出索引树
	tree() resourcePool.GTree
	// 显示文件（blob）/邻居网络（neighbor）的详细信息
	stat(string) (resourcePool.GFile, error)
}