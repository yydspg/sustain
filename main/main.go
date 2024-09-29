package main

import "github.com/yydspg/sustain"

func main() {
	context := sustain.Prepare()
	sustain.Run(context)
}
