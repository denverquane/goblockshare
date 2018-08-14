package common

type TorrentRep struct {
	TotalReports 	uint64
	ValidReports    uint64
	QualityReports  uint64
	AccurateReports uint64
}

type LayerRep struct {
	//How many times this layer has been shared w/ someone else on the blockchain
	SharedQuantity uint64

	NotReceived uint64

	//How many times this layer has been reported as valid upon receipt
	ValidReports uint64
}

type ReputationSummary struct {
	//Reputation of all torrents, indexed by their hashes
	TorrentRep map[string]TorrentRep

	//Reputation of all layers, indexed by hash
	LayerRep map[string]LayerRep
}

type JSONRepSummary struct {
	ValidTorrFraction float64
	QualityTorrFraction float64
	AccurateTorrFraction float64

	NotReceivedLayerFraction float64
	ValidLayerFraction float64
}

func (summary ReputationSummary) ToJSONSummary() JSONRepSummary {
	json := JSONRepSummary{0, 0, 0, 0, 0}
	var total = 0.0
	for _, v := range summary.TorrentRep {
		total += float64(v.TotalReports)
		json.ValidTorrFraction += float64(v.ValidReports)
		json.AccurateTorrFraction += float64(v.AccurateReports)
		json.QualityTorrFraction += float64(v.QualityReports)
	}
	if total != 0.0 {
		json.ValidTorrFraction = json.ValidTorrFraction / total
		json.QualityTorrFraction = json.QualityTorrFraction / total
		json.AccurateTorrFraction = json.AccurateTorrFraction / total
	}

	total = 0
	for _, v := range summary.LayerRep {
		total += float64(v.SharedQuantity)
		json.NotReceivedLayerFraction += float64(v.NotReceived)
		json.ValidLayerFraction += float64(v.ValidReports)
	}
	if total != 0.0 {
		json.ValidLayerFraction = json.ValidLayerFraction / total
		json.NotReceivedLayerFraction = json.NotReceivedLayerFraction / total
	} else {
		json.NotReceivedLayerFraction = 1.0
	}

	return json
}
