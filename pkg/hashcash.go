package algo

import "math/big"

type Hcash struct {
	target *big.Int
}

func NewHCash() *Hcash {
	return &Hcash{
		target: big.NewInt(1),
	}
}
