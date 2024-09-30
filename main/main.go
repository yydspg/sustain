package main

import "github.com/yydspg/sustain"
import _ "github.com/yydspg/sustain/example"

func main() {
	context := sustain.Prepare()
	sustain.Run(context)
}
