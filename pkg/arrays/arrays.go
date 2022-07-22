package arrays

func Map[T any, M any](ary []T, f func(T, int) M) []M {
	n := make([]M, len(ary))
	for i, item := range ary {
		n[i] = f(item, i)
	}
	return n
}

func Filter[T any](ary []T, f func(T, int) bool) []T {
	n := make([]T, 0)
	for i, item := range ary {
		if f(item, i) {
			n[i] = item
		}
	}
	return n
}

func ForEach[T any](ary []T, f func(T, int)) {
	for i, item := range ary {
		f(item, i)
	}
}
