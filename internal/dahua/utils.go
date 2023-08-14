package dahua

func drainChannel(c <-chan struct{}) {
	for {
		select {
		case <-c:
		default:
			return
		}
	}
}
