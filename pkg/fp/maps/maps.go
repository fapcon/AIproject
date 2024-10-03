package maps

import (
	"golang.org/x/exp/maps"
	"studentgit.kata.academy/quant/torque/pkg/fp"
)

func Map1[K1, K2 comparable, V1, V2 any](m map[K1]V1, convert func(K1, V1) (K2, V2)) map[K2]V2 {
	res := make(map[K2]V2, len(m))
	for k1, v1 := range m {
		k2, v2 := convert(k1, v1)
		res[k2] = v2
	}
	return res
}

func Map2[K1, K2 comparable, V1, V2 any](m map[K1]V1, key func(K1) K2, value func(V1) V2) map[K2]V2 {
	res := make(map[K2]V2, len(m))
	for k, v := range m {
		res[key(k)] = value(v)
	}
	return res
}

func MapKeys[K1, K2 comparable, V any](m map[K1]V, convert func(K1) K2) map[K2]V {
	return Map2(m, convert, fp.Identity[V])
}

func MapValues[K comparable, V1, V2 any](m map[K]V1, convert func(V1) V2) map[K]V2 {
	return Map2(m, fp.Identity[K], convert)
}

func ToSlice[K comparable, V, E any](m map[K]V, convert func(K, V) E) []E {
	res := make([]E, 0, len(m))
	for k, v := range m {
		res = append(res, convert(k, v))
	}
	return res
}

func Keys[K comparable, V any](m map[K]V) []K {
	return maps.Keys(m)
}

func Values[K comparable, V any](m map[K]V) []V {
	return maps.Values(m)
}

func Clone[K comparable, V any](m map[K]V) map[K]V {
	return maps.Clone(m)
}
