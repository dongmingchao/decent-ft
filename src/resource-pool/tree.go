package resourcePool

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/dongmingchao/decent-ft/JSlike"
	"io"
	"log"
)

type GTree struct {
	Version   uint32
	FileCount uint32
	Files     []GBaseInterface
	Checksum  [20]byte
}

func (t GTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSlike.Object{
		"Version":   t.Version,
		"FileCount": t.FileCount,
		"Files":     t.Files,
		"Checksum":  fmt.Sprintf("%x", t.Checksum),
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
		each.Write(w)
	}
	binary.Write(w, binary.BigEndian, &t.Checksum)
}

func (t *GTree) Read(r io.Reader) {
	binary.Read(r, binary.BigEndian, &t.Version)
	binary.Read(r, binary.BigEndian, &t.FileCount)
	t.Files = make([]GBaseInterface, t.FileCount)
	for i := uint32(0); i < t.FileCount; i++ {
		t.Files[i].Read(r)
	}
	binary.Read(r, binary.BigEndian, &t.Checksum)
}
