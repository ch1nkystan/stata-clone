package server

func f64n(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case float64:
		return n
	}

	return 0
}

func div(a, b interface{}) float64 {
	af := f64n(a)
	bf := f64n(b)

	if bf == 0 {
		return 0
	}

	return af / bf
}
