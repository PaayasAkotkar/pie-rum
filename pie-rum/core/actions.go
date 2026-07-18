package pierum

const (
	activateAction = "activate:"
	deactiveAction = "deactive:"
	swapAction     = "swap:"
	depthOne       = "profile"
	depthTwo       = "kit"
	depthThree     = "service"
	depthFour      = "dispatcher"
	depthFive      = "event"
)

func depthFinder(depth int) string {
	switch depth {
	case 1:
		return depthOne
	case 2:
		return depthTwo
	case 3:
		return depthThree
	case 4:
		return depthFour
	case 5:
		return depthFive
	}
	return depthOne
}
