package momos

type SSIAttributes map[string]string

type SSIElement struct {
	Tag         string
	Pos         int
	Len         int
	HasErrorTag bool
	ErrorHTML   string
	Attributes  SSIAttributes
}
