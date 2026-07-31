package pierum

import "fmt"

func activationError(reason string) error {
	return fmt.Errorf("[activate-read-fail] reason: %s", reason)
}

func deactivationError(reason string) error {
	return fmt.Errorf("[deactive-read-fail] reason: %s", reason)
}

func swapError(reason string) error {
	return fmt.Errorf("[swap-read-fail] reason: %s", reason)
}

func removeError(reason string) error {
	return fmt.Errorf("[remove-read-fai] reason: %s", reason)
}
