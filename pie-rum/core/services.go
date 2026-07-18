// Package pierum...
// grpc erviers that just act as a gateway nothing else basically triggers the channels which than triggers the hub
package pierum

import (
	"context"
	"encoding/json"
	"log"
	rumrpc "pie-rum-sdk/misc/rum"
)

// POST publishes the paper
func (r *PieRum[In, Out]) POST(ctx context.Context, in *rumrpc.IPostRequest) (*rumrpc.IPostResponse, error) {
	log.Println("in push")
	var s = make([]ILink[In, Out], 0, len(in.Post))
	for _, x := range in.Post {
		var input In
		if err := json.Unmarshal(x.Profile.Input, &input); err != nil {
			continue
		}

		if err := r.handleActions(IConfigRequest{
			Action:  activateAction,
			Target:  depthFinder(1),
			Profile: x.Profile.Profile,
		}); err != nil {
			continue
		} else {
			s = append(s, ILink[In, Out]{
				Seq: ISequence[In]{
					Profile: x.Profile.Profile,
					Input:   &input,
				},
			})
			r.store.AddProfileUsage(x.Profile.Profile)
			r.post <- ILinks[In, Out]{Links: s, Clean: true}
		}

	}

	return &rumrpc.IPostResponse{Succeed: &rumrpc.ISucceed{Succeed: true}}, nil
}

func (r *PieRum[In, Out]) MonitorTag(ctx context.Context, in *rumrpc.IMonitorTagRequest) (*rumrpc.IMonitorTagResponse, error) {

	r.monitorTags <- in.Tag

	return &rumrpc.IMonitorTagResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (r *PieRum[In, Out]) Release(context.Context, *rumrpc.ReleaseRequest) (*rumrpc.ReleaseResponse, error) {

	r.release <- true

	return &rumrpc.ReleaseResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

// configuations ------------------

// profile

func (c *PieRum[In, Out]) ActivateProfile(ctx context.Context, in *rumrpc.IActivateProfileRequest) (*rumrpc.IActivateProfileResponse, error) {

	cnfs := in.Activate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
		})
	}

	c.activateProfile <- reqs

	return &rumrpc.IActivateProfileResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) DeactivateProfile(ctx context.Context, in *rumrpc.IDeactivateProfileRequest) (*rumrpc.IDeactivateProfileResponse, error) {
	cnfs := in.Deactivate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
		})
	}
	c.deactivateProfile <- reqs

	return &rumrpc.IDeactivateProfileResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) SwapProfile(ctx context.Context, in *rumrpc.ISwapProfileRequest) (*rumrpc.ISwapProfileResponse, error) {

	cnfs := in.Swap.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
		})
	}
	c.swapProfile <- reqs

	return &rumrpc.ISwapProfileResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

// end

// kit

func (c *PieRum[In, Out]) DeactivateKit(ctx context.Context, in *rumrpc.IDeactivateKitRequest) (*rumrpc.IDeactivateKitResponse, error) {

	cnfs := in.Deactivate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
		})
	}
	c.deactivateKit <- reqs

	return &rumrpc.IDeactivateKitResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) ActivateKit(ctx context.Context, in *rumrpc.IActivateKitRequest) (*rumrpc.IActivateKitResponse, error) {

	cnfs := in.Activate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
		})
	}
	c.activateKit <- reqs

	return &rumrpc.IActivateKitResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) SwapKit(ctx context.Context, in *rumrpc.ISwapKitRequest) (*rumrpc.ISwapKitResponse, error) {

	cnfs := in.Swap.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
		})
	}
	c.swapKit <- reqs

	return &rumrpc.ISwapKitResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

// end

// service

func (c *PieRum[In, Out]) DeactivateService(ctx context.Context, in *rumrpc.IDeactivateServiceRequest) (*rumrpc.IDeactivateServiceResponse, error) {

	cnfs := in.Deactivate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
			Service: r.Service,
		})
	}
	c.deactivateService <- reqs

	return &rumrpc.IDeactivateServiceResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) ActivateService(ctx context.Context, in *rumrpc.IActivateServiceRequest) (*rumrpc.IActivateServiceResponse, error) {
	cnfs := in.Activate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
			Service: r.Service,
		})
	}
	c.activateService <- reqs

	return &rumrpc.IActivateServiceResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) SwapService(ctx context.Context, in *rumrpc.ISwapServiceRequest) (*rumrpc.ISwapServiceResponse, error) {
	cnfs := in.Swap.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
			Service: r.Service,
		})
	}
	c.swapService <- reqs
	return &rumrpc.ISwapServiceResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

// end

// dispatcher

func (c *PieRum[In, Out]) DeactivateDispatcher(ctx context.Context, in *rumrpc.IDeactivateDispatcherRequest) (*rumrpc.IDeactivateDispatcherResponse, error) {
	cnfs := in.Deactivate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile:    r.Profile,
			Kit:        r.Kit,
			Service:    r.Service,
			Dispatcher: r.Dispatcher,
		})
	}
	c.deactivateDispatcher <- reqs
	return &rumrpc.IDeactivateDispatcherResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) ActivateDispatcher(ctx context.Context, in *rumrpc.IActivateDispatcherRequest) (*rumrpc.IActivateDispatcherResponse, error) {
	cnfs := in.Activate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile:    r.Profile,
			Kit:        r.Kit,
			Service:    r.Service,
			Dispatcher: r.Dispatcher,
		})
	}
	c.activateDispatcher <- reqs

	return &rumrpc.IActivateDispatcherResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) SwapDispatcher(ctx context.Context, in *rumrpc.ISwapDispatcherRequest) (*rumrpc.ISwapDispatcherResponse, error) {
	cnfs := in.Swap.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile:    r.Profile,
			Kit:        r.Kit,
			Service:    r.Service,
			Dispatcher: r.Dispatcher,
		})
	}
	c.swapDispatcher <- reqs

	return &rumrpc.ISwapDispatcherResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

// end

// event

func (c *PieRum[In, Out]) DeactivateEvent(ctx context.Context, in *rumrpc.IDeactivateEventRequest) (*rumrpc.IDeactivateEventResponse, error) {
	cnfs := in.Deactivate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
			Service: r.Service,
			Event:   r.Event,
		})
	}
	c.deactivateEvent <- reqs

	return &rumrpc.IDeactivateEventResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) ActivateEvent(ctx context.Context, in *rumrpc.IActivateEventRequest) (*rumrpc.IActivateEventResponse, error) {
	cnfs := in.Activate.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
			Service: r.Service,
			Event:   r.Event,
		})
	}
	c.activateEvent <- reqs
	return &rumrpc.IActivateEventResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

func (c *PieRum[In, Out]) SwapEvent(ctx context.Context, in *rumrpc.ISwapEventRequest) (*rumrpc.ISwapEventResponse, error) {
	cnfs := in.Swap.Config
	reqs := make([]IConfigRequest, 0, len(cnfs))

	for _, r := range cnfs {
		reqs = append(reqs, IConfigRequest{
			Profile: r.Profile,
			Kit:     r.Kit,
			Service: r.Service,
			Event:   r.Event,
		})
	}
	c.swapEvent <- reqs
	return &rumrpc.ISwapEventResponse{
		Succeed: &rumrpc.ISucceed{Succeed: true},
	}, nil
}

// end

// end -------------

// metadata

func (c *PieRum[In, Out]) GetProfileMetadata(ctx context.Context, in *rumrpc.IProfileMetadataRequest) (*rumrpc.IProfileMetadataResponse, error) {

	packets := make([]*rumrpc.IPackage, 0, len(in.Name))

	for _, r := range in.Name {
		packet := c.GetStore().GetProfileMetadata(r.Profile).JSON()

		packets = append(packets, &rumrpc.IPackage{
			Package: packet,
			Name:    r,
		})
	}
	return &rumrpc.IProfileMetadataResponse{
		Packet: packets,
	}, nil
}

func (c *PieRum[In, Out]) GetKitMetadata(ctx context.Context, in *rumrpc.IKitMetadataRequest) (*rumrpc.IKitMetadataResponse, error) {
	packets := make([]*rumrpc.IPackage, 0, len(in.Name))
	for _, r := range in.Name {
		packet := c.GetProfiles()[r.Profile].GetKitMetadata(r.Kit).JSON()

		packets = append(packets, &rumrpc.IPackage{
			Package: packet,
			Name:    r,
		})
	}
	return &rumrpc.IKitMetadataResponse{
		Packet: packets,
	}, nil
}

func (c *PieRum[In, Out]) GetServiceMetadata(ctx context.Context, in *rumrpc.IServiceMetadataRequest) (*rumrpc.IServiceMetadataResponse, error) {
	packets := make([]*rumrpc.IPackage, 0, len(in.Name))
	for _, r := range in.Name {
		packet := c.GetProfiles()[r.Kit].registry[r.Kit].GetServiceMetadata(r.Service).JSON()

		packets = append(packets, &rumrpc.IPackage{
			Package: packet,
			Name:    r,
		})
	}
	return &rumrpc.IServiceMetadataResponse{
		Packet: packets,
	}, nil
}

func (c *PieRum[In, Out]) GetDispatcherMetadata(ctx context.Context, in *rumrpc.IDispatcherMetadataRequest) (*rumrpc.IDispatcherMetadataResponse, error) {
	packets := make([]*rumrpc.IPackage, 0, len(in.Name))
	for _, r := range in.Name {
		packet := c.GetProfiles()[r.Kit].registry[r.Service].registry[r.Dispatcher].GetDispatcherMetadata(r.Dispatcher).JSON()
		packets = append(packets, &rumrpc.IPackage{
			Package: packet,
			Name:    r,
		})
	}
	return &rumrpc.IDispatcherMetadataResponse{
		Packet: packets,
	}, nil
}

func (c *PieRum[In, Out]) GetEventMetadata(ctx context.Context, in *rumrpc.IEventMetadataRequest) (*rumrpc.IEventMetadataResponse, error) {
	packets := make([]*rumrpc.IPackage, 0, len(in.Name))
	for _, r := range in.Name {
		packet := c.GetProfiles()[r.Kit].registry[r.Service].registry[r.Dispatcher].registry[r.Event].GetEventMetadata(r.Event).JSON()
		packets = append(packets, &rumrpc.IPackage{
			Package: packet,
			Name:    r,
		})
	}
	return &rumrpc.IEventMetadataResponse{
		Packet: packets,
	}, nil
}

func (c *PieRum[In, Out]) PullProfileMetadata(ctx context.Context, in *rumrpc.IPullProfileMetadataRequest) (*rumrpc.IPullProfileMetadataResponse, error) {
	pack := c.GetStore().GetProfilesMetadata().JSON()
	resp := &rumrpc.IPullProfileMetadataResponse{
		Packet: &rumrpc.IPackage{
			Package: pack,
		},
	}

	return resp, nil
}

func (c *PieRum[In, Out]) PullKitMetadata(ctx context.Context, in *rumrpc.IPullKitMetadataRequest) (*rumrpc.IPullKitMetadataResponse, error) {
	pack := c.GetProfiles()[in.Name.Kit].GetKitsMetadata().JSON()
	resp := &rumrpc.IPullKitMetadataResponse{
		Packet: &rumrpc.IPackage{
			Package: pack,
			Name:    in.Name,
		},
	}
	return resp, nil
}

func (c *PieRum[In, Out]) PullServiceMetadata(ctx context.Context, in *rumrpc.IPullServiceMetadataRequest) (*rumrpc.IPullServiceMetadataResponse, error) {
	pack := c.GetProfiles()[in.Name.Kit].registry[in.Name.Service].GetServicesMetadata().JSON()
	resp := &rumrpc.IPullServiceMetadataResponse{
		Packet: &rumrpc.IPackage{
			Package: pack,
			Name:    in.Name,
		},
	}
	return resp, nil
}

func (c *PieRum[In, Out]) PullDispatcherMetadata(ctx context.Context, in *rumrpc.IPullDispatcherMetadataRequest) (*rumrpc.IPullDispatcherMetadataResponse, error) {
	pack := c.GetProfiles()[in.Name.Kit].registry[in.Name.Service].registry[in.Name.Dispatcher].GetDispatchersMetadata().JSON()
	resp := &rumrpc.IPullDispatcherMetadataResponse{
		Packet: &rumrpc.IPackage{
			Package: pack,
			Name:    in.Name,
		},
	}
	return resp, nil
}

func (c *PieRum[In, Out]) PullEventMetadata(ctx context.Context, in *rumrpc.IPullEventMetadataRequest) (*rumrpc.IPullEventMetadataResponse, error) {
	pack := c.GetProfiles()[in.Name.Kit].registry[in.Name.Service].registry[in.Name.Dispatcher].registry[in.Name.Event].GetEventsMetadata().JSON()
	resp := &rumrpc.IPullEventMetadataResponse{
		Packet: &rumrpc.IPackage{
			Package: pack,
			Name:    in.Name,
		},
	}
	return resp, nil
}

// end ------------
