package decimal

import "fmt"

type Decimal struct {
	form byte
	neg  bool
	mant []byte
	exp  int
}

func (d *Decimal) String() string {
	if len(d.mant) == 0 {
		return "0"
	}

	var buf []byte
	switch {
	case d.exp <= 0:
		// 0.00ddd
		buf = append(buf, "0."...)
		buf = appendZeros(buf, -d.exp)
		buf = append(buf, d.mant...)

	case /* 0 < */ d.exp < len(d.mant):
		// dd.ddd
		buf = append(buf, d.mant[:d.exp]...)
		buf = append(buf, '.')
		buf = append(buf, d.mant[d.exp:]...)

	default: // len(x.mant) <= x.exp
		// ddd00
		buf = append(buf, d.mant...)
		buf = appendZeros(buf, d.exp-len(d.mant))
	}

	return string(buf)
}

func appendZeros(buf []byte, n int) []byte {
	for ; n > 0; n-- {
		buf = append(buf, '0')
	}
	return buf
}

func (d *Decimal) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	return d.form, d.neg, d.mant, int32(d.exp)
}

func (d *Decimal) Compose(form byte, negative bool, coefficient []byte, exponent int32) error {
	switch form {
	default:
		return fmt.Errorf("unknown form %d", form)
	case 1, 2:
		d.form = form
		d.neg = negative
		return nil
	case 0:
	}
	d.form = form
	d.neg = negative
	d.exp = int(exponent)

	d.mant = coefficient

	return nil
}
