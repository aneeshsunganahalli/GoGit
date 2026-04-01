package main

import (
	"fmt"

	"github.com/aneeshsunganahalli/GoGit/cmd"
	"github.com/aneeshsunganahalli/GoGit/internal"
)

func main() {

	cmd.Execute()
	var b internal.Blob
	str := internal.ObjectHashing(b, "Hello")
	fmt.Println(str)
}
