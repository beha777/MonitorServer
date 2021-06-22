package jobs

import (
	"github.com/reiver/go-telnet"
	"os/exec"
	"strings"
)

func CheckTelnet(hostPort string) error {
	_, err := telnet.DialTo(hostPort)
	return err
}

func CheckPing(hostIp string) error {
	out, err := exec.Command("ping", hostIp).Output()
	if err != nil && !strings.Contains(string(out), "Lost = 0") {
		return err
	}
	return nil
}
