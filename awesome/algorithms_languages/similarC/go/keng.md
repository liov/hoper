交叉编译的bug
正常编译可以，交叉编译就报包找不到(cannot find module for path github.com/360EntSecGroup-Skylar/excelize)
main里下划线导入不报找不到包(https://juejin.im/post/5d776830f265da03e05b3c45),内部包找不到了
cgo的锅
set CGO_ENABLED=1
测试不是github.com/360EntSecGroup-Skylar/excelize/v2的锅

应该是cgo的原因，但是那个项目里的包都是常见的包啊，难以定位哪里用了带cgo的包

交叉编译时，CGO_ENABLED=0是会自动忽略带cgo的包，这个有bug，1.14会修复[https://github.com/golang/go/issues/35873]

main包匿名导入提示找不到路径的包又不报这个错，报内部包的函数undefine
无法复现

排查了半天，真的让人哭笑不得
真的跟cgo有关
那个引用找不到路径的包的包多了个import "C"，不知道什么时候加上去的

---p1.go
package p1

import "C"
import github.com/user/p2

---go.mod
github.com/user/p2
