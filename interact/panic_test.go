package interact

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func PanicRecover() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Panicing %s\r\n", e)
		}
	}()
	panic(errors.New("has wrong"))
	fmt.Printf("After bad call\r\n") // <-- wordt niet bereikt
}

func TestPanicRecover(t *testing.T) {
	fmt.Println("start")
	PanicRecover()
	fmt.Println("end")
	fmt.Println(os.Environ())
}
