package export

type ColumnNumber int

const (
	A ColumnNumber = iota
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
	L
	M
	N
	O
	P
	Q
	R
	S
	T
	U
	V
	W
	X
	Y
	Z
	AA
	AB
	AC
)

// 只拓展到两位列ZZ
func (c ColumnNumber) Sting() string {
	if c < 26 {
		return string(rune(c + 'A'))
	}

	return (c/26 - 1).Sting() + (c % 26).Sting()
}

var ColumnLetter = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC"}

const Style = `{"alignment":{"horizontal":"center","vertical":"center"}}`
const HeaderStyle = `{"font":{"bold":true},"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"000000","style":1}],"alignment":{"horizontal":"center","vertical":"center"},"fill":{"type":"pattern","color":["#bfbfbf"],"pattern":1}}`
