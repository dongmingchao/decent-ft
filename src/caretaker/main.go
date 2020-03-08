package caretaker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	resource_pool "github.com/dongmingchao/decent-ft/src/resource-pool"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

func watchDir(dirname string, handler func(event fsnotify.Event)) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

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
	return watcher
}

func (watcher *careWatcher) watchHandler(event fsnotify.Event) {
	length := len(watcher.fileNames)
	id := sort.SearchStrings(watcher.fileNames, event.Name)
	println("search file", id)
	//if event.Op&fsnotify.Create == fsnotify.Create {
	//	if id == length {
	//	}
	//}
	if event.Op&fsnotify.Write == fsnotify.Write {
		if id == length {
			watcher.stashAppend(event.Name)
		} else {
			gfile := gStash.Files[id]
			oldMarkStr := fmt.Sprintf("%x", gfile.Checksum)
			readFile(event.Name, func(file *os.File) {
				obj := stashFile(file)
				gfile.Checksum = obj.Mark
			})
			hashDir := stashDir+"/"+oldMarkStr[0:2]
			stashPath := hashDir+"/"+oldMarkStr[2:38]
			os.Remove(stashPath)
			println("remove", stashPath)
			_ = os.Remove(hashDir)
		}
	}
	if event.Op&fsnotify.Rename == fsnotify.Rename {

	}
		//f, err := os.Open(event.Name)
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	text, _ := ioutil.ReadAll(f)
	//	println(string(text))
	//	f.Close()
	//}
}

func (watcher *careWatcher) stashAppend(filename string) {
	var gfile resource_pool.GFile
	readFile(filename, func(file *os.File) {
		fName := file.Name()
		obj := stashFile(file)
		gfile = resource_pool.GFile{
			FileName: fName,
			FileNameLen: uint16(len(fName)),
			Checksum: obj.Mark,
		}
	})
	gStash.Files = append(gStash.Files, &gfile)
	watcher.fileNames = append(watcher.fileNames, filename)
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

func removeStash() {

}

const (
	focusDir       = "./src/resource-pool/sample-pool"
	stashDir       = "./objects"
	stashIndexFile = stashDir + "/index"
)

func ReadIndex() resource_pool.GTree {
	stash := resource_pool.GTree{}
	if _, err := os.Stat(stashDir); os.IsNotExist(err) {
		os.Mkdir(stashDir, os.ModeDir|0700)
		fmt.Println("Create Stash Dir: ", stashDir)
	}
	if _, err := os.Stat(stashIndexFile); os.IsNotExist(err) {
		os.Create(stashIndexFile)
		fmt.Println("Create Stash Index file: ", stashIndexFile)
		stash.Version = 1
	} else {
		stashIndex, _ := os.Open(stashIndexFile)
		stash.Read(stashIndex)
		stashIndex.Close()
		fmt.Println("Found Stash Index file")
		fmt.Println(stash)
	}
	return stash
}

func SaveIndex(stash resource_pool.GTree) {
	//var fileHashSet [][20]byte
	//for ei, ef := range stash.Files {
	//	fileHashSet[ei] = ef.Checksum
	//}

	//dir, err := ioutil.ReadDir(focusDir)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, each := range dir {
	//	fileName := focusDir + "/" + each.Name()
	//}
	stash.FileCount = uint32(len(stash.Files))

	var allBytes bytes.Buffer
	allBytes.Write(UInt32ToBytes(stash.Version))
	allBytes.Write(UInt32ToBytes(stash.FileCount))
	binary.Write(&allBytes, binary.BigEndian, stash.Files)
	println(allBytes.String())
	stash.Checksum = resource_pool.Sha1CheckSum(allBytes.Bytes())
	fmt.Println(stash)

	stashIndex, _ := os.OpenFile(stashIndexFile, os.O_CREATE | os.O_RDWR, 0644)
	stash.Write(stashIndex)
	stashIndex.Close()
}
var gStash resource_pool.GTree
type careWatcher struct {
	fileNames []string
}

func newCareWatcher(stash resource_pool.GTree) *careWatcher {
	var watcher careWatcher
	for _, file := range stash.Files {
		watcher.fileNames = append(watcher.fileNames, file.FileName)
	}
	return &watcher
}



func WatchDir() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var fileWatcher *fsnotify.Watcher
	go func() {
		gStash = ReadIndex()
		watcher := newCareWatcher(gStash)
		fileWatcher = watchDir(focusDir, watcher.watchHandler)
	}()
	fmt.Println("[File Watcher] Start")
	<-sigs
	err := fileWatcher.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[File Watcher] Stop")
	SaveIndex(gStash)
}

func UInt32ToBytes(i uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}
