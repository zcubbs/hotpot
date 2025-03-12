package must

import (
	"fmt"
	"os"
)

func Succeed(err error) {
	if err != nil {
		GetTheHeckOut(err)
	}
}

func GetTheHeckOut(err error) {
	fmt.Println(err)
	os.Exit(1)
}
