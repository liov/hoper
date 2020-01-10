package decimal

import "fmt"

type dec struct {
	form        byte
	neg         bool
	coefficient [16]byte
	exponent    int32
}

func (d dec) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	coef := make([]byte, 16)
	copy(coef, d.coefficient[:])
	return d.form, d.neg, coef, d.exponent
}

func (d *dec) Compose(form byte, negative bool, coefficient []byte, exponent int32) error {
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
	d.exponent = exponent

	// This isn't strictly correct, as the extra bytes could be all zero,
	// ignore this for this test.
	if len(coefficient) > 16 {
		return fmt.Errorf("coefficent too large")
	}
	copy(d.coefficient[:], coefficient)

	return nil
}
