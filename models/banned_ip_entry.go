package models

import "time"

type BannedIPEntry struct {
	IP         string    `json:"ip"`
	BannedTill time.Time `json:"banned_till"`
}
