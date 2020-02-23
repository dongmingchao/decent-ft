package caretaker

import (
	resource_pool "decent-ft/src/resource-pool"
	"encoding/binary"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
)

func watchDir(dirname string, handler func(event fsnotify.Event)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				handler(event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dirname) // "./src/resource-pool/sample-pool"
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func watchHandler(event fsnotify.Event) {
	f, err := os.OpenFile(event.Name, os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
	} else {
		text, _ := ioutil.ReadAll(f)
		println(string(text))
		f.Close()
	}
}

func readFile(fileName string, cb func(*os.File)) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		log.Println(err)
	} else {
		cb(f)
		f.Close()
	}
}

func stashFile(file *os.File) resource_pool.GHash {
	text, _ := ioutil.ReadAll(file)
	obj := resource_pool.NewGHash(text)
	os.MkdirAll(stashDir+"/"+obj.MarkStr[0:2], os.ModeDir|0700)
	f, err := os.OpenFile(stashDir+"/"+obj.MarkStr[0:2]+"/"+obj.MarkStr[2:38], os.O_CREATE|os.O_RDWR, 0644)
	binary.Write(f, binary.BigEndian, obj.FullBody.Bytes())
	if err != nil {
		log.Println(err)
	}
	return obj
}

const (
	focusDir       = "./src/resource-pool/sample-pool"
	stashDir       = "./objects"
	stashIndexFile = stashDir + "/index"
)

func main() {
	if _, err := os.Stat(stashDir); os.IsNotExist(err) {
		os.Mkdir(stashDir, os.ModeDir|0700)
	}
	stash := resource_pool.GTree{}
	if _, err := os.Stat(stashIndexFile); os.IsNotExist(err) {
		os.Create(stashIndexFile)
		stash.Version = 1
	}
	stashIndex, _ := os.Open(stashIndexFile)
	stash.Read(stashIndex)
	stashIndex.Close()
	fmt.Println(stash)
	//fileHashSet := [][20]byte{}
	//for ei, ef := range stash.Files {
	//	fileHashSet[ei] = ef.Checksum
	//}

	//dir, err := ioutil.ReadDir(focusDir)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, each := range dir {
	//	fileName := focusDir + "/" + each.Name()
	//	gfile := resource_pool.GFile{
	//		FileName: fileName,
	//		FileNameLen: uint16(len(fileName)),
	//	}
	//	readFile(fileName, func(file *os.File) {
	//		obj := stashFile(file)
	//		gfile.Checksum = obj.Mark
	//	})
	//	stash.Files = append(stash.Files, gfile)
	//}
	//stash.FileCount = uint32(len(stash.Files))
	//
	//var allBytes bytes.Buffer
	//allBytes.Write(UInt32ToBytes(stash.Version))
	//allBytes.Write(UInt32ToBytes(stash.FileCount))
	//binary.Write(&allBytes, binary.BigEndian, stash.Files)
	//println(allBytes.String())
	//stash.Checksum = resource_pool.Sha1CheckSum(allBytes.Bytes())
	//fmt.Println(stash)
	//
	//stashIndex, _ = os.OpenFile(stashIndexFile, os.O_CREATE | os.O_RDWR, 0644)
	//stash.Write(stashIndex)
	//stashIndex.Close()

	//dir, _ = ioutil.ReadDir(stashDir)
	//for _, each := range dir {
	//	println(each.Name())
	//}
	//watchDir(focusDir, watchHandler)
}

func UInt32ToBytes(i uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}
