package resourcePool

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/dongmingchao/decent-ft/JSlike"
	"io"
	"log"
)

type GFile struct {
	FileSize uint32 // 存储文件时文件大小占24位（16MB），但是在内存中时占32位（8G）
	*GBase
}

func (file GFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSlike.Object{
		"FileNameLen": file.FileNameLen,
		"FileName":    file.FileName,
		"Type":        file.FileType.String(),
		"Size":        file.FileSize,
		"Checksum":    fmt.Sprintf("%x", file.Checksum),
	})
}

func (file GFile) Write(w io.Writer) {
	binary.Write(w, binary.BigEndian, file.FileNameLen)
	binary.Write(w, binary.BigEndian, []byte(file.FileName))
	binary.Write(w, binary.BigEndian, file.FileType)
	if file.FileSize > 16777215 {// 0xffffff 24位最大
		log.Println(file.FileName, "文件超过最大单元长度")
	} else {
		binary.Write(w, binary.BigEndian, IntTo3Bytes(int(file.FileSize)))
	}
	binary.Write(w, binary.BigEndian, file.Checksum)
}

func (file *GFile) Read(r io.Reader) {
	binary.Read(r, binary.BigEndian, &file.FileNameLen)
	var filename = make([]byte, file.FileNameLen)
	binary.Read(r, binary.BigEndian, &filename)
	file.FileName = string(filename)
	binary.Read(r, binary.BigEndian, &file.FileType)
	var size [3]byte
	binary.Read(r, binary.BigEndian, &size)
	file.FileSize = BytesToUInt32(size[:])
	binary.Read(r, binary.BigEndian, &file.Checksum)
}
