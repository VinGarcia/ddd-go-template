package maps

// This package is a helper package:
// (1) it is very simple
// (2) it depends on no other packages
//
// It's so simple in fact that I don't care about decoupling from it,
// so I wont write any interfaces here.

type Body = map[string]interface{}

func Merge(baseMap *Body, maps ...Body) {
	if *baseMap == nil {
		*baseMap = Body{}
	}

	for _, m := range maps {
		for k, v := range m {
			(*baseMap)[k] = v
		}
	}
}
