package logparse

import (
	"goipban/config"
	"testing"
	"time"
)

const blacklistLineSample = "2025/01/13 12:34:56 " +
	"from 12.34.56.78:9999 " +
	"accepted " +
	"tcp:https://ya.ru:443 " +
	"[inbound -> blacklist] " +
	"email: user_mail"

func TestParseXRayLogLine(t *testing.T) {
	// Test cases
	test := struct {
		source   string
		expected XRayLogEntry
	}{
		source: blacklistLineSample,

		expected: XRayLogEntry{
			Time:     time.Date(2025, time.January, 13, 12, 34, 56, 0, time.Local),
			FromIP:   "12.34.56.78",
			FromPort: "9999",
			To:       "tcp:https://ya.ru:443",
			Inbound:  "inbound",
			Outbound: "blacklist",
		},
	}

	t.Run("", func(t *testing.T) {
		result := parseXRayLogLine(test.source)
		if result == nil {
			t.Errorf("parseXRayLogLine(%s) could not find any log entry", test.source)
		} else if *result != test.expected {
			t.Errorf("parseXRayLogLine(%s) = %s; want %s", test.source, *result, test.expected)
		}
	})
}

func TestGetBlacklistedXRayLogEntries(t *testing.T) {
	t.Run("", func(t *testing.T) {
		config, err := config.LoadConfig("")
		if err != nil {
			t.Fatalf("Can not load config %s", err)
		}
		logParser := CreateXRayLogParser(*config)

		result, err := logParser.GetBlacklistedXRayLogEntries()
		if err != nil {
			t.Errorf("GetBlacklistedXRayLogEntries returns error %s", err)
		} else {
			t.Logf("Found %d IPs to ban", len(result))
			for i := range result {
				t.Logf("IP %s to ban from %s", result[i].FromIP, result[i].Time)
			}
		}
	})
}
