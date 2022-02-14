package todo

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

//病毒路径和md5
type virusSignature struct {
	path string
	md5  string
}

var virusDetection bool

//virus md5 list
var applicationMd5List = []string{
	"671404b85bd72598d8a38e05d96cb8d7",
	"b999aed32ffbf9f5c3265f1cf29ea552",
	"2b8638d1a6f3fac3c99893b878bf0bb3",
	"a108cfdc6960eb103df6ccc6107e444f",
	"63110145a6335c2952e225caad679db0",
	"b5caf2728618441906a187fc6e90d6d5",
	"1095e840b31609c1178164f405f65020",
	"4e245480d5f4136a49f6b0f50c6152c4",
	"81866e7a8bec80097a9d07355e5e67ff",
	"3355134160aa4a5bc86fe42185ef0677",
}

var applicationList = []string{
	"FinalShell",
	"Jump",
	"Remote",
	"Navicat",
	"SecureCRT",
	"Sourcetree",
	"iTerm",
}

func TestVirusScan(t *testing.T) {
	fmt.Println("正在检测是否存在病毒...")
	virusDetection = false
	//检查应用列表是否存在病毒
	listApplication()
	//检查是否已经被植入后门
	googleUpdate1 := new(virusSignature)
	googleUpdate1.path = "/tmp/GoogleUpdate"
	googleUpdate1.md5 = "96013240f62f846de82304fbcad8b653"
	googleUpdate2 := new(virusSignature)
	googleUpdate2.path = "/tmp/GoogleUpdate"
	googleUpdate2.md5 = "47d774e5307215c7c11151211c8d3ce2"
	gpy := new(virusSignature)
	gpy.path = "/tmp/g.py"
	gpy.md5 = "2786ebc3b917866d30e622325fc6f5f3"
	detection(googleUpdate1)
	detection(googleUpdate2)
	detection(gpy)
	if !virusDetection {
		fmt.Println("[安全]您的Mac电脑暂未安装风险软件")
	} else {
		fmt.Println("[危险]您的Mac电脑存在风险软件，请联系涂鸦安全进行排查!")
	}
}

//判断是否被植入后门
func detection(virus *virusSignature) {

	if CalcFileMD5(virus.path) == virus.md5 {
		virusDetection = true
	}
}

// 计算文件MD5

func CalcFileMD5(filename string) string {
	f, err := os.Open(filename) //打开文件
	if nil != err {
		return ""
	}
	defer f.Close()

	md5Handle := md5.New()         //创建 md5 句柄
	_, err = io.Copy(md5Handle, f) //将文件内容拷贝到 md5 句柄中
	if nil != err {
		fmt.Println(err)
		return ""
	}
	md := md5Handle.Sum(nil)        //计算 MD5 值，返回 []byte
	md5str := fmt.Sprintf("%x", md) //将 []byte 转为 string
	return md5str
}

//读取Application列表
func listApplication() {
	baseDir := "/Applications"
	fileInfos, err := ioutil.ReadDir(baseDir)
	if err == nil {
		for _, fir := range fileInfos {
			filename := baseDir + "/" + fir.Name() + "/Contents/Frameworks/"
			packageInfos, e := ioutil.ReadDir(filename)
			if e == nil {
				for _, f := range packageInfos {
					if strings.Contains(f.Name(), "libcrypto") {
						applicationMD5 := CalcFileMD5(filename + f.Name())
						if in(applicationMD5, applicationMd5List) || infile(filename+f.Name(), "erdou") {
							fmt.Println(fmt.Sprintf("【危险】发现存在风险软件：%s,请联系涂鸦安全！", fir.Name()))
							virusDetection = true
						}
					}
				}

			}

		}
	} else {
		fmt.Println(fmt.Sprintf("读取应用列表错误！%s", err))
	}
}

//判读md5是否在数组中
func in(md5 string, md5List []string) bool {
	for _, element := range md5List {
		if md5 == element {
			return true
		}
	}
	return false
}

//判读是否存在关键字
func infile(path string, word string) bool {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(b)
	// //check whether s contains substring text
	return strings.Contains(s, word)
	return false
}
