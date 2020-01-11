package decimal

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
)

type Decimal1 struct {
	Int uint64
	dec uint64
}

//居然没有运算符重载

func New(Dec string) (dec Decimal1, err error) {
	nums := strings.Split(Dec, ".")
	if len(nums) == 1 {
		dec.Int, err = strconv.ParseUint(nums[0], 10, 64)
		if err != nil {
			return Decimal1{}, err
		}
	}

	if len(Dec) > 19 {
		return Decimal1{}, errors.New("小数最多19位")
	}
	log.Info(reverse(Dec))
	dec.dec, err = strconv.ParseUint(reverse(Dec), 10, 64)
	if err != nil {
		return Decimal1{}, err
	}
	return dec, nil
}

func reverse(s string) string {
	bytes := make([]byte, len(s), len(s))
	for i := range bytes {
		bytes[i] = s[len(s)-i-1]
	}
	return string(bytes)
}

func (d *Decimal1) Add(v Decimal1) {
	d.Int += v.Int
	d.dec += v.dec
}

func (d *Decimal1) Multi(v Decimal1) {
	/*	i := d.Int * v.Int
		Decimal := d.Decimal * v.Int
		Decimal = Decimal + d.Int*v.Decimal + (d.Decimal*d.Decimal)/(int(exponent(10, uint64(d.effective*2))))
		i = i + Decimal/(int(exponent(10, uint64(d.effective))))
		d.Int = i
		d.Decimal = Decimal % (int(exponent(10, uint64(d.effective))))*/
}

func (d Decimal1) String() string {
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

type decimal struct {
	mant []byte // mantissa ASCII digits, big-endian
	exp  int    // exponent
}

func (x *decimal) String() string {
	if len(x.mant) == 0 {
		return "0"
	}

	var buf []byte
	switch {
	case x.exp <= 0:
		// 0.00ddd
		buf = append(buf, "0."...)
		buf = appendZeros(buf, -x.exp)
		buf = append(buf, x.mant...)

	case /* 0 < */ x.exp < len(x.mant):
		// dd.ddd
		buf = append(buf, x.mant[:x.exp]...)
		buf = append(buf, '.')
		buf = append(buf, x.mant[x.exp:]...)

	default: // len(x.mant) <= x.exp
		// ddd00
		buf = append(buf, x.mant...)
		buf = appendZeros(buf, x.exp-len(x.mant))
	}

	return string(buf)
}

func appendZeros(buf []byte, n int) []byte {
	for ; n > 0; n-- {
		buf = append(buf, '0')
	}
	return buf
}

type Decimal2 big.Float

/*func (d *Decimal2) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	coef := make([]byte, 16)
	copy(coef, d.coefficient[:])
	return d.form, d.neg, coef, d.exponent
}

func (d *Decimal2) Compose(form byte, negative bool, coefficient []byte, exponent int32) error {
	f :=(*big.Float)(d)
	negative =f.Signbit()

}*/

type Decimal3 struct {
	mant uint64
	exp  int
	neg  bool
}
