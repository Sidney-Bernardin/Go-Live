package src

func MapSlice[A, B any](aa []A, cb func(A) B) []B {
	bb := make([]B, len(aa))
	for i, a := range aa {
		bb[i] = cb(a)
	}
	return bb
}
