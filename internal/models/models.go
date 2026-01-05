package models

// GatewayResponse data structures matching the T-Mobile Gateway JSON
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
	FourG   ConnectionStats `json:"4g"`
	Generic GenericInfo     `json:"generic"`
}

type ConnectionStats struct {
	AntennaUsed string   `json:"antennaUsed"`
	Bands       []string `json:"bands"`
	Bars        float64  `json:"bars"`
	CID         int      `json:"cid"`
	GNBID       int      `json:"gNBID"`
	EID         int      `json:"eid"`
	PCID        int      `json:"pcid"`
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

// PingStats represents the latency statistics.
type PingStats struct {
	Min      float64 `json:"min"`
	Avg      float64 `json:"avg"`
	Max      float64 `json:"max"`
	StdDev   float64 `json:"stddev"`
	Loss     float64 `json:"loss"`
	LastRTT  float64 `json:"last_rtt"`
	Sent     int     `json:"sent"`
	Received int     `json:"received"`
}

// CombinedStats represents the full set of monitored data.
type CombinedStats struct {
	Gateway GatewayResponse `json:"gateway"`
	Ping    PingStats       `json:"ping"`
}