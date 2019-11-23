package utils

import "fmt"

type Decimal struct {
	Int int
	Dec int
	effective int8
}

//居然没有运算符重载

func New(Int int,Dec int, eff int8) Decimal {
	return Decimal{
		Int:       Int,
		Dec:       Dec,
		effective: eff,
	}
}

func (d *Decimal)Add(v Decimal)  {
	d.Int +=v.Int
	d.Dec +=v.Dec
}

func (d *Decimal)Multi(v Decimal)  {
	i:= d.Int * v.Int
	dec:= d.Dec * v.Int
	dec = dec + d.Int * v.Dec + (d.Dec * d.Dec)/(int(exponent(10,uint64(d.effective*2))))
	i = i + dec/(int(exponent(10,uint64(d.effective))))
	d.Int = i
	d.Dec = dec%(int(exponent(10,uint64(d.effective))))
}

func (d Decimal)String() string {
	return fmt.Sprintf("%d.%d",d.Int,d.Dec)
}

func exponent (a,n uint64) uint64 {
	result := uint64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}
