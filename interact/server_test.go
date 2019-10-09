package interact

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var port = ":4321"

func TestTcpServer(t *testing.T) {
	fmt.Println("start tcp server....")
	listener, err := net.Listen("tcp", port)
	fmt.Printf("listener:%s, err:%s", listener, err)
	// 监听并接受来自客户端的连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go doServerStuff(conn)
	}
}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		fmt.Printf("n:%x,buf-str:%s", n, string(buf))
		if err != nil {
			fmt.Println("Error reading", err.Error())
			conn.Write([]byte("Error reading"))
			return //终止程序
		}
		//
		if bytes.Compare([]byte{buf[0], buf[1]}, []byte{13, 10}) == 0 {
			conn.Write([]byte("pls send data"))
			conn.Close() //终止程序
			return
		}
		fmt.Printf("Received data: %v", string(buf))
		conn.Write([]byte(fmt.Sprintf("Received data: %v", string(buf))))
	}
}

func TestTcpClient(t *testing.T) {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		//由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}
	conn.SetDeadline(time.Now().Add(10 * time.Second)) //设置发送接收数据超时

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("First, what is your name?")
	clientName, _ := inputReader.ReadString('\n')

	trimmedClient := strings.Trim(clientName, "\r\n") // Windows 平台下用 "\r\n"，Linux平台下使用 "\n"
	// 给服务器发送信息直到程序退出：
	for {
		fmt.Println("What to send to the server? Type Q to quit.")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\r\n")
		// fmt.Printf("input:--s%--", input)
		// fmt.Printf("trimmedInput:--s%--", trimmedInput)
		if trimmedInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(trimmedClient + " says: " + trimmedInput))
	}
}

func TestUdpClient(t *testing.T) {
	conn, err := net.Dial("tcp", port) // tcp ipv4
	checkConnection(conn, err)
	conn, err = net.Dial("udp", "www.baidu.com") // udp
	checkConnection(conn, err)
	conn, err = net.Dial("tcp", "[2620:0:2d0:200::10]:80") // tcp ipv6
	checkConnection(conn, err)
}
func checkConnection(conn net.Conn, err error) {
	if err != nil {
		fmt.Printf("error %v connecting!", err)
		os.Exit(1)
	}
	fmt.Printf("Connection is made with %v", conn)
}

func TestSocket(t *testing.T) {
	var (
		host          = "www.apache.org"
		port          = "80"
		remote        = host + ":" + port
		msg    string = "GET / \n"
		data          = make([]uint8, 4096)
		read          = true
		count         = 0
	)
	// 创建一个socket
	con, err := net.Dial("tcp", remote)
	// 发送我们的消息，一个http GET请求
	io.WriteString(con, msg)
	// 读取服务器的响应
	for read {
		count, err = con.Read(data)
		read = (err == nil)
		fmt.Printf(string(data[0:count]))
	}
	con.Close()
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside HelloServer handler")
	fmt.Fprintf(w, "Hello,"+req.URL.Path[1:])
}
func TestHttpServer(t *testing.T) {
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func TestHttpHead(t *testing.T) {
	var urls = []string{
		"http://www.google.com/",
		"http://golang.org/",
		"http://blog.golang.org/",
	}
	// Execute an HTTP HEAD request for all url's
	// and returns the HTTP status string or an error string.
	for _, url := range urls {
		resp, err := http.Head(url)
		if err != nil {
			fmt.Println("Error:", url, err)
		}
		fmt.Print(url, ": ", resp.Status)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Get : %v", err)
	}
}
func TestHttpGet(t *testing.T) {
	res, err := http.Get("http://www.google.com")
	checkError(err)
	data, err := ioutil.ReadAll(res.Body)
	checkError(err)
	fmt.Printf("Got: %q", string(data))

	// 初始化XML返回值的结构
	user := User{xml.Name{"", "user"}, Status{""}}
	// 将XML解析为我们的结构
	xml.Unmarshal(data, &user)
	fmt.Printf("status: %s", user.Status.Text)
}

func Test10(t *testing.T) {
	fmt.Println(string(byte(13)))
	fmt.Printf("13:%x,10:%x", byte(13), byte(10))
}

type Status struct {
	Text string
}

type User struct {
	XMLName xml.Name
	Status  Status
}

const form = `
    <html><body>
        <form action="#" method="post" name="bar">
            <input type="text" name="in" />
            <input type="submit" value="submit"/>
        </form>
    </body></html>
`

/* handle a simple get request */
func SimpleServer(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "<h1>hello, world</h1>")
}

func FormServer(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch request.Method {
	case "GET":
		/* display the form to the user */
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, form)
	case "POST":
		/* handle the form data, note that ParseForm must
		   be called before we can extract form data */
		//request.ParseForm();
		//io.WriteString(w, request.Form["in"][0])
		io.WriteString(w, request.FormValue("in"))
	}
}

func TestHttpForm(t *testing.T) {
	http.HandleFunc("/test1", SimpleServer)
	http.HandleFunc("/test2", FormServer)
	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}
