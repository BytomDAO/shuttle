package main

import (
	"flag"
	"fmt"
)

func main() {
	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")
	//用程序中已有的参数来声明一个标志也是可以的。注意在标志声明函数中需要使用该参数的指针。
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	//所有标志都声明完成以后，调用 flag.Parse() 来执行命令行解析。
	//这是关键步骤，不要忘记！！！
	flag.Parse()
	//这里我们将仅输出解析的选项以及后面的位置参数。
	//注意，我们需要使用类似 *wordPtr 这样的语法来对指针解引用，从而得到选项的实际值。
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}
