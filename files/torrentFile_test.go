package files

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMakeTorrentFile(t *testing.T) {
	now := time.Now()
	torr, err := MakeTorrentFileFromFile(kilobyte, "../../../../../../Downloads/CentOS-7-x86_64-Minimal-1708.iso")
	after := time.Now()
	fmt.Println("Took " + strconv.FormatFloat(after.Sub(now).Seconds(), 'f', 2, 64) + " seconds to make torrent")

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(torr.ValidateHashes())
	//var total int64 = 0
	//for _, v := range torr.segmentHashMap {
	//	if v > 1 {
	//		//fmt.Println("Duplicate " + i + "occurs " + strconv.FormatInt(v, 10) + "times")
	//		total += v
	//	}
	//}
	//var ratio = float64(total) / float64(len(torr.segmentHashes))
	//fmt.Println("Can reduce size by " + strconv.FormatFloat(ratio * 100.0, 'f', 5, 64) + "%")
}

func TestAreSameTorrentBytes(t *testing.T) {
	sampleA := []byte("asdfasdf")
	sampleB := []byte("asdfasdf")

	if !AreSameTorrentBytes(2, sampleA, sampleB) {
		t.Fail()
	}
	sampleC := []byte("asdfasd")

	if AreSameTorrentBytes(2, sampleA, sampleC) {
		t.Fail()
	}

	sampleD := []byte("asdfasdg")

	if AreSameTorrentBytes(10, sampleA, sampleD) {
		t.Fail()
	}
}
