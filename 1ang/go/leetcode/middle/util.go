package middle

func Abs(a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}

func Swap(heap []int, i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}

func Reverse(bytes []byte) []byte {
	ret := make([]byte, len(bytes))
	l := len(bytes) - 1
	for i := 0; i < len(bytes); i++ {
		ret[i] = bytes[l-i]
	}
	return ret
}

func ReverseInPlace(bytes []byte) {
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1] = bytes[len(bytes)-1], bytes[i]
	}
}

func intSlice(i int, s []int) bool {
	for _, n := range s {
		if n == i {
			return true
		}
	}
	return false
}

var enNumbers = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
var numberStrings = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var numberBytes = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func parent(i int) int {
	return (i - 1) / 2
}
func leftChild(i int) int {
	return i*2 + 1
}

func rightChild(i int) int {
	return i*2 + 2
}
