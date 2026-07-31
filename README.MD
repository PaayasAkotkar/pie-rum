## PIE RUM


+ ██████╗░██╗███████╗  ██████╗░██╗░░░██╗███╗░░░███╗
+ ██╔══██╗██║██╔════╝  ██╔══██╗██║░░░██║████╗░████║
+ ██████╔╝██║█████╗░░  ██████╔╝██║░░░██║██╔████╔██║
+ ██╔═══╝░██║██╔══╝░░  ██╔══██╗██║░░░██║██║╚██╔╝██║
+ ██║░░░░░██║███████╗  ██║░░██║╚██████╔╝██║░╚═╝░██║
+ ╚═╝░░░░░╚═╝╚══════╝  ╚═╝░░╚═╝░╚═════╝░╚═╝░░░░░╚═╝

<img src="/assets/logo/logo.png" alt="pie-rum-logo"/>

### Polyglot Integrated Environment: Run Your Models

A developer-centric framework AI orchestration, and work engine. It bridges the gap between the implemented system layers and developer’s communication. Reducing the complexity of tools or communication apps and providing one complete codex for the organization.
Supports all the fields in the list of computing or programming.

---

## How PIE RUM SDK works?

1. **core layer server** -> focuses on the orchestration -> store [ profiles [M:N] [ kits [M:N] [ services [M:N] [ dispatchers [M:N] [ events ] ] ] ] ]
   * each mapping field is many-to-many resulting into structured configuration wiring storage.
2. **the developer must set the default profile** for search safety 
3. **each request being made** to the implemented pie-rum layer is being filtered to dedicated search params namely for the store [….prof…]..]
4. **pie-rum runs all the dedicated kit** of that service & of that dispatchers events in the monitor method if config; one-by-one and pushes the result either failure or succeed asap to the core pie-rum monitor method without needing to wait to complete all the methods
5. **the result is stored** in the pie-rum server itself the developer decides what to do with it

---

## Why use PIE RUM SDK?

* 1) focus directly on building the core product layer & connect the dots together
* 2) achieve meta-info of those mapping fields
* 3) config the mapping fields not only a single field but many at once
* 4) structure configuration system to pull meta-info or config the mapping fields

---

## Glimpse

```golang
// filling the pie-rum store
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
			}) 
		})
	}).
		Build()

      // fetching the results
      pie.GoMonitor(ctx, func(result *pierum.IResults) {
			if result.IsReady {
				log.Println("res: ", result.Resuts)
				for _, r := range result.Resuts {
					var rs structure.Res
					var in structure.Req
					log.Println("input: ", string(r.Input))
					log.Println("output: ", string(r.Output))
					log.Println("dog-report", string(r.DogReport))
					log.Println("meta-report", r.MetaInfo)

					if err := json.Unmarshal(r.Input, &in); err != nil {
						log.Println("Error unmarshaling input:", err)
						continue
					}
					if err := json.Unmarshal(r.Output, &rs); err != nil {
						log.Println("Error unmarshaling output:", err)
						continue
					}

					if in.Student.ID == nil {
						log.Println("Student ID is nil")
						continue
					}

					if rs.Err != nil {
						log.Println("Error from AI model:", rs.Err)
						x := rs.Err.Error()
						errorReply := &model.OnChessCoachReply{
							Information: &model.ChessCoachReply{
								Desc: &x,
							},
						}
						app.cheetah.Publish(*in.Student.ID, errorReply)
						continue
					}

					if rs.Coach != nil {
						app.cheetah.Publish(*in.Student.ID, rs.Coach)
					} else {
						log.Println("Coach result is nil")
					}
				}
			}
		})
```

> [!NOTE]
> 
> 1) **Testing for the meta-leaderboard is left** > Do test it and report how you guys liked it what more is needed
> 2) The **Dog Monitor** interface might get replaced via promethus
> If you feel like not to remove or we can create more using that let me know that too
> 3) The dependency injection might get replaced by uber-fx 
> If you feel like not t remove or we can create more using that let me know that too
> 3) I am only the solo developer who improved the codes, many things are yet to be tested by my side and I am not fully sure whether it is ready or not cause I kept on focusing on writing the core than testing out things 
> Please do test the thing and do comment down so that i can know whats missing
> More work is still coming by my side regarding monitoring funcs to work more roboustly than using the wait-group stuffs
> 4) most exampels are written by AI do make sure to report if you find bugs
> 5) Most of the things are kinda imroved by calude so yeah shotout to claude
---

## Quick-Start
1. Pick any example from using the example package and run it to understand how things working
2. simply run the complie.bat for windows & compile.sh for linux & mac

## New 🔥 🌟
1. playground: you can now use the chat-app to work around frontend and backend
2. plugin: allows you to get the reports in your own [finally 🤗] 🔥
3. refresh: keeps at least one field in the [m:n] config active
4. only supports direct config part

## Changes 🪛
1. removed the meta-data stuffs
2. removed the grpc active|deative|inactive|swap services

## Upcoming 🔥 🚋
1. promethus addition
2. uber-fx