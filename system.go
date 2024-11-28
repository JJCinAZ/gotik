package gotik

import (
	"fmt"
	"time"
)

type Routerboard struct {
	Routerboard     bool   `json:"routerboard"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serial-number"`
	FirmwareType    string `json:"firmware-type"`
	FactoryFirmware string `json:"factory-firmware"`
	CurrentFirmware string `json:"current-firmware"`
	UpgradeFirmware string `json:"upgrade-firmware"`
}

type Resources struct {
	Uptime               time.Duration `json:"uptime"`
	Version              string        `json:"version"`
	BuildTime            time.Time     `json:"build-time"`
	FactorySoftware      string        `json:"factory-software"`
	FreeMemory           int           `json:"free-memory"`
	TotalMemory          int           `json:"total-memory"`
	CPU                  string        `json:"cpu"`
	CPUCount             int           `json:"cpu-count"`
	CPUFrequency         int           `json:"cpu-frequency"` // in MHz
	CPULoad              int           `json:"cpu-load"`
	FreeHddSpace         int           `json:"free-hdd-space"`
	TotalHddSpace        int           `json:"total-hdd-space"`
	WriteSectSinceReboot int           `json:"write-sect-since-reboot"`
	WriteSectTotal       int           `json:"write-sect-total"`
	BadBlocks            int           `json:"bad-blocks"`
	ArchitectureName     string        `json:"architecture-name"`
	BoardName            string        `json:"board-name"`
	Platform             string        `json:"platform"`
}

type License struct {
	SoftwareId string `json:"software-id"`
	Level      int    `json:"nlevel"`
	Features   string `json:"features"`
}

func parseResources(props map[string]string) Resources {
	entry := Resources{
		Uptime:           parseDuration(props["uptime"]),
		Version:          props["version"],
		BuildTime:        parseTime(props["build-time"]),
		FactorySoftware:  props["factory-software"],
		FreeMemory:       parseInt(props["free-memory"]),
		TotalMemory:      parseInt(props["total-memory"]),
		CPU:              props["cpu"],
		CPUCount:         parseInt(props["cpu-count"]),
		CPUFrequency:     parseInt(props["cpu-frequency"]),
		CPULoad:          parseInt(props["cpu-load"]),
		FreeHddSpace:     parseInt(props["free-hdd-space"]),
		TotalHddSpace:    parseInt(props["total-hdd-space"]),
		ArchitectureName: props["architecture-name"],
		BoardName:        props["board-name"],
		Platform:         props["platform"],
	}
	return entry
}

func (c *Client) GetSystemResources() (Resources, error) {
	detail, err := c.RunCmd("/system/resource/print")
	if err == nil {
		r := parseResources(detail.Re[0].Map)
		return r, nil
	}
	return Resources{}, err
}

func parseRouterboard(props map[string]string) Routerboard {
	return Routerboard{
		Routerboard:     parseBool(props["routerboard"]),
		Model:           props["model"],
		SerialNumber:    props["serial-number"],
		FirmwareType:    props["firmware-type"],
		FactoryFirmware: props["factory-firmware"],
		CurrentFirmware: props["current-firmware"],
		UpgradeFirmware: props["upgrade-firmware"],
	}
}

func (c *Client) GetSystemRouterboard() (Routerboard, error) {
	detail, err := c.RunCmd("/system/routerboard/print")
	if err == nil {
		r := parseRouterboard(detail.Re[0].Map)
		return r, nil
	}
	return Routerboard{}, err
}

func (c *Client) GetSystemId() (string, error) {
	detail, err := c.RunCmd("/system/identity/print")
	if err == nil {
		name := detail.Re[0].Map["name"]
		return name, nil
	}
	return "", err
}

func parseLicense(props map[string]string) License {
	entry := License{
		SoftwareId: props["softwareid"],
		Level:      parseInt(props["nlevel"]),
		Features:   props["features"],
	}
	return entry
}

func (c *Client) GetSystemLicense() (License, error) {
	detail, err := c.RunCmd("/system/license/print")
	if err == nil {
		r := parseLicense(detail.Re[0].Map)
		return r, nil
	}
	return License{}, err
}

func (c *Client) CreateExport(targetname string, minFreeSpace int) error {
	var (
		found bool
	)
	if files, err := c.GetAllFiles(); err != nil {
		return err
	} else {
		for _, f := range files {
			if f.Name == targetname {
				found = true
				break
			}
		}
	}
	if !found {
		if r, err := c.GetSystemResources(); err != nil {
			return err
		} else if minFreeSpace != 0 && r.FreeHddSpace < minFreeSpace {
			return fmt.Errorf("not enough free space (< 256KB)")
		}
		if err := c.ExportConfig("/", targetname, false); err != nil {
			return err
		}
	}
	return nil
}
