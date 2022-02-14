package others

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

type Param map[string]interface{}
type Show struct {
	Param
}

func TestMap(t *testing.T) {
	param := make(Param, 1)
	param["1"] = 1
	fmt.Println(param)

	show := new(Show)
	show1 := Show{}
	fmt.Println(show)
	fmt.Println(show1)

	show.Param = make(Param) //map 创建方式1
	show.Param["111"] = 1
	//show1.Param["111"] = 1 //报错

	m := new(map[interface{}]int) //map 创建方式2，但要初始化
	fmt.Println(m == nil)
	//(*m)[1] = 1 //panic: assignment to entry in nil map
	fmt.Println(m)
	m1 := map[int]int{} //需要这种方法初始化 //map 创建方式3
	m1[1] = 1
	fmt.Println(m1)
}

func TestMapType(t *testing.T) {
	show := Show{}
	//constant 100000000000000000000000 overflows time.Duration
	//time.Sleep(time.Second * 100000000000000) //会报错，超出int64
	mapType(&show)
}

func mapType(v interface{}) {
	switch v.(type) {
	case Show:
		fmt.Println(v)
	case *Show:
		fmt.Println(v)
	default:
		fmt.Println("----")
	}
}

type Student struct {
	name string
	age  int
}

func TestMapValue(t *testing.T) {
	m := map[int]Student{}
	m[1] = Student{name: "1"}
	m[2] = Student{name: "1"}
	//m[1].name = "112" //编译报错，map的value是不可寻址

	m1 := map[int]*Student{1: {name: "1"}}
	fmt.Printf("m1 %v,m1[1] %v \n", m1, m1[1])
	m1[1].name = "2"
	fmt.Printf("m1 %v,m1[1] %v \n", m1, m1[1])

	stus := []Student{
		{name: "zhou", age: 24}, {name: "li", age: 23}, {name: "wang", age: 22},
	}
	fmt.Println(stus)

	for _, stu := range stus {
		m[stu.age] = stu
	}
	fmt.Println(m[24], m[22])

	for _, stu := range stus { //range stus 会先copy 下，stu是个变量
		m1[stu.age] = &stu //这个变量会一直变化，但引用确实一个
	}
	fmt.Println(m1[24], m1[22])

	for _, stu := range stus {
		s := stu
		m1[stu.age] = &s //s是每次都会创建新的
	}
	fmt.Println(m1[24], m1[22])

}

func TestOverflow(t *testing.T) {
	var i byte
	i = 255
	fmt.Println(i + 1) // byte 其实被 alias 到 uint8 上了
	for i = 0; i <= 255; i++ {
		fmt.Println(i)
	}
}

// 测试下，只允许一个协程，代码new了很多协程，哪些先执行
func TestGOMAXPROCS(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() { //闭包函数的引用，i永远是同一个
			fmt.Println("i: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("j: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	//分析
	//这个输出结果决定来自于调度器优先调度哪个G。
	// 从runtime的源码newproc可以看到，当创建一 个G时，
	// 会优先放入到下一个调度的 runnext 字段上作为下一次优先调度的G。
	// 因此， 最先输出的是最后创建的G，也就是9.
}

func TestEQ(t *testing.T) {
	s := Student{name: "1"}
	s1 := Student{name: "1"}
	fmt.Println(s == s1) //指针比指针，非指针比各个属性值

	ss := &Student{name: "1"}
	ss1 := &Student{name: "1"}
	fmt.Println(ss == ss1) //指针比指针

	//数组只能与相同纬度⻓度以及类型的其他数组比较，切片之间不能直接比较
	//[...] 数组表示方式
	fmt.Println([...]string{"1"} == [...]string{"1"}) //
	a := []string{"1"}
	b := []string{"1"}
	fmt.Println(&a == &b)
}
