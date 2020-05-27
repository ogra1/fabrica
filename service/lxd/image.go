package lxd

// GetImageAlias checks of an image alias is available
func (lx *LXD) GetImageAlias(name string) error {
	// Connect to the lxd service
	c, err := lx.connect()
	if err != nil {
		return err
	}

	// Check if the alias exists (the image could still be loading)
	_, _, err = c.GetImageAlias(name)
	return err
}
