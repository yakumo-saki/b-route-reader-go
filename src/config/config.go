package config

import "os"

var B_ROUTE_ID = ""
var B_ROUTE_PASSWORD = ""
var SERIAL = "/home/yakumo/ttyAMA0"

var ACTIVE_SCAN_COUNT = 10

// 環境変数からconfigをセット
func Initialize() {
	B_ROUTE_ID = os.Getenv("B_ROUTE_ID")
	B_ROUTE_PASSWORD = os.Getenv("B_ROUTE_PASSWORD")
	SERIAL = os.Getenv("SERIAL")
}
