package yap

import "math/big"

func NewFloatFromInt(i int) *big.Float {
	return new(big.Float).SetInt(big.NewInt(int64(i)))
}
