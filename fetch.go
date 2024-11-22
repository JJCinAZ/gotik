package gotik

// Fetch will cause the router to download a file from an
func (c *Client) Fetch(filename string, hideSensitive bool) error {
	a := []string{"/tool/fetch", "=file=" + filename, "=compact="}
	if hideSensitive {
		a = append(a, "=hide-sensitive=yes")
	} else {
		a = append(a, "=hide-sensitive=no")
	}
	_, err := c.Run(a...)
	return err
}
