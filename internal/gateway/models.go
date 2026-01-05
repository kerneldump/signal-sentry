package gateway

// Data structures matching the T-Mobile Gateway JSON

type GatewayResponse struct {
	Device DeviceInfo `json:"device"`
	Signal SignalInfo `json:"signal"`
	Time   TimeInfo   `json:"time"`
}

type DeviceInfo struct {
	HardwareVersion string `json:"hardwareVersion"`
	MacID           string `json:"macId"`
	Manufacturer    string `json:"manufacturer"`
	Model           string `json:"model"`
	Role            string `json:"role"`
	Serial          string `json:"serial"`
	SoftwareVersion string `json:"softwareVersion"`
}

type SignalInfo struct {
	FiveG   ConnectionStats `json:"5g"`
	FourG   ConnectionStats `json:"4g"` // Assuming 4G structure is similar if present
	Generic GenericInfo     `json:"generic"`
}

type ConnectionStats struct {
	AntennaUsed string   `json:"antennaUsed"`
	Bands       []string `json:"bands"`
	Bars        float64  `json:"bars"`
	CID         int      `json:"cid"`
	GNBID       int      `json:"gNBID"` // 5G uses gNBID
	EID         int      `json:"eid"`   // 4G often uses eid or similar, but structure usually shares fields. We'll stick to common ones.
	PCID        int      `json:"pcid"`  // 4G Physical Cell ID
	RSRP        int      `json:"rsrp"`
	RSRQ        int      `json:"rsrq"`
	RSSI        int      `json:"rssi"`
	SINR        int      `json:"sinr"`
}

type GenericInfo struct {
	APN          string `json:"apn"`
	HasIPv6      bool   `json:"hasIPv6"`
	Registration string `json:"registration"`
}

type TimeInfo struct {
	LocalTime     int64  `json:"localTime"`
	LocalTimeZone string `json:"localTimeZone"`
	UpTime        int    `json:"upTime"`
}
