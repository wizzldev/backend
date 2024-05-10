package utils

func RemoveFromSlice[S []E, E comparable](slice S, elem E) S {
	for i := 0; i < len(slice); i++ {
		if slice[i] == elem {
			// Remove the element by slicing the slice to exclude it
			slice = append(slice[:i], slice[i+1:]...)
			// Decrement i to account for the removed element
			i--
		}
	}
	return slice
}
