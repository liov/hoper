package graphql

//性能强迫症适合用go吗，一个graphql-go/graphql实现，schema定义占一大波内存，完了取值是用过反射遍历字段，
//各种用反射，各种循环，是为了功能又费内存又费cpu
//还不如graph-gophers/graphql-go，虽然定义比较多，在生成schema的反射期记录方法索引，虽然取值也是反射！
//两个实现性能差不多，毫秒间，而且说不上谁性能好，好像graph-gophers/graphql-go更好一些
//本来想基于反射（go写工具已经离不开反射了）写一个graphql-go/graphql自动生成schema
func NewObject() {

}
