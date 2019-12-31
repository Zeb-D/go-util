package etcd

import (
	"context"
	"flag"
	"fmt"
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/todo/etcd/pb"
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
	err = Register(*serv, localHost, *port, *reg, time.Second*10, 15)
	time.Sleep(time.Second)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		UnRegister()
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
	r := NewResolver(*serv)
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

// server is used to implement pb.HelloServiceServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Hello(ctx context.Context, in *pb.String) (*pb.String, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Value)
	return &pb.String{Value: "Hello " + in.Value}, nil
}
