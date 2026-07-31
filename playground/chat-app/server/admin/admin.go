package admin

import (
	"app/server/rag/flow"
	nvideamodel "app/server/rag/nvidea/model"
	"app/server/server/graph/model"
	"context"
	"encoding/json"
	"log"
	pierum "pie-rum-sdk/pie-rum/core"
	rumsdk "pie-rum-sdk/pie-rum/sdk"
	"time"
)

// JSONError wraps error for proper JSON marshaling/unmarshaling
type JSONError struct {
	Msg string
}

func (e *JSONError) Error() string {
	return e.Msg
}

func (e *JSONError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	return json.Marshal(e.Msg)
}

func (e *JSONError) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var msg string
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	e.Msg = msg
	return nil
}

// Req think it as the input to be recieve
type Req struct {
	Student *model.ChessStudentRequest
	//Store   *store.Store
}

func (r *Req) Pack() []byte {
	p, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return nil
	}
	return p
}

// Res think it as the output given as per the req
type Res struct {
	Coach *model.OnChessCoachReply
	Err   *JSONError
}

// Init returns the flow of the model placement
func Init(ctx context.Context) *rumsdk.PieRum[Req, Res] {

	store := rumsdk.NewStore[Req, Res](ctx)

	// Configure settings for AI model operations with longer timeout
	dispatcherSettings := pierum.Settings{
		Base:      2 * time.Minute, // 2 minutes for AI model operations
		SleepTime: 100 * time.Millisecond,
	}

	store.AddProfile("nvidea", func(p *rumsdk.Profile[Req, Res]) {
		p.SetConfig(&pierum.IConfig{
			Activate: true,
		})
		p.AddKit("nvidea-kit", func(kit *rumsdk.Kit[Req, Res]) {
			kit.SetConfig(&pierum.IConfig{
				Activate: true,
			})
			kit.AddService("nvidea-service", func(service *rumsdk.Service[Req, Res]) {
				service.AddDispatcher("nvidea-dispatcher", func(d *rumsdk.Dispatcher[Req, Res]) {
					d.SetSettings(dispatcherSettings)
					d.AddEvent("nvidea-event", func(event *rumsdk.Event[Req, Res]) {
						event.RegisterFunc(func(ctx context.Context, in Req) (Res, error) {
							c, err := flow.Flow(ctx, nvideamodel.GK(ctx), in.Student)
							var jsonErr *JSONError
							if err != nil {
								log.Println(err)
								jsonErr = &JSONError{Msg: err.Error()}
							}
							return Res{
								Coach: c,
								Err:   jsonErr,
							}, nil
						})
					})
				})
			}).SetConfig(&pierum.IConfig{
				Activate: false,
			})
		}).AddKit("nvidea-omni", func(kit *rumsdk.Kit[Req, Res]) {
			kit.AddService("nvidea-omni-service", func(service *rumsdk.Service[Req, Res]) {
				service.AddDispatcher("nvidea-omni-dispatcher", func(d *rumsdk.Dispatcher[Req, Res]) {
					d.SetSettings(dispatcherSettings)
					d.AddEvent("nvidea-omni-event", func(event *rumsdk.Event[Req, Res]) {
						event.RegisterFunc(func(ctx context.Context, in Req) (Res, error) {
							c, err := flow.Flow(ctx, nvideamodel.GKOmni(ctx), in.Student)
							var jsonErr *JSONError
							if err != nil {
								log.Println(err)
								jsonErr = &JSONError{Msg: err.Error()}
							}
							return Res{
								Coach: c,
								Err:   jsonErr,
							}, nil
						})
					})
				})
			})
			kit.SetConfig(&pierum.IConfig{
				Activate: true,
			}) // keep it disable
		})
	}).
		Build()
	pie := rumsdk.New(ctx, store)
	pie.SetPlugin(&pierum.IPlugin{
		Key: "http://127.0.0.1:1234/pie-rum-plugin",
		On:  true,
		Org: "lexuxa.pvt.ltd",
	})
	return pie
}
