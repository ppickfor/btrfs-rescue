package main

import (
	"fmt"

	"github.com/monnand/GoLLRB/llrb"
)

func main() {
	//	tree := llrb.New(lessInt)
	tree := llrb.New()
	tree.ReplaceOrInsert(llrb.Int(3))
	tree.ReplaceOrInsert(llrb.Int(4))
	tree.ReplaceOrInsert(llrb.Int(1))
	tree.ReplaceOrInsert(llrb.Int(2))

	tree.DeleteMin()
	tree.Delete(llrb.Int(4))
	fmt.Println("Ascend")
	tree.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%v\n", i)
		return true
	})
	fmt.Println("Descend")
	tree.DescendLessOrEqual(llrb.Inf(1), func(i llrb.Item) bool {
		fmt.Printf("%v\n", i)
		return true
	})
	fmt.Println("Has")
	if tree.Has(llrb.Int(3)) {
		i := tree.Get(llrb.Int(3))
		fmt.Printf("%v\n", i)
	}

}
