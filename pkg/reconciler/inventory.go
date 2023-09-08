package reconciler

type inventory map[string]bool

func (i *inventory) count() (ready int, total int) {
	if i == nil {
		return 0, 0
	}

	for _, item := range *i {
		total++
		if item {
			ready++
		}
	}
	return
}

func (i *inventory) ready() bool {
	if i == nil {
		return true
	}

	for _, item := range *i {
		if !item {
			return false
		}
	}
	return true
}
