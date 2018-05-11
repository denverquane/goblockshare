package files

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMakeTorrentFile(t *testing.T) {
	err, torr := MakeTorrentFileFromFile(kilobyte, "../../../../../../Downloads/CentOS-7-x86_64-Minimal-1708.iso")
	if err != nil {
		fmt.Print(err)
		return
	}
	var total int64 = 0
	for i, v := range torr.segmentHashMap {
		if v > 1 {
			fmt.Println("Duplicate " + i + "occurs " + strconv.FormatInt(v, 10) + "times")
			total += v
		}
	}
	fmt.Println(strconv.FormatInt(total, 10) + " kilobyte reduction ")
}
