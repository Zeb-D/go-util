package common

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"path"
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
