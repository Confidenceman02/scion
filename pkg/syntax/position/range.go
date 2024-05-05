package position

type Location struct {
	Offset int // offset, starting at 0
	Row    int // line number, stqarting at 1
	Column int // column number, starting at 1
}

type Range struct {
	start Location
	end   Location
}
