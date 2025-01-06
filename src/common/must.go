package common

func Must(err error) {
	Checke(err)
}

func Must2[V any](value V, err error) V {
	Checke(err)
	return value
}

func Must3[V1 any, V2 any](value1 V1, value2 V2, err error) (V1, V2) {
	Checke(err)
	return value1, value2
}
