package todo

var A Wb
var B Wb

type Wb struct {
	Obj *int
}

func simpleSet(c *int) {
	A.Obj = nil
	B.Obj = c

	//if GC Begin
	A.Obj = c
	B.Obj = nil
	//scan B
}

//安装：go get -u github.com/google/gops
//代码显示注入：agent.Listen(agent.Options{
//		ShutdownCleanup: true, // automatically closes on os.Interrupt
//	});
//可以使用gops:=golang版( jps + jstack + jstat + jinfo )
