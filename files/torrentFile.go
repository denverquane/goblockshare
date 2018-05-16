package files

import (
	"crypto/sha256"
	"os"
	"strconv"
)

type TorrentFile struct {
	SegmentByteSize int
	SegmentHashes   [][]byte
	RawData         []byte
}

var kilobyte = 1000
var megabyte = 1000000

func (torr TorrentFile) ToString() string {
	return "torrent size: " + strconv.FormatInt(int64(torr.SegmentByteSize), 10) + " data: " + string(torr.RawData)
}

func MakeTorrentFileFromFile(segByteSize int, url string) (TorrentFile, error) {
	file, err := os.Open(url)
	defer file.Close()
	if err != nil {
		return TorrentFile{}, err
	}

	torr := TorrentFile{segByteSize, make([][]byte, 0), make([]byte, 0)}
	readbytes := segByteSize

	for offset := int64(0); readbytes == segByteSize; {
		buffer := make([]byte, segByteSize)
		readbytes, _ = file.ReadAt(buffer, offset)
		if err != nil {
			return TorrentFile{}, err
		}
		torr.appendNewSegment(buffer[0:readbytes])
		offset += int64(readbytes)
		//fmt.Println("Read " + strconv.FormatInt(off / 1000, 10) + " kilobytes so far")
	}
	return torr, nil
}

func AreSameTorrentBytes(segByteSize int, fileA []byte, fileB []byte) bool {
	if len(fileA) != len(fileB) || segByteSize > len(fileA){
		return false
	}

	var readBytes int

	if len(fileA) % segByteSize == 0 {
		readBytes = segByteSize
	} else {
		readBytes = len(fileA) % segByteSize
	}

	for offset := int64(0); offset < int64(len(fileA)); {
		segA := fileA[offset : offset + int64(segByteSize)]
		segB := fileB[offset : offset + int64(segByteSize)]
		if !doSegmentsHashToSame(segA, segB) {
			return false
		}
		offset += int64(readBytes)
	}
	return true
}

func (torr TorrentFile) ValidateHashes() bool {
	segs := arrToSegments(torr.SegmentByteSize, torr.RawData)

	for i, v := range torr.SegmentHashes {
		hash := hashSegment(segs[i])
		for ii, vv := range v {
			if hash[ii] != vv {
				return false
			}
		}
	}
	return true
}

func arrToSegments(size int, arr []byte) [][]byte {
	doubleArr := make([][]byte, 0)
	var offset int
	for offset = 0; offset+size < len(arr); offset += size {
		doubleArr = append(doubleArr, arr[offset : offset + size])
	}
	doubleArr = append(doubleArr, arr[offset : ])
	return doubleArr
}

func doSegmentsHashToSame(segA []byte, segB []byte) bool {
	if len(segA) != len(segB) {
		return false
	}

	hashA := hashSegment(segA)
	hashB := hashSegment(segB)

	for i, v := range hashA {
		if v != hashB[i] {
			return false
		}
	}
	return true
}

func (file *TorrentFile) appendNewSegment(segData []byte) {
	hashed := hashSegment(segData)
	file.SegmentHashes = append(file.SegmentHashes, hashed)
	//_, ok := file.segmentHashMap[string(hashed)]
	//if ok {
	//	file.segmentHashMap[string(hashed)]++
	//} else {
	//	file.segmentHashMap[string(hashed)] = 1
	//}
	file.RawData = append(file.RawData, segData...)
}

func hashSegment(seg []byte) []byte {
	h := sha256.New()
	h.Write(seg)
	return h.Sum(nil)
}
