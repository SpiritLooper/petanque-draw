package utils

func LeftRotation(a []int, rotation int) []int {
	size := len(a)
	var newArray []int
	for range rotation {
		newArray = a[1:size]
		newArray = append(newArray, a[0])
		a = newArray
	}
	return a
}
