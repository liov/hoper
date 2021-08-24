package number

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"math/big"
	"testing"
	"time"
	"unsafe"
)

func Test_Dec(t *testing.T) {
	fmt.Printf("%#v", Decimal{exp: 6})
}

const maxUint64 uint64 = 1<<64 - 1

func Test_Decimal(t *testing.T) {
	a, _ := New1("2.555")
	b, _ := New1(" 0.06")
	fmt.Println(a, b, a.Add(b))
	//fmt.Println(new(big.Int).Lsh(big.NewInt(1),100))
	c, _ := New11("2.11", 8)
	d, _ := New11("0.89", 8)
	fmt.Println(c.Add(d))
	e := New3(205, -2, true)
	f := New3(2, 2, false)
	fmt.Println(e.Sub(*f))
	fmt.Println(e.Mul(*f))
	fmt.Println(e.Div(*f, HALFUP))
	g := &Decimal{mant: []byte("123456"), exp: 10}
	fmt.Println(g)
	ViewBin(9007199254741004.0)
	h := uint64(111)
	fmt.Println(*(*float64)(unsafe.Pointer(&h)))
}

func Test_Float(t *testing.T) {
	//var a = 0.1000000000000000055511151231257827021181583404541015625
	a := 0.11
	b := 0.1
	f1 := big.NewFloat(a)
	f2 := big.NewFloat(b)
	f3 := f2.Mul(f1, f2)
	data, _ := f3.MarshalText()
	fmt.Println(a*b, f3, string(data))
}

func Test_DB(t *testing.T) {
	type DecTest struct {
		Id   uint64
		Dec  Decimal3 `gorm:"type:decimal(10,2)"`
		Time time.Time
	}
	db := &gorm.DB{}
	tx := db.Begin()
	/*	tx.DropTable(&DecTest{})
		tx.CreateTable(&DecTest{})*/
	var dec = DecTest{Dec: Decimal3{123, -2, false}}
	tx.Save(&dec)
	log.Println(dec.Id)
	var dec1 DecTest
	tx.First(&dec1)
	tx.Commit()
	fmt.Println(dec1)
}
