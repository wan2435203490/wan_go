package test

import (
	"fmt"
	"github.com/timandy/routine"
	"testing"
	"time"
)

var threadLocal = routine.NewThreadLocal()
var inheritableThreadLocal = routine.NewInheritableThreadLocal()

func TestRoutine(t *testing.T) {
	goid := routine.Goid()
	fmt.Println(goid)
	go func() {
		i := routine.Goid()
		fmt.Println(i)
	}()

	threadLocal.Set("threadLocal")
	inheritableThreadLocal.Set("inheritableThreadLocal")
	fmt.Println("threadLocal:", threadLocal.Get())
	fmt.Println("inheritableThreadLocal:", inheritableThreadLocal.Get())

	// The child coroutine cannot read the previously assigned "hello world".
	go func() {
		fmt.Println("threadLocal in goroutine:", threadLocal.Get())
		fmt.Println("inheritableThreadLocal in goroutine:", inheritableThreadLocal.Get())
	}()

	//mac 环境下面两个都是nil
	// However, a new sub-coroutine can be started via the Go/GoWait/GoWaitResult function, and all inheritable variables of the current coroutine can be passed automatically.
	routine.Go(func() {
		//nil
		fmt.Println("threadLocal in goroutine by Go:", threadLocal.Get())
		//nil
		fmt.Println("inheritableThreadLocal in goroutine by Go:", inheritableThreadLocal.Get())
	})

	//// You can also create a task via the WrapTask/WrapWaitTask/WrapWaitResultTask function, and all inheritable variables of the current coroutine can be automatically captured.
	//task := routine.WrapTask(func() {
	//	fmt.Println("threadLocal in task by WrapTask:", threadLocal.Get())
	//	fmt.Println("inheritableThreadLocal in task by WrapTask:", inheritableThreadLocal.Get())
	//})
	//go task.Run()

	// Wait for the sub-coroutine to finish executing.
	time.Sleep(time.Second)
}
