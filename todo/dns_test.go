package todo

import (
	"fmt"
	"net"
	"os"
	"testing"
)

//	strace常用来跟踪进程执行时的系统调用和所接收的信号。
//	在Linux世界，进程不能直接访问硬件设备，当进程需要访问硬件设备(比如读取磁盘文件，接收网络数据等等)时，
//	必须由用户态模式切换至内核态模式，通 过系统调用访问硬件设备。
//	strace可以跟踪到一个进程产生的系统调用,包括参数，返回值，执行消耗的时间。

//	在solaris下，对应的是dtrace
//	在mac下，对应的命令是：dtruss 需要root
//  sudo  dtruss go test -v -run TestDns
func TestDns(t *testing.T) {
	// go/src/net/dnsclient_unix.go
	ns, err := net.LookupHost("www.baidu.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s", err.Error())
		return
	}

	for index, n := range ns {
		_, _ = fmt.Fprintf(os.Stdout, "%v--%s\n", index, n)
	}
}
