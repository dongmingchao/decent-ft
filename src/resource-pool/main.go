package resource_pool

//Go 在多个 crypto/* 包中实现了一系列散列函数。
import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/dongmingchao/decent-ft/JSlike"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

/**
echo "blob 13\0test content" | shasum
d670460b4b4aece5915caf5c68d12f560a9fe3e4
echo 'test content' | git hash-object -w --stdin
git cat-file -p d670460b4b4aece5915caf5c68d12f560a9fe3e4
test content
*/

const (
	Blob   = "blob"
	Commit = "commit"
)

type GHash struct {
	GType    string // blob | commit
	GLen     int
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
		GLen:  bin.Len(),
		GBin:  bin.Bytes(),
	}
	hash.FullBody = makeFullBody(hash.GType, hash.GLen, hash.GBin)
	hash.Mark = Sha1CheckSum(hash.FullBody.Bytes())
	hash.MarkStr = fmt.Sprintf("%x", hash.Mark)
	return hash
}

func makeFullBody(gType string, gLen int, gBin []byte) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(gType)
	buffer.WriteString(" ")
	buffer.WriteString(strconv.Itoa(gLen))
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

type GFile struct {
	ctime       uint64 // 提交时间的16进制秒数
	mtime       uint64
	fileSize    uint32
	FileNameLen uint16
	FileName    string // 变长，为fileNameLen
	Checksum    [20]byte
}

func (file GFile) String() string {
	s, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(s)
}

func (file GFile) MarshalJSON() ([]byte, error){
	return json.Marshal(JSlike.Object{
		"FileNameLen": file.FileNameLen,
		"FileName": file.FileName,
		"Checksum": fmt.Sprintf("%x", file.Checksum),
	})
}

func (file GFile) Write(w io.Writer) {
	binary.Write(w, binary.BigEndian, &file.FileNameLen)
	binary.Write(w, binary.BigEndian, []byte(file.FileName))
	binary.Write(w, binary.BigEndian, &file.Checksum)
}

func (file *GFile) Read(r io.Reader) {
	binary.Read(r, binary.BigEndian, &file.FileNameLen)
	var filename = make([]byte, file.FileNameLen)
	binary.Read(r, binary.BigEndian, &filename)
	file.FileName = string(filename)
	binary.Read(r, binary.BigEndian, &file.Checksum)
}

type GTree struct {
	Version   uint32
	FileCount uint32
	Files     []*GFile
	Checksum  [20]byte
}

func (t GTree) MarshalJSON() ([]byte, error){
	return json.Marshal(JSlike.Object{
		"Version": t.Version,
		"FileCount": t.FileCount,
		"Files": t.Files,
		"Checksum": fmt.Sprintf("%x", t.Checksum),
	})
}

func (t GTree) String() string {
	s, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(s)
}

func (t GTree) Write(w io.Writer) {
	binary.Write(w, binary.BigEndian, &t.Version)
	binary.Write(w, binary.BigEndian, &t.FileCount)
	for _, each := range t.Files {
		//fmt.Println(each)
		each.Write(w)
	}
	binary.Write(w, binary.BigEndian, &t.Checksum)
}

func (t *GTree) Read(r io.Reader) {
	binary.Read(r, binary.BigEndian, &t.Version)
	binary.Read(r, binary.BigEndian, &t.FileCount)
	t.Files = make([]*GFile, t.FileCount)
	for i := uint32(0); i < t.FileCount; i++ {
		t.Files[i].Read(r)
	}
	binary.Read(r, binary.BigEndian, &t.Checksum)
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
