module test

go 1.18

require (
	bou.ke/monkey v1.0.2
	github.com/actliboy/hoper/server/go/lib v1.0.0
	github.com/bcongdon/corral v0.0.0-20190319213343-1f4ad49dee34
	github.com/brahma-adshonor/gohook v1.1.9
	github.com/google/flatbuffers v2.0.0+incompatible
	github.com/lucasb-eyer/go-colorful v1.0.2
	github.com/unicorn-engine/unicorn v0.0.0-20191119163456-3cea38bff7bf
	github.com/xuri/excelize/v2 v2.4.1
	golang.org/x/exp v0.0.0-20220314205449-43aec2f8a4e7
	golang.org/x/mobile v0.0.0-20201217150744-e6ae53a27f4f
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654
)

require github.com/mattetti/filebuffer v1.0.0 // indirect

replace github.com/actliboy/hoper/server/go/lib => ../../server/go/lib
