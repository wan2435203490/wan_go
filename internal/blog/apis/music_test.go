package apis

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
	"wan_go/internal/blog/vo/music"
)

// BenchmarkA-8表示执行 BenchmarkA 时，所用的最大P的数量为8
// 1: 表示hello()方法在达到这个执行次数时，等于或超过了1秒
// 1009659209 ns/op： 表示每次执行A()所消耗的平均执行时间
// 2.386s：表示测试总共用时
func BenchmarkA(b *testing.B) {
	//BenchmarkA-8   	       		1	1 009 659 209 ns/op
	//100 BenchmarkA-8   	       	1	10 099 091 000 ns/op
	for i := 0; i < b.N; i++ {
		A()
	}
	b.Log("NNNNN:", b.N)
}
func BenchmarkB(b *testing.B) {
	//BenchmarkB-8   	      	  10	 101 028 358 ns/op
	//100 BenchmarkB-8   	      10	 101 992 004 ns/op
	//10000 BenchmarkB-8   	       9	 111 919 273 ns/op
	for i := 0; i < b.N; i++ {
		B()
	}
	b.Log("NNNNN:", b.N)
}
func BenchmarkC(b *testing.B) {
	//BenchmarkC-8   	      	  10	 101 098 800 ns/op
	//100 BenchmarkC-8   	      10	 101 472 883 ns/op
	//10000 BenchmarkC-8   	       9	 112 539 806 ns/op
	for i := 0; i < b.N; i++ {
		C()
	}
	b.Log("NNNNN:", b.N)
}

func A() {

	count := 10000

	ch := make(chan *music.Song, count)

	go func() {
		for i := 0; i < count; i++ {
			//怎么改成异步
			Request()
			ch <- &music.Song{Name: strconv.Itoa(i)}
		}
		close(ch)
	}()

	result := make([]*music.Song, 0)
	for song := range ch {
		result = append(result, song)
		//fmt.Println(song.Name)
	}

	//close(ch)

	//for _, v := range result {
	//	fmt.Println(v.Name)
	//}
}

func Request() {
	time.Sleep(time.Millisecond * 100)
}

func B() {

	count := 10000

	result := make([]*music.Song, 0)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < count; i++ {
		wg.Add(1)
		//goroutine请求
		go func() {
			//请求
			Request()
			response := &music.Song{Name: strconv.Itoa(i)}
			mu.Lock()
			result = append(result, response)
			mu.Unlock()
			wg.Done()
		}()

	}

	wg.Wait()

}

func C() {

	count := 10000

	ch := make(chan *music.Song, count)

	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		//goroutine请求
		go func() {
			//请求
			Request()
			response := &music.Song{Name: strconv.Itoa(i)}
			ch <- response
			wg.Done()
		}()
	}

	wg.Wait()

	close(ch)

	result := make([]*music.Song, 0)
	for song := range ch {
		result = append(result, song)
		//fmt.Println(song.Name)
	}

}

func genPrimeNumber() (ch chan int) {
	ch = make(chan int, 1)
	go func() {
		for i := 2; ; i++ {
			ch <- i
			fmt.Printf("[gen %v ] \n", i)
		}
	}()
	return
}

func primeNumberFilter(ch <-chan int, p, i int) (out chan int) {
	out = make(chan int, 1)
	go func(i int) {
		for {
			n := <-ch
			if (n % p) != 0 {
				out <- n
			}
		}
	}(i)
	return
}

func TestGet(t *testing.T) {

	runtime.GOMAXPROCS(1)

	ch := genPrimeNumber()

	for i := 0; i < 4; i++ {
		n := <-ch
		ch = primeNumberFilter(ch, n, i)
	}

	fmt.Println()
}
