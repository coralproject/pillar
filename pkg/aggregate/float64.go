package aggregate

type Float64 map[string]float64

func NewFloat64() Float64 {
	return make(map[string]float64)
}

func (a Float64) Add(key string, delta float64) {
	a[key] += delta
}

func (a Float64) Total(key string) float64 {
	return a[key]
}

func (a Float64) Ratio(numerator, denominator string) float64 {
	if a[denominator] == 0 {
		return .0
	}
	return a[numerator] / a[denominator]
}

func (a Float64) Keys() []string {
	keys := make([]string, 0, len(a))
	for key, _ := range a {
		keys = append(keys, key)
	}
	return keys
}

func (a Float64) Combine(b Float64) {
	for key, value := range b {
		a[key] += value
	}
}
