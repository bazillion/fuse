package fuse

func localVolume(conf *mountConfig) error {
	return nil
}

func volumeName(name string) MountOption {
	return dummyOption
}

func daemonTimeout(name string) MountOption {
	return func(conf *mountConfig) error {
		conf.options["timeout"] = name
		return nil
	}
}

func noAppleXattr(conf *mountConfig) MountOption {
	return nil
}

func noAppleDouble(conf *mountConfig) MountOption {
	return nil
}
