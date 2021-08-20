package unbounded

import (
	"fmt"
	"time"
)

func ExampleNewChan() {
	in, out := NewChan()

	go func() {
		s := "hello world"

		for _, c := range s {
			in <- c
		}
		close(in)

		fmt.Println("Sending done!")
	}()

	for c := range out {
		time.Sleep(time.Millisecond)
		fmt.Printf("%c", c)
	}

	fmt.Println("\nReceiving done!")

	// Output:
	// Sending done!
	// hello world
	// Receiving done!
}
