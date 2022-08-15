package rediscache

import "time"

const (
	Maxidle = 64
	Maxactive = 64
	Idletimeout = time.Second * 3
	Dialreadtimeout = time.Second * 3
	Dialwritetimeout = time.Second * 3
	Dialconnecttimeout = time.Second * 3
)