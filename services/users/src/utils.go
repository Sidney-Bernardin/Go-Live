package src

type mapFunc[A, B any] func(A) B

func MapSlice[A, B any](aa []A, cb mapFunc[A, B]) []B {
	bb := make([]B, len(aa))
	for i, a := range aa {
		bb[i] = cb(a)
	}
	return bb
}
