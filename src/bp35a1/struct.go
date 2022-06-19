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
