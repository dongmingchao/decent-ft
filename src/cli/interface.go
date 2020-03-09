package cli

import (
	"bytes"
	resourcePool "github.com/dongmingchao/decent-ft/src/resource-pool"
)

type Cli interface {
	cat(filename string) (bytes.Buffer, error)
	list() []resourcePool.GFile
	tree() resourcePool.GTree
	stat(filename string) (resourcePool.GFile, error)
}