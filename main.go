package main

import (
	// "fmt"

	"github.com/aneeshsunganahalli/GoGit/cmd"
	"github.com/aneeshsunganahalli/GoGit/internal"
)

func main() {

	cmd.Execute()
	// str := internal.ObjectHashing("Hello")
	// fmt.Println(str)
	// fmt.Println(str[:2])
	// internal.WriteObject("blob", "Hello Worlds")
	// internal.LoadIndex(".gogit/index.json")
	internal.RecursiveWalk("./cmd")
}
