package pierum

import (
	"fmt"
)

// handleActions typically related to searching for the activate, deactive & swap stuff
func (r *PieRum[In, Out]) handleActions(token IConfigRequest) error {
	action := token.Action
	if len(action) > 0 && action[len(action)-1] != ':' {
		action += ":"
	}
	key := action + token.Target
	entry, ok := r.actions[key]
	if !ok {
		return fmt.Errorf("unknown action %s", key)
	}
	res, err := r.walk(token, entry.depth)
	if err != nil {
		return err
	}
	return entry.action(res, token)
}

func (r *PieRum[In, Out]) handleDispatch(seq ISequence[In]) {
	r.autoWrite(seq)
}

func (r *PieRum[In, Out]) onPost(seq ISequence[In]) {

	r.handleDispatch(seq)
}

func (r*PieRum[In, Out]) walk(token IConfigRequest, depth int) (*Resolved[In, Out], error) {
	res := &Resolved[In, Out]{}

	if !r.store.IsProfileActive(token.Profile) {
		return nil, activationError(fmt.Sprintf("profile %s not active", token.Profile))
	}
	if depth == 1 {
		return res, nil
	}

	res.prf = r.store.registry[token.Profile]
	if !res.prf.IsKitActive(token.Kit) {
		return nil, activationError(fmt.Sprintf("kit %s not active", token.Kit))
	}
	if depth == 2 {
		return res, nil
	}

	res.kit = res.prf.GetKit(token.Kit)
	if !res.kit.IsServiceActive(token.Service) {
		return nil, activationError(fmt.Sprintf("service %s not active", token.Service))
	}
	if depth == 3 {
		return res, nil
	}

	res.svc = res.kit.GetService(token.Service)
	if !res.svc.IsDispatcherActive(token.Dispatcher) {
		return nil, activationError(fmt.Sprintf("dispatcher %s not active", token.Dispatcher))
	}
	if depth == 4 {
		return res, nil
	}

	res.dt = res.svc.GetDispatcher(token.Dispatcher)
	if !res.dt.IsEventActive(token.Event) {
		return nil, activationError(fmt.Sprintf("event %s not active", token.Event))
	}
	return res, nil
}
