package types

type FindTxResponse struct {
	Mode   string   `json:"mode"`
	Hashes []string `json:"hashes"`
	Cursor struct {
		Node    int  `json:"node"`
		HasMore bool `json:"hasMore"`
		Perma   int  `json:"perma"`
		NextInt int  `json:"nextInt"`
	} `json:"cursor"`
}

type TxTrytesResponse struct {
	MilestoneIndexes []int64  `json:"milestoneIndexes"`
	Trytes           []string `json:"trytes"`
}

type StateInfo struct {
	Accountname string
	Addresses   []AddrInfo
}

type AddrInfo struct {
	Address      string
	TxHashes     []string
	TxInfos      []TxInfo
	BundleHashes []string
	Balance      uint64
}

type TxInfo struct {
	Hash                string
	Address			string
	Value               int64
	AttachmentTimestamp int64
	Timestamp           uint64
	ConfMilestone       int64
	Confirmed           bool
	Tag                 string
	ObsoleteTag         string
	MessageASCII        string
	Bundlehash          string
	rawTrytes           string
}

type ValueBundle struct {
	Hash     string
	Date     string
	Addresses  []string
	Outgoing bool
	Internal bool
	Value    int64
	TxInfos  []TxInfo
}
