如何从json中反序列化一个带函数的结构体
之前一直以为go的func channel之类的类型是不能反序列化的,原因是不能序列化,直到看到这样的写法

```go
// A LevelEncoder serializes a Level to a primitive type.
type LevelEncoder func(Level, PrimitiveArrayEncoder)

// LowercaseLevelEncoder serializes a Level to a lowercase string. For example,
// InfoLevel is serialized to "info".
func LowercaseLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
    enc.AppendString(l.String())
}

// LowercaseColorLevelEncoder serializes a Level to a lowercase string and adds coloring.
// For example, InfoLevel is serialized to "info" and colored blue.
func LowercaseColorLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
    s, ok := _levelToLowercaseColorString[l]
    if !ok {
    s = _unknownLevelColor.Add(l.String())
    }
    enc.AppendString(s)
}

// CapitalLevelEncoder serializes a Level to an all-caps string. For example,
// InfoLevel is serialized to "INFO".
func CapitalLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
    enc.AppendString(l.CapitalString())
}

// CapitalColorLevelEncoder serializes a Level to an all-caps string and adds color.
// For example, InfoLevel is serialized to "INFO" and colored blue.
func CapitalColorLevelEncoder(l Level, enc PrimitiveArrayEncoder) {
    s, ok := _levelToCapitalColorString[l]
    if !ok {
    s = _unknownLevelColor.Add(l.CapitalString())
    }
    enc.AppendString(s)
}
func (e *LevelEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "capital":
		*e = CapitalLevelEncoder
	case "capitalColor":
		*e = CapitalColorLevelEncoder
	case "color":
		*e = LowercaseColorLevelEncoder
	default:
		*e = LowercaseLevelEncoder
	}
	return nil
}
```
上面是日志库zap里的一段写法,利用接口,在反序列化的过程中生成函数,当然并不是创建了个动态函数,只是赋值了函数指针,利用了go序列化库中都会提供的接口Unmarshalxxx,同理可以写一个UnmarshalJson版本,直接从json为函数赋值,适用于有几个选项的函数