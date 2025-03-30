package cache

import "github.com/fluffy-melli/RouteNX/pkg/config"

var Config *config.RouteNX

var TX = int64(0)
var RX = int64(0)

var Label = make([]int64, 0)
var TXBPS = make([]int64, 0)
var RXBPS = make([]int64, 0)
