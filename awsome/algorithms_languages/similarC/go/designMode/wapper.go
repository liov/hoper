package main

type Food struct {
	name string
}

func (f *Food) make() string {
	return f.name
}

type Breed struct {
	Food
}

func (b *Breed) make() string {
	return b.name + "面包"
}

type Cream struct {
	Food
}

func (c *Cream) make() string {
	return c.name + "奶油"
}

//组合大于继承，对于非OOP语言，组合直接就是装饰器啊哈哈哈
func main() {

}
