package decimal

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
)

//放弃，省空间，但计算时间浪费，来回不停转字符串
type Decimal1 struct {
	Int uint64
	//小数部分翻转 0.001 =》 100
	dec uint64
}

//居然没有运算符重载

func New1(Dec string) (dec Decimal1, err error) {
	nums := strings.Split(Dec, ".")
	dec.Int, err = strconv.ParseUint(nums[0], 10, 64)
	if len(nums) == 1 {
		if err != nil {
			return
		}
		return
	}

	if len(nums[1]) > 19 {
		err = errors.New("小数最多19位")
		log.Error(err)
		return
	}

	err = dec.SetDec(nums[1])
	if err != nil {
		return
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

func (d *Decimal1) Dec() string {
	dec := strconv.FormatUint(d.dec, 10)
	return reverse(dec)
}

func (d *Decimal1) DecInt() uint64 {
	dec, _ := strconv.ParseUint(d.Dec(), 10, 64)
	return dec
}

func (d *Decimal1) SetDec(dec string) error {
	var err error
	dec = reverse(dec)
	d.dec, err = strconv.ParseUint(dec, 10, 64)
	return err
}

func (d *Decimal1) SetDecInt(dec uint64) error {
	var err error
	decStr := strconv.FormatUint(dec, 10)
	decStr = reverse(decStr)
	d.dec, err = strconv.ParseUint(decStr, 10, 64)
	return err
}

func (d *Decimal1) Add(v Decimal1) Decimal1 {
	var dec Decimal1
	dec.Int = d.Int + v.Int
	d1 := d.Dec()
	d2 := v.Dec()
	if i := len(d1) - len(d2); i > 0 {
		d2 = d2 + strings.Repeat("0", i)
	} else {
		d1 = d1 + strings.Repeat("0", -i)
	}
	decStr := strconv.FormatUint(d.DecInt()+v.DecInt(), 10)

	if len(decStr)-len(d.Dec()) > 0 {
		dec.SetDec(decStr[1:])
		dec.Int += 1
	} else {
		dec.SetDec(decStr)
	}

	return dec
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
	return fmt.Sprintf("%d.%s", d.Int, d.Dec())
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

type Decimal11 struct {
	Int       uint64
	dec       uint64
	effective int
}

func New11(Dec string, eff int) (dec *Decimal11, err error) {
	dec = new(Decimal11)
	nums := strings.Split(Dec, ".")
	dec.Int, err = strconv.ParseUint(nums[0], 10, 64)
	if len(nums) == 1 {
		if err != nil {
			return
		}
		return
	}

	if eff > 19 || len(nums[1]) > 19 {
		err = errors.New("小数最多19位")
		log.Error(err)
		return
	}
	dec.effective = eff
	if len(nums[1]) >= eff {
		nums[1] = nums[1][0:eff]
	} else {
		nums[1] = nums[1] + strings.Repeat("0", eff-len(nums[1]))
	}
	dec.dec, err = strconv.ParseUint(nums[1], 10, 64)
	return
}

func (x *Decimal11) String() string {
	dec := strconv.FormatUint(x.dec, 10)
	dec = dec + strings.Repeat("0", x.effective-len(dec))
	return fmt.Sprintf("%d.%s", x.Int, dec)
}

func (x *Decimal11) Add(v *Decimal11) *Decimal11 {
	var dec = *x
	dec.Int += v.Int
	if x.effective > v.effective {
		dec.effective = v.effective
		dec.dec = x.dec / exponent(uint64(x.effective-v.effective), 10)
	} else if x.effective < v.effective {
		dec.dec = x.dec / exponent(uint64(v.effective-x.effective), 10)
	}
	d := dec.dec + v.dec
	dStr := strconv.FormatUint(d, 10)

	if len(dStr) > x.effective {
		dec.dec, _ = strconv.ParseUint(dStr[1:], 10, 64)
		dec.Int += 1
	} else {
		dec.dec = d
	}
	return &dec
}

type Decimal3 struct {
	mant uint64
	exp  int
	neg  bool
}

func New3(mant uint64, exp int, neg bool) *Decimal3 {
	return &Decimal3{
		mant: mant,
		exp:  exp,
		neg:  neg,
	}
}

func (x *Decimal3) Add(v Decimal3) *Decimal3 {
	var dec = *x
	if x.exp < 0 || v.exp < 0 {
		if x.exp > v.exp {
			dec.exp = v.exp
			dec.mant = dec.mant * uint64(math.Pow10(x.exp-v.exp))
		} else if x.exp < v.exp {
			v.mant = v.mant * uint64(math.Pow10(v.exp-x.exp))
			v.exp = x.exp
		}
	}

	if x.neg == v.neg {
		dec.mant += v.mant
	} else {
		if x.mant >= v.mant {
			dec.mant -= v.mant
		} else {
			dec.mant = v.mant - dec.mant
			dec.neg = v.neg
		}
	}

	return &dec
}

func (x *Decimal3) Sub(v Decimal3) *Decimal3 {
	v.neg = !v.neg
	return x.Add(v)
}

func (x *Decimal3) Mul(v Decimal3) *Decimal3 {
	if x.mant == 0 || v.mant == 0 {
		return &Decimal3{}
	}
	v.mant *= x.mant
	v.exp += x.exp
	return &v
}

func (x *Decimal3) Div(v Decimal3) *Decimal3 {
	if x.mant == 0 {
		return &Decimal3{}
	}
	if v.mant == 0 {
		panic("除数不能为0")
	}
	if v.exp == 0 {
		return &*x
	}
	d1 := x.mant
	d2 := v.mant

	if v.exp < 0 {
		d1 = x.mant * uint64(math.Pow10(0-v.exp))
	} else {
		v.exp = x.exp - v.exp
	}

	v.mant = d1 / d2
	return &v
}

func (x *Decimal3) String() string {
	if x.mant == 0 {
		return "0"
	}
	d := x.mant
	if x.exp > 0 {
		d = x.mant * uint64(math.Pow10(x.exp))
	}

	str := strconv.FormatUint(d, 10)
	var in, de string
	if x.neg == true {
		in = "-"
	}
	if len(str)+x.exp < 0 {
		in += "0"
		de = "." + strings.Repeat("0", -x.exp-len(str)) + str
	} else {
		if x.exp >= 0 {
			in += str
		} else {
			in += str[:len(str)+x.exp]
			de = "." + str[len(str)+x.exp:]
		}

	}

	for i := len(de) - 1; i >= 0; i-- {
		if de[i] == '0' {
			x.mant /= 10
			x.exp += 1
			de = de[:len(de)-1]
		} else {
			break
		}
	}
	return fmt.Sprintf("%s%s", in, de)
}
