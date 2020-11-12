package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var lenSet int
var set []string
var fdLock chan struct{}

func main() {
	var count uint64
	var root string
	{
		var setStr string
		var fdCap uint
		flag.Uint64Var(&count, "count", 3, "depth of folders to create")
		flag.StringVar(&setStr, "set", "a,b,c", "csv set to use at each step")
		flag.StringVar(&root, "root", "t", "first directory to be created")
		flag.UintVar(&fdCap, "fdcap", 1024, "maximum number of concurrently opened files")
		flag.Parse()
		set = strings.Split(setStr, ",")
		lenSet = len(set)
		fdLock = make(chan struct{}, int(fdCap))
		for fdCap > 0 {
			fdLock <- struct{}{}
			fdCap--
		}
	}
	uint64LenSet := uint64(lenSet)
	var total uint64 = uint64LenSet
	for i := count; i > 0; i-- {
		total = total * uint64LenSet
	}
	var totalGoroutine uint64 = 1
	var previous uint64 = 1
	for i := count; i > 0; i-- {
		previous = previous * uint64LenSet
		totalGoroutine += previous
	}
	if totalGoroutine > math.MaxInt32 {
		fmt.Println(strconv.FormatUint(totalGoroutine, 10) + " total goroutines is too big, max number is " + strconv.FormatInt(int64(math.MaxInt32), 10) + ".")
		os.Exit(1)
	}
	if lenSet == 0 {
		fmt.Println("Can't use an empty set.")
		os.Exit(1)
	}
	fmt.Println("Creating " + strconv.FormatUint(total, 10) + " files.")
	wg.Add(int(totalGoroutine))
	err := os.Mkdir(root, 0777)
	if err != nil {
		panic(err)
	}
	do(uint(count), root)
	wg.Wait()
	fmt.Println("Done !")
}

func do(count uint, path string) {
	i := lenSet
	pretotal := path + "/"
	if count == 0 {
		<-fdLock
		for i > 0 {
			i--
			total := pretotal + set[i]
			err := ioutil.WriteFile(total, []byte(total), 0644)
			if err != nil {
				panic(err)
			}
		}
		fdLock <- struct{}{}
	} else {
		for i > 1 {
			i--
			total := pretotal + set[i]
			err := os.Mkdir(total, 0777)
			if err != nil {
				panic(err)
			}
			go do(count-1, total)
		}
		total := pretotal + set[0]
		err := os.Mkdir(total, 0777)
		if err != nil {
			panic(err)
		}
		do(count-1, total)
	}
	wg.Done()
}
