package pair

type PairKey string

func NewPairKey(base, quote string) PairKey {
	return PairKey(base + "/" + quote)
}
