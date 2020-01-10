package decimal

import (
	"errors"
	"fmt"
	"strconv"
)

type Decimal struct {
	Int int
	dec uint64
}

//居然没有运算符重载

func New(Int int, Dec string) (Decimal, error) {
	var dec = Decimal{
		Int: Int,
	}
	if len(Dec) > 10 {
		return Decimal{}, errors.New("小数最多10位")
	}
	dec.dec, _ = strconv.ParseUint(reverse(Dec), 10, 64)
	return dec, nil
}

func reverse(s string) string {
	bytes := make([]byte, len(s))
	for b := range s {
		bytes = append(bytes, byte(b))
	}
	return string(bytes)
}

func (d *Decimal) Add(v Decimal) {
	d.Int += v.Int
	d.dec += v.dec
}

func (d *Decimal) Multi(v Decimal) {
	/*	i := d.Int * v.Int
		dec := d.dec * v.Int
		dec = dec + d.Int*v.dec + (d.dec*d.dec)/(int(exponent(10, uint64(d.effective*2))))
		i = i + dec/(int(exponent(10, uint64(d.effective))))
		d.Int = i
		d.dec = dec % (int(exponent(10, uint64(d.effective))))*/
}

func (d Decimal) String() string {
	return fmt.Sprintf("%d.%d", d.Int, d.dec)
}

func exponent(a, n uint64) uint64 {
	result := uint64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}
