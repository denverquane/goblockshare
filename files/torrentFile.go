package files

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
)

type TorrentFile struct {
	Name 			string
	SegmentByteSize int
	TotalHash       string
	//Total hash is dependent on the name, the segment size, the raw data, and nothing else
	//-the other aspects of a torrent, such as the hashed duplicates and ordering of segments should be LOCAL
	//(we already hash the raw data in order)

	SegmentHashKeys []string //fine to expose publicly to say what layers the torrent has

	segmentHashMaps map[string][]byte //don't expose; reveals entire torrent structure
	duplicatesMaps  map[string]int //don't reveal; should only be used internally/locally
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

func (torr TorrentFile) ToHash() []byte {
	h := sha256.New()

	for _, v := range torr.SegmentHashKeys {
		h.Write([]byte(v))
		h.Write(torr.segmentHashMaps[v]) //get the raw bytes associated with the hash
	}

	return h.Sum(nil)
}

func MakeTorrentFileFromFile(segByteSize int, url string, name string) (TorrentFile, error) {
	file, err := os.Open(url)
	defer file.Close()
	if err != nil {
		return TorrentFile{}, err
	}

	torr := TorrentFile{name, segByteSize, "", make([]string, 0),
		make(map[string][]byte, 0), make(map[string]int)}
	readbytes := segByteSize
	total := sha256.New()
	total.Write([]byte(name))
	//TODO write segbytesize (dont allow torrents to change segmentation size)
	for offset := int64(0); readbytes == segByteSize; {
		buffer := make([]byte, segByteSize)
		readbytes, _ = file.ReadAt(buffer, offset)
		if err != nil {
			return TorrentFile{}, err
		}
		torr.appendNewSegment(buffer[0:readbytes])
		total.Write(buffer[0:readbytes])
		offset += int64(readbytes)
		//fmt.Println("Read " + strconv.FormatInt(off / 1000, 10) + " kilobytes so far")
	}
	torr.TotalHash = hex.EncodeToString(total.Sum(nil))
	return torr, nil
}

func MakeTorrentFromBytes(segByteSize int, data []byte, name string) (TorrentFile, error) {
	if segByteSize > len(data) {
		return TorrentFile{}, errors.New("Segment too long")
	}

	torr := TorrentFile{name, segByteSize, "", make([]string, 0), make(map[string][]byte, 0), make(map[string]int),}

	var offset int
	total := sha256.New()
	total.Write([]byte(name))
	//TODO write segbytesize (dont allow torrents to change segmentation size)
	for offset = 0; offset+segByteSize < len(data); {
		segment := data[offset : offset+segByteSize]

		offset += segByteSize
		torr.appendNewSegment(segment)
		total.Write(segment)
	}
	torr.appendNewSegment(data[offset:])
	total.Write(data[offset:])
	torr.TotalHash = hex.EncodeToString(total.Sum(nil))
	return torr, nil
}

func (torr1 TorrentFile) Equals(torr2 TorrentFile) bool {
	h1 := torr1.ToHash()
	h2 := torr2.ToHash()

	for i, v := range h1 {
		if v != h2[i] {
			return false
		}
	}
	return true
}

func (torr TorrentFile) Validate() bool {
	for hash, raw := range torr.segmentHashMaps {
		if hex.EncodeToString(hashSegment(raw)) != hash {
			return false
		}
	}
	return true
}

func (file TorrentFile) GetDuplicatesTotal() int {
	total := 0
	for _, v := range file.duplicatesMaps {
		total += v
	}
	return total
}

//appendNewSegment adds new raw data to the torrentfile by hashing it and storing it in a map. If the same hash has
//been computed previously, then the counter is incremented for that particular entry (to allow for analyzing if there
//are common layers of file, and thus save on storage/bandwidth)
func (file *TorrentFile) appendNewSegment(segData []byte) {
	hexHashed := hex.EncodeToString(hashSegment(segData))          //hash the data
	file.SegmentHashKeys = append(file.SegmentHashKeys, hexHashed) //record the hash in order

	if _, ok := file.segmentHashMaps[hexHashed]; ok { //we've generated this hash before
		if _, okk := file.duplicatesMaps[hexHashed]; okk { //we've made the entry before
			file.duplicatesMaps[hexHashed]++
		} else {
			file.duplicatesMaps[hexHashed] = 1 //this is the 2nd occurrence (counter starts at 1 for "1st duplicate")
		}
	} else {
		file.segmentHashMaps[hexHashed] = segData
	}
}

func hashSegment(seg []byte) []byte {
	h := sha256.New()
	h.Write(seg)
	return h.Sum(nil)
}
