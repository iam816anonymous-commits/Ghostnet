package common

import (
	"time"
)

type TrafficClass string

const (
	ClassSmallWeb    TrafficClass = "Small Web Requests"
	ClassAPI         TrafficClass = "API Requests"
	ClassImage       TrafficClass = "Image Downloads"
	ClassVideo       TrafficClass = "Video Streaming"
	ClassFileUpload  TrafficClass = "File Uploads"
	ClassDNS         TrafficClass = "DNS Requests"
	ClassBackground  TrafficClass = "Background Sync"
	ClassInteractive TrafficClass = "Interactive Sessions"
	ClassUnknown     TrafficClass = "Unknown"
)

type PrivacyMode string

const (
	ModeOff            PrivacyMode = "OFF"
	ModeLow            PrivacyMode = "LOW"
	ModeMedium         PrivacyMode = "MEDIUM"
	ModeHigh           PrivacyMode = "HIGH"
	ModeMaximumPrivacy PrivacyMode = "MAXIMUM PRIVACY"
)

type Packet struct {
	ID        string
	Payload   []byte
	Size      int
	Type      TrafficClass
	IsDummy   bool
	Metadata  Metadata
	Timestamp time.Time
}

type Metadata struct {
	Source      string
	Destination string
	IdentityID  string
	Route       []string
	HopIndex    int
}
