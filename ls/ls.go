package ls

type Flags struct {
	L bool
}

type LS struct {
	Dir string
	Flags
}

func (l *LS) ListDir() error {

	return nil
}
