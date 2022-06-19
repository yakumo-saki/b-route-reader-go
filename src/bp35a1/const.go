package bp35a1

const RET_OK = "OK"
const RET_SCAN_COMPLETE = "EVENT 22 "
const RET_SCAN_FOUND = "EPANDESC"
const RET_FAIL = "FAIL ER"

var RET_STOP_WORDS = []string{RET_OK, RET_SCAN_COMPLETE, RET_FAIL}

// PANDESC
const RET_PAN_CHANNEL = "Channel:"
const RET_PAN_CHANNEL_PAGE = "Channel Page:"
const RET_PAN_ID = "Pan ID:"
const RET_PAN_ADDR = "Addr:"
const RET_PAN_LQI = "LQI:"
const RET_PAN_PAIR_ID = "PairID:"
