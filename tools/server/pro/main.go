package main

import (
	"os"
	"sync"
	"time"
)

func main() {
	option := 'n'
	start := 330000
	end := 340000
	sd := NewSpeed(loop)
	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)
		f, _ := os.Create(commonDir + "fail_" + time.Now().Format("2006_01_02_15_04_05") + `.txt`)
		for txt := range sd.fail {
			f.WriteString(txt + "\n")
		}
		f.Close()
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		f, _ := os.Create(commonDir + "fail_pic_" + time.Now().Format("2006_01_02_15_04_05") + `.txt`)
		for txt := range sd.failPic {
			f.WriteString(txt + "\n")
		}
		f.Close()
		wg.Done()
	}()
	switch option {
	case 'n':
		normal(start, end, sd)
	case 'g':
		gif(start, end, sd)
	}
	sd.Wait()
	close(sd.fail)
	close(sd.failPic)
	wg.Wait()
}

func normal(start, end int, sd *speed) {
	for i := start; i < end; i++ {
		sd.WebAdd(1)
		go fetch(i, sd)
		time.Sleep(interval)
	}
}
