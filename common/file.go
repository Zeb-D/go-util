package common

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

// @param path 带路径的文件名
// @param info 文件的 os.FileInfo 信息
// @param content 文件内容
// @return 返回nil或者content原始内容，则不修改文件内容，返回其他内容，则会覆盖重写文件
type WalkFunc func(path string, info os.FileInfo, content []byte, err error) []byte

// 遍历访问指定文件夹下的文件
// @param root 需要遍历访问的文件夹
// @param recursive 是否递归访问子文件夹
// @param suffix 指定文件名后缀进行过滤，如果为""，则不过滤
func Walk(root string, recursive bool, suffix string, walkFn WalkFunc) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkFn(path, info, nil, err)
			return nil
		}
		if !recursive && info.IsDir() && path != root {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if suffix != "" && !strings.HasSuffix(info.Name(), suffix) {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			walkFn(path, info, content, err)
			return nil
		}
		newContent := walkFn(path, info, content, nil)
		if newContent != nil && bytes.Compare(content, newContent) != 0 {
			if err = ioutil.WriteFile(path, newContent, 0755); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// 文件尾部添加内容
func AddTailContent(content []byte, tail []byte) []byte {
	if !bytes.HasSuffix(content, []byte{'\n'}) {
		content = append(content, '\n')
	}
	return append(content, tail...)
}

// 文件头部添加内容
func AddHeadContent(content []byte, head []byte) []byte {
	if !bytes.HasSuffix(head, []byte{'\n'}) {
		head = append(head, '\n')
	}
	return append(head, content...)
}

// 行号范围
// 1表示首行，-1表示最后一行
type LineRange struct {
	From int
	To   int
}

var ErrLineRange = errors.New("go-util.file: line range error")

func calcLineRange(len int, lr LineRange) (LineRange, error) {
	// 换算成从0开始的下标
	if lr.From < 0 {
		lr.From = len + lr.From
	} else if lr.From > 0 {
		lr.From = lr.From - 1
	} else {
		return lr, ErrLineRange
	}
	if lr.To < 0 {
		lr.To = len + lr.To
	} else if lr.To > 0 {
		lr.To = lr.To - 1
	} else {
		return lr, ErrLineRange
	}

	// 排序交换
	if lr.From > lr.To {
		lr.From, lr.To = lr.To, lr.From
	}

	if lr.From < 0 || lr.From >= len || lr.To < 0 || lr.To >= len {
		return lr, ErrLineRange
	}

	return lr, nil
}

func DeleteLines(content []byte, lr LineRange) ([]byte, error) {
	lines := bytes.Split(content, []byte{'\n'})
	length := len(lines)
	nlr, err := calcLineRange(length, lr)
	if err != nil {
		return content, err
	}
	var nlines [][]byte
	if nlr.From > 0 {
		nlines = append(nlines, lines[:nlr.From]...)
	}
	if nlr.To < length-1 {
		nlines = append(nlines, lines[nlr.To+1:]...)
	}
	return bytes.Join(nlines, []byte{'\n'}), nil
}
