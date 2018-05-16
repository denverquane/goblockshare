package files

import (
	"crypto/sha256"
	"os"
	"strconv"
	"encoding/hex"
)

type TorrentFile struct {
	SegmentByteSize int

	SegmentHashKeys   []string
	SegmentHashMap	map[string][]byte
}

var kilobyte = 1000
var megabyte = 1000000

func (torr TorrentFile) ToString() string {
	a := "torrent segment size: " + strconv.FormatInt(int64(torr.SegmentByteSize), 10) + "\n"
	for _, v := range torr.SegmentHashKeys {
		a += v + "\n"
	}
	return a
}

func MakeTorrentFileFromFile(segByteSize int, url string) (TorrentFile, error) {
	file, err := os.Open(url)
	defer file.Close()
	if err != nil {
		return TorrentFile{}, err
	}

	torr := TorrentFile{segByteSize, make([]string, 0), make(map[string][]byte, 0)}
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
	for hash, raw := range torr.SegmentHashMap {
		if hex.EncodeToString(hashSegment(raw)) != hash {
			return false
		}
	}
	return true
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
	hexHashed := hex.EncodeToString(hashSegment(segData))
	file.SegmentHashKeys = append(file.SegmentHashKeys, hexHashed)
	file.SegmentHashMap[hexHashed] = segData
}

func hashSegment(seg []byte) []byte {
	h := sha256.New()
	h.Write(seg)
	return h.Sum(nil)
}
