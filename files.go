package gotik

import (
	"time"
)

type File struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	FileType       string    `json:"type"`
	Size           int       `json:"size"`
	CreationTime   time.Time `json:"creation-time"`
	PackageBldTime time.Time `json:"package-build-time"`
	PackageName    string    `json:"package-name"`
	PackageVersion string    `json:"package-version"`
	PackageArch    string    `json:"package-architecture"`
}

func parseFile(props map[string]string) File {
	entry := File{
		ID:             props[".id"],
		Name:           props["name"],
		FileType:       props["type"],
		Size:           parseInt(props["size"]),
		CreationTime:   parseTime(props["creation-time"]),
		PackageBldTime: parseTime(props["package-build-time"]),
		PackageName:    props["package-name"],
		PackageVersion: props["package-version"],
		PackageArch:    props["package-architecture"],
	}
	return entry
}

func (c *Client) GetAllFiles() ([]File, error) {
	list := make([]File, 0, 1024)
	detail, err := c.RunCmd("/file/print")
	if err != nil {
		return list, err
	}
	for _, re := range detail.Re {
		list = append(list, parseFile(re.Map))
	}
	return list, nil
}

func (c *Client) AddFile(name, contents string) error {
	if c.majorVersion > 7 || (c.majorVersion == 7 && c.minorVersion >= 9) {
		_, err := c.Run("/file/add", "=name="+name, "type=file", "=contents="+contents)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrVersionTooOld
}

func (c *Client) RemoveFileByName(name string) error {
	_, err := c.Run("/file/remove", "=name="+name)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveFileByID(id string) error {
	_, err := c.Run("/file/remove", "=.id="+id)
	if err != nil {
		return err
	}
	return nil
}
