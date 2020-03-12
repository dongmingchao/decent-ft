package resourcePool

import (
	"encoding/json"
	"io"
	"log"
)

type GBase struct {
	ctime    uint64 // 提交时间的16进制秒数
	mtime    uint64
	FileType
	FileNameLen uint16 // 文件名限长0xffff
	FileName    string // 变长，为fileNameLen
	Checksum    [20]byte
}

type GBaseInterface interface {
	MarshalJSON() ([]byte, error)
	Read(r io.Reader)
	Write(w io.Writer)
}

func (b GBase) String() string {
	s, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(s)
}
