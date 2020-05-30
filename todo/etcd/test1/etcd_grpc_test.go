package test1

import (
	"context"
	"flag"
	"fmt"
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/todo/etcd"
	"github.com/Zeb-D/go-util/todo/etcd/pb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

// etcd解读 https://www.infoq.cn/article/etcd-interpretation-application-scenario-implement-principle/

var (
	serv = flag.String("service", "pb.HelloService", "service name")
	port = flag.Int("port", 50001, "listening port")
	reg  = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

func TestEtcdGrpcServer(t *testing.T) {

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}
	localHost, err := common.LocalIp()
	fmt.Println(localHost, " host->", err)
	err = etcd.Register(*serv, localHost, *port, *reg, time.Second*10, 15)
	time.Sleep(time.Second)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		etcd.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %d", *port)
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	s.Serve(lis)
	fmt.Println("start???")
}

func TestGrpcClient(t *testing.T) {
	flag.Parse()
	fmt.Println("serv", *serv)
	r := etcd.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}
	fmt.Println("conn...")

	time.Sleep(time.Second) //	为什么要给点时间建立连接？watcher机制没触发？
	for i := 0; i < 10; i++ {
		client := pb.NewHelloServiceClient(conn)
		resp, err := client.Hello(context.Background(), &pb.String{Value: "world " + strconv.Itoa(time.Now().Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", time.Now(), resp.Value)
		} else {
			fmt.Println("has err->", err)
		}
		time.Sleep(time.Second)
	}
	fmt.Println("client end")
}

func TestSearchRequest(t *testing.T) {
	url := &UrlVO{
		Url:   "https://github.com/Zeb-D",
		Title: "my-review",
	}
	ss := &SearchRequest{RunMode: 32904, BizType: "123", Url: url}
	bs, _ := proto.Marshal(ss)
	fmt.Println(bs)
	// 第一个长度大小、第二个长度开始
	fmt.Println(string(bs))
	bs1, _ := proto.Marshal(url)

	fmt.Println(bs1)
	fmt.Println([]byte(url.Url))

	url1 := &UrlVO{}
	bs2 := []byte{10, 24, 104, 116, 116, 112, 115, 58, 47, 47, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 90, 101, 98, 45, 68, 18, 9, 109, 121, 45, 114, 101, 118, 105, 101, 119}
	proto.Unmarshal(bs2, url1)
	fmt.Println(url1.Url)
}

// server is used to implement pb.HelloServiceServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Hello(ctx context.Context, in *pb.String) (*pb.String, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Value)
	return &pb.String{Value: "Hello " + in.Value}, nil
}

func TestMy(t *testing.T) {
	fmt.Println(1 >> 31)
}
