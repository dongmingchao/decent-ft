package resourcePool

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

/**
echo "blob 13\0test content" | shasum
d670460b4b4aece5915caf5c68d12f560a9fe3e4
echo 'test content' | git hash-object -w --stdin
git cat-file -p d670460b4b4aece5915caf5c68d12f560a9fe3e4
test content
*/
type FileType byte

const (
	Blob   FileType = 1 << iota
	Neighbor
)

func (t FileType) String() string {
	switch t {
	case Blob:
		return "blob"
	case Neighbor:
		return "neighbor"
	}
	return "unknown"
}

type GHash struct {
	GType    FileType
	GLen     uint32
	GBin     []byte
	Mark     [20]byte
	MarkStr  string
	FullBody bytes.Buffer
}

func (hash GHash) String() string {
	return hash.FullBody.String()
}

func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func (hash GHash) GitRecordFile() {
	sHash := []byte(hash.String())
	fmt.Printf("zlib: % x\n", DoZlibCompress(sHash))
	f, err := os.OpenFile("./test/.git/objects/d6/70460b4b4aece5915caf5c68d12f560a9fe3e4", os.O_CREATE|os.O_RDWR, 0444)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func NewGHash(content []byte) GHash {
	var bin bytes.Buffer
	bin.Write(content)
	bin.WriteByte('\n')
	hash := GHash{
		GType: Blob,
		GLen:  uint32(bin.Len()),
		GBin:  bin.Bytes(),
	}
	fLength := IntTo3Bytes(bin.Len())
	hash.FullBody = makeFullBody(hash.GType, fLength[:], hash.GBin)
	hash.Mark = Sha1CheckSum(hash.FullBody.Bytes())
	hash.MarkStr = fmt.Sprintf("%x", hash.Mark)
	return hash
}

func makeFullBody(gType FileType, gLen []byte, gBin []byte) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteByte(byte(gType))
	buffer.Write(gLen)
	buffer.WriteByte('\000')
	buffer.Write(gBin)
	return buffer
}

func Sha1CheckSum(u []byte) [20]byte {
	h := sha1.New()
	h.Write(u)
	res := h.Sum(nil)
	var ret [20]byte
	copy(ret[:], res[:20])
	return ret
}

/**
Git 暂存区 .git/index
*/
type GitIndex struct {
	DIRC      string // "DIRC"
	version   uint32
	fileCount uint32

	ctime        uint64 // 提交时间的16进制秒数
	mtime        uint64
	uid          uint32
	gid          uint32
	fileSize     uint32
	fileChecksum [20]byte
	flags        uint
	/**
	flagAssumeUnchanged 1 bit
	flagExtended        1 bit // version 2之后为0
	flagStage           2 bit // 0 regular 1 base 2 ours 3 theirs
	*/
	fileNameLen [3]uint
	// extends     [2]byte // version 2之后取消了
	fileName []byte // 变长，为fileNameLen

	zero     []byte // 1-8 padding
	checksum [20]byte
}

func main() {
	f, err := os.OpenFile("./test/.git/objects/16/ebff13d6420ad7f81a6ee0614d788b28f9e4c0", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	content, _ := ioutil.ReadAll(f)
	content = DoZlibUnCompress(content)
	f.Close()
	head := make([]byte, 24)
	bytes.NewReader(content).Read(head)
	gt := bytes.Split(head, []byte{'\000'})[0]
	println(string(gt))
	//outf, err := os.OpenFile("./output.zip")
	//
	//println(string())
}
