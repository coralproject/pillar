package aggregate

type Int map[string]int

func NewInt() Int {
	return make(map[string]int)
}

func (a Int) Add(key string, delta int) {
	a[key] += delta
}

func (a Int) Total(key string) int {
	return a[key]
}

func (a Int) Ratio(numerator, denominator string) float64 {
	if a[denominator] == 0 {
		return .0
	}
	return float64(a[numerator]) / float64(a[denominator])
}

func (a Int) Keys() []string {
	keys := make([]string, 0, len(a))
	for key, _ := range a {
		keys = append(keys, key)
	}
	return keys
}

func (a Int) Combine(b Int) {
	for key, value := range b {
		a[key] += value
	}
}
