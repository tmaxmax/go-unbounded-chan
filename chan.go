// Package unbounded provides utilities for creating
// channels with variable size buffers, which are helpful
// when you need sends that never block the sender's goroutine.
package unbounded

// NewChan creates a channel with a variable sized buffer. It returns
// two actual channels: the first one is for sending values, the second
// one is for receiving them. Closing the send channel will automatically
// close the receive channel after the receiver consumes all the remaining
// values in the buffer. The send channel will never block, as it is always
// ready for receiving. If the receive channel is ready, an attempt to send
// the received value from the send channel directly to it, without storing
// it in the intermediary buffer, will be made.
func NewChan() (chan<- interface{}, <-chan interface{}) {
	var buf []interface{}
	in := make(chan interface{})
	out := make(chan interface{})

	nextValue := func() interface{} {
		if len(buf) == 0 {
			return nil
		}
		return buf[0]
	}
	outChan := func() chan<- interface{} {
		if len(buf) == 0 {
			return nil
		}
		return out
	}

	go func() {
		defer close(out)

	loop:
		for {
		selector:
			select {
			case v, ok := <-in:
				if !ok {
					break loop
				}
				if len(buf) == 0 {
					select {
					case out <- v:
						break selector
					default:
					}
				}
				buf = append(buf, v)
			case outChan() <- nextValue():
				buf = buf[1:]
			}
		}
		for _, v := range buf {
			out <- v
		}
	}()

	return in, out
}
