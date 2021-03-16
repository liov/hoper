buntdb 本地db
freecache 号称0gc
gcache 用的人少
go-cache 实现简单,最快的

ristretto 宣称性能比freecache好,为啥确定存在的值还有命中失败的时候

BenchmarkFree
BenchmarkFree-4        	 4511440	       259.5 ns/op
BenchmarkGCache
BenchmarkGCache-4      	 8045618	       152.3 ns/op
BenchmarkGoCache
BenchmarkGoCache-4     	12903280	        95.33 ns/op
BenchmarkRistretto
BenchmarkRistretto-4   	 2098572	       476.9 ns/op