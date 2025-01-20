package logparse

import (
	"bufio"
	"fmt"
	"goipban/config"
	"net"
	"os"
	"time"
)

type XRayLogEntry struct {
	Time     time.Time `json:"time"`
	FromIP   string    `json:"from_ip"`
	FromPort string    `json:"from_port"`
	To       string    `json:"to"`
	Inbound  string    `json:"inbound"`
	Outbound string    `json:"outbound"`
}

// tool to parse xray logs
type XRayLogParser struct {
	config config.Config
}

// NewBanDB creates or opens the SQLite database and initializes the table
func CreateXRayLogParser(config config.Config) XRayLogParser {
	return XRayLogParser{config: config}
}

// log file already parsed from beginning to this position
func (p *XRayLogParser) GetBlacklistedXRayLogEntries() ([]XRayLogEntry, error) {
	// list of banned ip entries
	var blacklistedIPEntries []XRayLogEntry

	// Open the logs file in read-only mode
	logFile, err := os.Open(p.config.XRayLogsFile)
	if err != nil {
		return blacklistedIPEntries, fmt.Errorf("error opening file: %v", err)
	}
	defer logFile.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(logFile)

	// collect ip's last mentions as addresses to be banned
	var lastBaningIPsOccurences = make(map[string]XRayLogEntry)
	for scanner.Scan() {
		// Get the current line
		line := scanner.Text()

		// Parse current line
		if parsedLineEntry := parseXRayLogLine(line); parsedLineEntry != nil {
			// check if line is about ban
			if parsedLineEntry.Outbound == p.config.XRayBlacklistOutbound {
				lastBaningIPsOccurences[parsedLineEntry.FromIP] = *parsedLineEntry
			}
		}
	}

	// select IPs which was banned not too much time ago
	earliestBannableTime := time.Now().Add(-p.config.BanDuration.TimeDuration())
	for _, logEntry := range lastBaningIPsOccurences {
		if logEntry.Time.After(earliestBannableTime) {
			blacklistedIPEntries = append(blacklistedIPEntries, logEntry)
		}
	}

	return blacklistedIPEntries, nil
}

/*
BlacklistLineSample: `2025/13/01 12:34:56 // 0, 1

	from                                  // 2
	12.34.56.78:9999                      // 3
	accepted tcp:https://ya.ru:443        // 4, 5
	[inbound -> blacklist]                // 6, 7, 8
	email: user_mail`,                    // 9
*/
func parseXRayLogLine(logLine string) *XRayLogEntry {
	spacePositions := findCharPositions(logLine, ' ')

	if len(spacePositions) != 10 {
		return nil
	}

	// parse time
	const timeLayout = "2006/01/02 15:04:05"
	time, err := time.ParseInLocation(timeLayout, logLine[0:spacePositions[1]], time.Local)
	if err != nil {
		return nil
	}

	// parse ip, port
	ip, port, err := net.SplitHostPort(logLine[spacePositions[2]+1 : spacePositions[3]])
	if err != nil {
		return nil
	}

	// parse to, inbound, outbound
	toAddress := logLine[spacePositions[4]+1 : spacePositions[5]]
	inbound := logLine[spacePositions[5]+2 : spacePositions[6]]
	outbound := logLine[spacePositions[7]+1 : spacePositions[8]-1]

	return &XRayLogEntry{
		Time:     time,
		FromIP:   ip,
		FromPort: port,
		To:       toAddress,
		Inbound:  inbound,
		Outbound: outbound,
	}
}

func findCharPositions(s string, c rune) []int {
	var positions []int
	for i, char := range s {
		if char == c {
			positions = append(positions, i)
		}
	}
	return positions
}
