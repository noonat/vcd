package vcd

import (
	"net"
	"time"
)

// VesselPilotClick tracks when someone clicks to pilot a vessel.
type VesselPilotClick struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IP        net.IP    `json:"ip"`
	Referrer  string    `json:"referrer"`
}
