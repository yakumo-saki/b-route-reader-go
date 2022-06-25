package bp35a1

import "fmt"

type SmartMeter struct {
	Channel     string
	ChannelPage string
	PanId       string
	Addr        string
	LQI         string
	PairId      string
}

func (sm SmartMeter) String() string {
	return fmt.Sprintf("Channel: %s (page: %s)\nPanID: %s\nAddr: %s\nPairID:%s",
		sm.Channel, sm.ChannelPage, sm.PanId, sm.Addr, sm.PairId)
}

type ElectricData struct {
	RphaseAmp float64 // R相電流アンペア Signed Short (0.1A単位)
	TphaseAmp float64 // T相電流アンペア Signed Short (0.1A単位)
	TotalAmp  float64 // 合計電流アンペア Signed Short (0.1A単位)
	DeltakWh  float64 // 積算電力量
	Watt      int     // 瞬間電力量
}

func (ed ElectricData) String() string {
	return fmt.Sprintf("ElectricData: Amps Total=%.1f(R=%.1f T=%.1f) Watt=%d Delta=%.2f",
		ed.TotalAmp, ed.RphaseAmp, ed.TphaseAmp, ed.Watt, ed.DeltakWh)
}
