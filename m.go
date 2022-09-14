package tl

type M map[string]any

// filter used when we want to add or replace a nested translation
// data is old translation which we want to replace or add
// and (m M) is new nested translation
func (m M) filter(data any, override bool) any {
	// pl hold filtered translation
	var pl any

	// check data type
	switch data.(type) {
	case string:
		if override {
			pl = m
		} else {
			pl = data
		}
	case map[string]any:
		// if data type of map[string]any it means translation which we want to
		// replace or add is nested, so we check each old translation first
		tls := data.(map[string]any)

		// m hold new nested translation
		for tag, newTl := range m {
			// check old translation, if old translation doesn't have the translation
			// we want to add, then just add
			if _, exists := tls[tag]; !exists {
				tls[tag] = newTl
				continue
			}

			// if translation already exists but override is true
			// then replace old nested translation with new one
			if override {
				tls[tag] = newTl
			}
		}

		// assign filtered translation to pl
		pl = (M)(tls)
	default:
		panic("value of nested translation must be string.")
	}

	return pl
}
