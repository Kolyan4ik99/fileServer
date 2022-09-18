package model

type FileInfo struct {
	Id   int
	Name string
	Size int
}

type TMPFILE []FileInfo

func (d TMPFILE) Len() int {
	return len(d)
}

func (d TMPFILE) Less(i, j int) bool {
	return d[i].Id < d[j].Id
}

func (d TMPFILE) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
