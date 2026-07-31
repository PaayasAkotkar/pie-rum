package injection

import (
	injection "pie-rum-sdk/di/core"
	"reflect"
)

type ServiceRequest struct {
	Type    reflect.Type
	Factory injection.Factory
}
