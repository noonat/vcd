package vcd

import (
	"net"
	"time"
)

// VesselClick tracks when someone clicks on a vessel.
type VesselClick struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IP        net.IP    `json:"ip"`
	Referrer  string    `json:"referrer"`
}
