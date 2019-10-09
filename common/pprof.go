package common

import (
	"os"
	"runtime/pprof"
)

func StartCPUProfile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	pprof.StartCPUProfile(f)
	//pprof.WriteHeapProfile(f)
	return nil
}
