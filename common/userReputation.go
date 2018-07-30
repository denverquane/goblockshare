package common

type TorrentRep struct {
	ValidReports	uint64
	QualityReports  uint64
	AccurateReports uint64
}

type LayerRep struct {
	//How many times this layer has been shared w/ someone else on the blockchain
	SharedQuantity uint64

	NotReceived	   uint64

	//How many times this layer has been reported as valid upon receipt
	ValidReports   uint64
}

type ReputationSummary struct {
	//Reputation of all torrents, indexed by their hashes
	TorrentRep	 map[string]TorrentRep

	//Reputation of all layers, indexed by hash
	LayerRep		 map[string]LayerRep
}
