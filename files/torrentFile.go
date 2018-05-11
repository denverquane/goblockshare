package files

import (
	"crypto/sha256"
	"os"
)

type TorrentFile struct {
	segmentByteSize int
	segmentData     [][]byte
	segmentHashes   [][]byte
	rawData         []byte
	segmentHashMap  map[string]int64
}

var kilobyte = 1000
var megabyte = 1000000

func MakeTorrentFileFromFile(segsize int, url string) (error, TorrentFile) {
	file, err := os.Open(url)
	defer file.Close()
	if err != nil {
		return err, TorrentFile{}
	}

	torr := TorrentFile{segsize, make([][]byte, 0), make([][]byte, 0), make([]byte, 0), make(map[string]int64)}
	readbytes := segsize

	for off := int64(0); readbytes == segsize; {
		buffer := make([]byte, segsize)
		readbytes, _ = file.ReadAt(buffer, off)
		if err != nil {
			return err, TorrentFile{}
		}
		torr.appendNewSegment(buffer[0:readbytes])
		off += int64(readbytes)
		//fmt.Println("Read " + strconv.FormatInt(off / 1000, 10) + " kilobytes so far")
	}
	return nil, torr
}

func (file *TorrentFile) appendNewSegment(segData []byte) {
	file.segmentData = append(file.segmentData, segData)
	hashed := hashSegment(segData)
	file.segmentHashes = append(file.segmentHashes, hashed)
	_, ok := file.segmentHashMap[string(hashed)]
	if ok {
		file.segmentHashMap[string(hashed)]++
	} else {
		file.segmentHashMap[string(hashed)] = 1
	}
	file.rawData = append(file.rawData, segData...)
}

func hashSegment(seg []byte) []byte {
	h := sha256.New()
	h.Write(seg)
	return h.Sum(nil)
}
