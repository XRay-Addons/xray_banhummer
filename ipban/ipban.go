package ipban

import (
	"fmt"
	"goipban/config"
	"log"
	"os/exec"
	"strings"
)

var commands = struct {
	Ban   string
	Unban string
	Apply string
}{
	Ban:   "ufw insert %d deny from %s",
	Unban: "ufw delete deny from %s",
	Apply: "ufw reload",
}

// tool to ban/unban IPs
type IPBanner struct {
	config config.Config
}

// NewBanDB creates or opens the SQLite database and initializes the table
func CreateIPBanner(config config.Config) IPBanner {
	return IPBanner{config: config}
}

// Ban the IP using UFW
func (b *IPBanner) BanIP(ip string) error {
	log.Printf("Banning IP: %s\n", ip)
	banRecordIdx := b.config.ProtectedUFWRulesCount + 1
	if err := b.runCommand(fmt.Sprintf(commands.Ban, banRecordIdx, ip)); err != nil {
		log.Printf("Error banning IP %s: %v\n", ip, err)
		return err
	}
	if err := b.runCommand(commands.Apply); err != nil {
		log.Printf("Error applying ban: %v\n", err)
		return err
	}
	return nil
}

// Unan the IP using UFW
func (b *IPBanner) UnbanIP(ip string) error {
	log.Printf("Unbanning IP: %s\n", ip)
	if err := b.runCommand(fmt.Sprintf(commands.Unban, ip)); err != nil {
		log.Printf("Error unbanning IP %s: %v\n", ip, err)
		return err
	}
	if err := b.runCommand(commands.Apply); err != nil {
		log.Printf("Error applying unban: %v\n", err)
		return err
	}
	return nil
}

// runCommand executes shell commands
func (b *IPBanner) runCommand(cmd string) error {
	log.Printf("Executing command: %s\n", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	args := parts[1:]

	if b.debugMode {
		return nil // don't ban, just print
	}

	out, err := exec.Command(head, args...).CombinedOutput()
	if err != nil {
		log.Printf("Command output: %s\n", string(out))
		return err
	}
	return nil
}
