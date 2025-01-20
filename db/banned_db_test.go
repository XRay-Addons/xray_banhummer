package db

import (
	"goipban/models"
	"os"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	dbFilename := "data/testing_db.sqlite"
	defer os.Remove(dbFilename)

	db, err := ConnectToBannedDB(dbFilename)
	if err != nil {
		t.Errorf("err connect database %v", err)
	}
	defer db.Close()
}

func TestInsertBannedIPs(t *testing.T) {
	dbFilename := "data/testing_db.sqlite"
	defer os.Remove(dbFilename)

	db, err := ConnectToBannedDB(dbFilename)
	if err != nil {
		t.Errorf("err connect database %v", err)
	}
	defer db.Close()

	bannedIP := models.BannedIPEntry{
		IP:         "1.2.3.4",
		BannedTill: time.Now(),
	}
	err = db.InsertBannedIP(bannedIP)
	if err != nil {
		t.Errorf("err inserting database %v", err)
		return
	}

	bannedIPs, err := db.GetExpiredEntries(time.Now())
	if err != nil {
		t.Errorf("err get values from database %v", err)
		return
	}
	if len(bannedIPs) != 1 {
		t.Errorf("get errors from banned before returns %d items, expected 1", len(bannedIPs))
	}

	err = db.RemoveBannedIP(bannedIP.IP)
	if err != nil {
		t.Errorf("err remove banned IP from database %v", err)
		return
	}
	bannedIPs, err = db.GetExpiredEntries(time.Now())
	if err != nil {
		t.Errorf("err get values from database %v", err)
		return
	}
	if len(bannedIPs) != 0 {
		t.Errorf("get errors from banned before returns %d items, expected 0", len(bannedIPs))
	}
}
