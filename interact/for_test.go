package interact

import (
	"fmt"
	"sync"
	"testing"
)

//通过查看go 编译源码[1]可以了解到, for-range 其实是语法糖，
//内部调用还是 for 循环，初始化会拷贝带遍历的列表（如 array，slice，map），
//然后每次遍历的v都是对同一个元素的遍历赋值。也就是说如果直接对v取地址，
//最终只会拿到一个地址，而对应的值就是最后遍历的那个元素所附给v的值。
func TestFor(t *testing.T) {
	//遍历取不到所有元素指针?
	arr := [2]int{1, 2}
	res := []*int{}
	for _, v := range arr {
		res = append(res, &v) //保存的是一个临时遍历地址
	}
	//expect: 1 2
	fmt.Println(*res[0], *res[1])
	//but output: 2 2
}

//遍历会停止么
func TestFor2(t *testing.T) {
	v := []int{1, 2, 3}
	for i := range v {
		fmt.Print(i)
		v = append(v, i) //期间对原来v的修改不会反映到遍历中
	}
	fmt.Println("end?", v)
}

//对大数组这样遍历有啥问题
//答案是【有问题】！遍历前的拷贝对内存是极大浪费啊 怎么优化？有两种
//对数组取地址遍历for i, n := range &arr
//对数组做切片引用for i, n := range arr[:]
//反思题：对大量元素的 slice 和 map 遍历为啥不会有内存浪费问题？
func TestFor3(t *testing.T) {
	//假设值都为1，这里只赋值3个
	var arr = [102400]int{1, 1, 1}
	for i, n := range arr {
		//just ignore i and n for simplify the example
		_ = i
		_ = n
	}

	//对大数组这样重置效率高么
	//go 对这种重置元素值为默认值的遍历是有优化的
	for i, _ := range &arr {
		arr[i] = 0
	}
}

// 对 map 遍历时删除元素能遍历到么
//答案是【不会】 map 内部实现是一个链式 hash 表，为保证每次无序，
//初始化时会随机一个遍历开始的位置[3], 这样，
// 如果删除的元素开始没被遍历到（上边once.Do函数内保证第一次执行时删除未遍历的一个元素），
// 那就后边就不会出现
func TestForMap(t *testing.T) {
	var m = map[int]int{1: 1, 2: 2, 3: 3}
	//only del key once, and not del the current iteration key
	var o sync.Once
	for i := range m {
		o.Do(func() {
			for _, key := range []int{1, 2, 3} {
				if key != i {
					fmt.Printf("when iteration key %d, del key %d\n", i, key)
					delete(m, key)
					break
				}
			}
		})
		fmt.Printf("%d%d ", i, m[i])
	}

	//这样遍历中起 goroutine 可以么
	//1\使用局部变量拷贝 2\以参数方式传入
	for i := range m {
		go func(i int) {
			fmt.Print(i)
		}(i)
	}
}
