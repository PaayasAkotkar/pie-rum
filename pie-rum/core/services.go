// Package pierum...
// grpc erviers that just act as a gateway nothing else basically triggers the channels which than triggers the hub
package pierum

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	rumrpc "pie-rum-sdk/misc/rum"
)

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
			//r.Store.AddProfileUsage(x.Profile.Profile)
			r.post <- ILinks[In, Out]{Links: s, Clean: true}
		}

	}

	return &rumrpc.IPostResponse{Succeed: &rumrpc.ISucceed{Succeed: true}}, nil
}

func (r *PieRum[In, Out]) GetDoc(ctx context.Context, in *rumrpc.IDocRequest) (*rumrpc.IDocResponse, error) {
	doc := r.JSON()
	if doc == nil {
		return &rumrpc.IDocResponse{
			Succeed: false,
			Doc:     "",
			Error:   "[something went wrong while packing the doc]",
		}, fmt.Errorf("[something went wrong while packing the doc]")
	}
	return &rumrpc.IDocResponse{
		Succeed: true,
		Doc:     string(doc),
		Error:   "",
	}, nil
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

// end

func (r *PieRum[In, Out]) UpdateProfileConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if err := r.handleActions(IConfigRequest{
		Action:  configAction,
		Target:  depthOne,
		Config:  rpcToPIERUM(in.Config),
		Profile: in.ProfileName,
	}); err != nil {
		return &rumrpc.IConfigResponse{
			Error:   err.Error(),
			Succeed: false,
		}, nil
	}
	return &rumrpc.IConfigResponse{
		Error:   "",
		Succeed: true,
	}, nil
}
func (r *PieRum[In, Out]) UpdateKitConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	log.Println("in: ", in.ProfileName, in.Config, in.KitName)
	if err := r.handleActions(IConfigRequest{
		Action:  configAction,
		Target:  depthTwo,
		Config:  rpcToPIERUM(in.Config),
		Profile: in.ProfileName,
		Kit:     in.KitName,
	}); err != nil {
		return &rumrpc.IConfigResponse{
			Error:   err.Error(),
			Succeed: false,
		}, nil
	}
	return &rumrpc.IConfigResponse{
		Error:   "",
		Succeed: true,
	}, nil
}
func (r *PieRum[In, Out]) UpdateServiceConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if err := r.handleActions(IConfigRequest{
		Action:  configAction,
		Target:  depthThree,
		Config:  rpcToPIERUM(in.Config),
		Profile: in.ProfileName,
		Kit:     in.KitName,
		Service: in.ServiceName,
	}); err != nil {
		return &rumrpc.IConfigResponse{
			Error:   err.Error(),
			Succeed: false,
		}, nil
	}
	return &rumrpc.IConfigResponse{
		Error:   "",
		Succeed: true,
	}, nil
}
func (r *PieRum[In, Out]) UpdateDispatcherConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if err := r.handleActions(IConfigRequest{
		Action:     configAction,
		Target:     depthFour,
		Config:     rpcToPIERUM(in.Config),
		Profile:    in.ProfileName,
		Kit:        in.KitName,
		Service:    in.ServiceName,
		Dispatcher: in.DispatcherName,
	}); err != nil {
		return &rumrpc.IConfigResponse{
			Error:   err.Error(),
			Succeed: false,
		}, nil
	}
	return &rumrpc.IConfigResponse{
		Error:   "",
		Succeed: true,
	}, nil
}
func (r *PieRum[In, Out]) UpdateEventConfig(ctx context.Context, in *rumrpc.IConfigRequest) (*rumrpc.IConfigResponse, error) {
	if err := r.handleActions(IConfigRequest{
		Action:     configAction,
		Target:     depthFive,
		Config:     rpcToPIERUM(in.Config),
		Profile:    in.ProfileName,
		Kit:        in.KitName,
		Service:    in.ServiceName,
		Dispatcher: in.DispatcherName,
		Event:      in.EventName,
	}); err != nil {
		return &rumrpc.IConfigResponse{
			Error:   err.Error(),
			Succeed: false,
		}, nil
	}
	return &rumrpc.IConfigResponse{
		Error:   "",
		Succeed: true,
	}, nil
}

// end -------------
