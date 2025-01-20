package ipban

import (
	"goipban/config"
	"testing"
)

func TestBanIP(t *testing.T) {
	config, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Can not load config %s", err)
	}
	ipBanner := CreateIPBanner(*config)
	if err := ipBanner.BanIP("11.22.33.44"); err != nil {
		t.Errorf("Banning error: %v", err)
	} else {
		t.Logf("Banned successfully!")
	}
}

func TestUnbanIP(t *testing.T) {
	config, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Can not load config %s", err)
	}
	ipBanner := CreateIPBanner(*config)
	if err := ipBanner.UnbanIP("11.22.33.44"); err != nil {
		t.Errorf("Unbanning error: %v", err)
	} else {
		t.Logf("Unbanned successfully!")
	}
}
