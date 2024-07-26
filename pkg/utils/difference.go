package utils

func Difference(a, b []uint) []uint {
	mb := make(map[uint]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []uint
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
