package common

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var digitRegexp = regexp.MustCompile("[0-9]+")

const NBUF = 512

//CatFile data from file write to the Writer
func CatFile(fileName string, w io.Writer) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(w, "error reading from %s: %s\n", fileName, err.Error())
		return
	}
	//r := bufio.NewReader(f)
	//for {
	//	buf, err := r.ReadBytes('\n')
	//	if err == io.EOF {
	//		break
	//	}
	//	time.Sleep(20 * time.Millisecond)
	//	fmt.Fprintf(w, "%s", buf)
	//}

	var buf [NBUF]byte
	for {
		switch nr, err := f.Read(buf[:]); true {
		case nr < 0:
			fmt.Fprintf(w, "cat: error reading: %s\n", err.Error())
			return
		case nr == 0: // EOF
			return
		case nr > 0:
			if nw, ew := w.Write(buf[0:nr]); nw != nr {
				fmt.Fprintf(w, "cat: error writing: %s\n", ew.Error())
			}
		}
	}
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

//FindDigits find digits from file
//切片的底层指向一个数组，该数组的实际体积可能要大于切片所定义的体积。
//只有在没有任何切片指向的时候，底层的数组内层才会被释放，这种特性有时会导致程序占用多余的内存。
func FindDigits(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	b = digitRegexp.Find(b)
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func Mkdir(path string) error {
	return os.Mkdir(path, 0777)
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	// os.Stat获取文件信息
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func ReadFile(filename string) []string {
	var words []string

	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}

		wordSlice := strings.Fields(line)
		for _, word := range wordSlice {
			if word = extractStr(strings.ToLower(word)); word != "" {
				words = append(words, word)
			}
		}
	}

	return words
}

func extractStr(str string) string {
	var res []rune
	for _, letter := range str {
		if letter >= 'a' && letter <= 'z' {
			res = append(res, letter)
		}
	}
	return string(res)
}
