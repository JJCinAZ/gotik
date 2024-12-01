package gotik

import (
	"fmt"
)

type Package struct {
	Name      string `json:"name" tik:"name"`
	Disabled  bool   `json:"disabled" tik:"disabled"`
	Version   string `json:"version" tik:"version"`
	BuildTime string `json:"build-time" tik:"build-time"`
	Scheduled string `json:"scheduled" tik:"scheduled"`
	Bundle    string `json:"bundle" tik:"bundle"`
}

func parsePackage(props map[string]string) Package {
	entry := Package{
		Name:      props["name"],
		Disabled:  parseBool(props["disabled"]),
		Version:   props["version"],
		BuildTime: props["build-time"],
		Scheduled: props["scheduled"],
		Bundle:    props["bundle"],
	}
	return entry
}

func (c *Client) GetPackages() ([]Package, error) {
	entries := make([]Package, 0)
	detail, err := c.Run("/system/package/print")
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parsePackage(detail.Re[i].Map))
		}
	} else {
		entries = nil
	}
	return entries, err
}

func parsePackageUpdate(props map[string]string) PackageUpdate {
	entry := PackageUpdate{
		Channel:   props["channel"],
		Installed: props["installed-version"],
		Latest:    props["latest-version"],
		Status:    props["status"],
	}
	return entry
}

// GetUpdateInfo will query the router for the current channel, installed and latest versions and the current status
// of the updates ("Downloaded, please reboot router to upgrade it","","System is already up to date", "finding out latest version...")
func (c *Client) GetUpdateInfo() (info PackageUpdate, err error) {
	var detail *Reply
	detail, err = c.Run("/system/package/update/print")
	if err != nil {
		return
	}
	for i := range detail.Re {
		info = parsePackageUpdate(detail.Re[i].Map)
		return
	}
	err = fmt.Errorf("invalid return")
	return
}

// CheckForUpdates will command the router to check with its package sources for the latest update on the current
// update channel (set with SetUpdateChannel).  The function returns before the check completes, so the value
// of info.Status may be "finding out latest version..."
func (c *Client) CheckForUpdates() (info PackageUpdate, err error) {
	var detail *Reply
	detail, err = c.Run("/system/package/update/check-for-updates")
	if err != nil {
		return
	}
	for i := range detail.Re {
		info = parsePackageUpdate(detail.Re[i].Map)
		return
	}
	err = fmt.Errorf("invalid return")
	return
}

func (c *Client) SetUpdateChannel(channel string) error {
	_, err := c.Run("/system/package/update/set", fmt.Sprintf("=channel=%s", channel))
	if err != nil {
		return err
	}
	return nil
}

// DownloadUpdates will initiate a download of any updates on the current "channel"
// The function will not return until the download is complete or fails so this function may take a long time
// to execute based on the speed at which the target device can download from mikrotik.com.
// The final status message will be returned in the info.Status variable.
func (c *Client) DownloadUpdates() (info PackageUpdate, err error) {
	var (
		detail *Reply
	)
	/*
	   The return from this Run looks like:
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `calculating download size...`} {`.section` `0`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `downloading...`} {`.section` `1`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 8% (0.9MiB)`} {`.section` `2`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 17% (2.0MiB)`} {`.section` `3`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 28% (3.1MiB)`} {`.section` `4`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 38% (4.3MiB)`} {`.section` `5`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 48% (5.4MiB)`} {`.section` `6`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 58% (6.5MiB)`} {`.section` `7`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 68% (7.7MiB)`} {`.section` `8`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 77% (8.7MiB)`} {`.section` `9`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 87% (9.8MiB)`} {`.section` `10`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded 98% (11.0MiB)`} {`.section` `11`}]
	      !re @ [{`channel` `long-term`} {`installed-version` `6.43.8`} {`latest-version` `6.46.8`} {`status` `Downloaded, please reboot router to upgrade it`} {`.section` `12`}]
	      !done @ []
	*/
	detail, err = c.Run("/system/package/update/download")
	if err != nil {
		return
	}
	if len(detail.Re) > 0 {
		info = parsePackageUpdate(detail.Re[len(detail.Re)-1].Map)
		return
	}
	err = fmt.Errorf("invalid return")
	return
}

// InstallUpdates will initiate a download of any updates on the current "channel" and will reboot the
// router if download is successful.
// The function will not return until the download is complete or fails so this function may take a long time
// to execute based on the speed at which the target device can download from mikrotik.com.
// The final status message will be returned in the info.Status variable.  For example:
//
//	routeros.PackageUpdate{Channel:"long-term", Installed:"6.43.8", Latest:"6.46.8", Status:"Downloaded, rebooting..."}
//
// Note that while there is still an open handle to the router, the router rebooted and further connection
// is no longer possible.  The router handle should be closed immediately after calling this function
// and getting a successful return.
func (c *Client) InstallUpdates() (info PackageUpdate, err error) {
	var (
		detail *Reply
	)
	detail, err = c.Run("/system/package/update/install")
	if err != nil {
		return
	}
	if len(detail.Re) > 0 {
		info = parsePackageUpdate(detail.Re[len(detail.Re)-1].Map)
		return
	}
	err = fmt.Errorf("invalid return")
	return
}
