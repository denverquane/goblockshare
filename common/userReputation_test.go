package common

import "testing"

func TestReputationSummary_ToJSONSummary(t *testing.T) {
	torrMap := make(map[string]TorrentRep, 10)
	layerMap := make(map[string]LayerRep, 10)

	var totalReports = 10.0

	for i := 0; i < int(totalReports); i++ {
		torrMap[string(i)] = TorrentRep{5, 4, 3, 2}
		layerMap[string(i)] = LayerRep{5, 4, 3}
	}

	repSum := ReputationSummary{torrMap, layerMap}
	json := repSum.ToJSONSummary()

	if json.NotReceivedLayerFraction != (4.0 * totalReports)/(5.0 * totalReports) {
		t.Fail()
	}
	if json.ValidLayerFraction != (3.0 * totalReports)/(5.0 * totalReports) {
		t.Fail()
	}

	if json.ValidTorrFraction != (4.0 * totalReports)/(5.0 * totalReports) {
		t.Fail()
	}
	if json.QualityTorrFraction != (3.0 * totalReports)/(5.0 * totalReports) {
		t.Fail()
	}
	if json.AccurateTorrFraction != (2.0 * totalReports)/(5.0 * totalReports) {
		t.Fail()
	}

}