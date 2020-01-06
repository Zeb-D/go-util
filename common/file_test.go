package common

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func createTestFile() string {
	f, _ := os.Create("test_file")
	return f.Name()
}

func cleanTestFile() {
	_ = os.Remove("test_file")
}

//test file exists
func TestExists(t *testing.T) {
	pwd, err := os.Getwd()
	println("1,", pwd, "->", err)
	assert.Nil(t, err)
	tf := createTestFile()
	defer cleanTestFile()

	p := path.Join(pwd, tf)
	assert.True(t, Exists(p))

	assert.False(t, Exists(p+"_no_exist"))
}

//test file path is dir
func TestIsDir(t *testing.T) {
	pwd, err := os.Getwd()
	assert.Nil(t, err)
	tf := createTestFile()
	defer cleanTestFile()

	p := path.Join(pwd, tf)
	assert.True(t, IsDir(pwd))
	assert.False(t, IsDir(p))
}

//test file path is file
func TestIsFile(t *testing.T) {
	pwd, err := os.Getwd()
	assert.Nil(t, err)
	tf := createTestFile()
	defer cleanTestFile()

	p := path.Join(pwd, tf)
	assert.True(t, IsFile(p))
	assert.False(t, IsFile(pwd))
}

func TestFindDigits(t *testing.T) {
	bs := FindDigits("../testdata/consumer_config.yml")
	fmt.Println(string(bs))
}

func TestReadFile(t *testing.T) {
	inputFile, inputError := os.Open("../testdata/consumer_config.yml")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			return
		}
		fmt.Printf("The input was: %s", inputString)
	}
}

func TestWriteFile(t *testing.T) {
	inputFile := "../testdata/consumer_config.yml"
	outputFile := "products_copy.txt"
	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		// panic(err.Error())
	}
	fmt.Printf("%s\n", string(buf))
	err = ioutil.WriteFile(outputFile, buf, 0x644)
	if err != nil {
		panic(err.Error())
	}
}

func TestBufRead(t *testing.T) {
	inputFile, inputError := os.Open("../testdata/consumer_config.yml")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	buf := make([]byte, 1024)
	n, err := inputReader.Read(buf)
	fmt.Printf("read file count:%d,err:%s", n, err)
	fmt.Println(buf)
}

func TestReadFileByClown(t *testing.T) {
	file, err := os.Open("../testdata/products2.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var col1, col2, col3 []string
	for {
		var v1, v2, v3 string
		_, err := fmt.Fscanln(file, &v1, &v2, &v3)
		// scans until newline
		if err != nil {
			break
		}
		col1 = append(col1, v1)
		col2 = append(col2, v2)
		col3 = append(col3, v3)
	}

	fmt.Println(col1)
	fmt.Println(col2)
	fmt.Println(col3)
}

func TestBufWriter(t *testing.T) {
	outputFile, outputError := os.OpenFile("../testdata/products.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputString := "hello world!\n"

	for i := 0; i < 10; i++ {
		outputWriter.WriteString(outputString)
	}
	outputWriter.Flush()
}

func TestStdout(t *testing.T) {
	os.Stdout.WriteString("hello, world\n")
	f, _ := os.OpenFile("../testdata/test", os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	f.WriteString("hello, world in a file\n")
}

func TestCopyFile(t *testing.T) {
	dstfile, srcfile := "../testdata/test.txt", "../testdata/products.txt"
	ret, err := CopyFile(dstfile, srcfile)
	fmt.Printf("ret:%d,err:%s", ret, err)
}

func TestCatFile(t *testing.T) {
	CatFile("../testdata/products.txt", os.Stdout)
}

var filenameToContent map[string][]byte

var head = `// Copyright %s, Chef.  All rights reserved.
// https://%s
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Zeb灬D (1406721322@qq.com)`

var tail = `
> author: xxx
> link: xxx
> license: xxx
`

// /<root>/
//     |-- /dir1/
//     |-- /dir2/
//         |-- file5
//         |-- file6
//         |-- file7.txt
//         |-- file8.txt
//     |-- file1
//     |-- file2
//     |-- file3.txt
//     |-- file4.txt
func prepareTestFile() (string, error) {
	filenameToContent = make(map[string][]byte)

	root, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	if root[len(root)-1] != '/' {
		root = root + "/"
	}
	fmt.Println("root:" + root)

	if err = os.Mkdir(filepath.Join(root, "dir1"), 0755); err != nil {
		return "", err
	}
	if err = os.Mkdir(filepath.Join(root, "dir2"), 0755); err != nil {
		return "", err
	}

	filenameToContent[root+"file1"] = []byte("hello")
	filenameToContent[root+"file2"] = []byte("hello")
	filenameToContent[root+"file3.txt"] = []byte("hello")
	filenameToContent[root+"file4.txt"] = []byte("hello")
	filenameToContent[root+"dir2/file5"] = []byte("hello")
	filenameToContent[root+"dir2/file6"] = []byte("hello")
	filenameToContent[root+"dir2/file7.txt"] = []byte("hello")
	filenameToContent[root+"dir2/file8.txt"] = []byte("hello")

	for k, v := range filenameToContent {
		if err = ioutil.WriteFile(k, v, 0755); err != nil {
			return "", err
		}
	}

	return root, nil
}

func testWalk(t *testing.T, recursive bool, suffix string) {
	root, err := prepareTestFile()
	assert.Equal(t, nil, err)
	defer os.RemoveAll(root)

	err2 := Walk(root, recursive, suffix, func(path string, info os.FileInfo, content []byte, err error) []byte {
		t.Logf("%+v %+v %s", path, info.Name(), string(content))

		v := filenameToContent[path]
		assert.Equal(t, v, content)
		delete(filenameToContent, path)

		return content
	})
	assert.Equal(t, nil, err2)
}

func TestWalk(t *testing.T) {
	testWalk(t, true, "")
	assert.Equal(t, 0, len(filenameToContent))

	testWalk(t, false, "")
	assert.Equal(t, 4, len(filenameToContent))

	testWalk(t, true, ".txt")
	assert.Equal(t, 4, len(filenameToContent))

	testWalk(t, false, ".txt")
	assert.Equal(t, 6, len(filenameToContent))

	testWalk(t, false, ".notexist")
	assert.Equal(t, 8, len(filenameToContent))
}

func TestAddContent(t *testing.T) {
	root, err := prepareTestFile()
	assert.Equal(t, nil, err)
	defer os.RemoveAll(root)

	err2 := Walk(root, true, ".txt", func(path string, info os.FileInfo, content []byte, err error) []byte {
		lines := bytes.Split(content, []byte{'\n'})
		t.Logf("%+v %d", path, len(lines))

		v := filenameToContent[path]
		assert.Equal(t, v, content)
		delete(filenameToContent, path)

		return AddHeadContent(AddTailContent(content, []byte(tail)), []byte(head))
	})
	assert.Equal(t, nil, err2)

	err2 = Walk(root, true, "", func(path string, info os.FileInfo, content []byte, err error) []byte {
		t.Logf("%+v %+v %s", path, info.Name(), string(content))
		return nil
	})
	assert.Equal(t, nil, err2)
}

func TestDeleteLines(t *testing.T) {
	origin := `111
222
333
444
555`

	content := []byte(origin)
	lines := bytes.Split(content, []byte{'\n'})
	assert.Equal(t, 5, len(lines))

	var (
		res []byte
		err error
	)

	// 常规操作
	res, err = DeleteLines(content, LineRange{From: 1, To: 1})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`222
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: -5, To: -5})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`222
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: 2, To: 2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: -4, To: -4})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: 4, To: 4})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
555`), res)

	res, err = DeleteLines(content, LineRange{From: -2, To: -2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
555`), res)

	res, err = DeleteLines(content, LineRange{From: 5, To: 5})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
444`), res)

	res, err = DeleteLines(content, LineRange{From: -1, To: -1})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
444`), res)

	res, err = DeleteLines(content, LineRange{From: 1, To: 3})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`444
555`), res)

	res, err = DeleteLines(content, LineRange{From: -5, To: -3})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`444
555`), res)

	res, err = DeleteLines(content, LineRange{From: 3, To: 5})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222`), res)

	res, err = DeleteLines(content, LineRange{From: -3, To: -1})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222`), res)

	res, err = DeleteLines(content, LineRange{From: 2, To: 4})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
555`), res)

	res, err = DeleteLines(content, LineRange{From: -4, To: -2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
555`), res)

	// 非常规操作
	res, err = DeleteLines(content, LineRange{From: 4, To: 2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
555`), res)

	res, err = DeleteLines(content, LineRange{From: 0, To: 1})
	assert.Equal(t, ErrLineRange, err)

	res, err = DeleteLines(content, LineRange{From: 1, To: 0})
	assert.Equal(t, ErrLineRange, err)

	res, err = DeleteLines(content, LineRange{From: 10, To: 20})
	assert.Equal(t, ErrLineRange, err)
}
