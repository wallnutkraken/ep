package progress

import "fmt"

const (
	volumeBlockChar = "█"
	emptyBlock      = "░"
)

type Bar struct {
	Value     int
	Increment int
}

func (b Bar) Draw() {
	fmt.Printf("\rVolume: [%s]", b.getInner())
}

func (b Bar) getInner() string {
	var text string
	index := 0
	for index < b.Value {
		text += volumeBlockChar
		index += b.Increment
	}
	index = b.Value * b.Increment
	for index < 100 {
		text += emptyBlock
		index += b.Increment
	}
	return text
}
