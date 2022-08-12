package rediscache

import "time"

const (
	Maxidle = 0
	Maxactive = 0
	Idletimeout = time.Second * 3
	Dialreadtimeout = time.Second * 3
	Dialwritetimeout = time.Second * 3
	Dialconnecttimeout = time.Second * 3
)