package pierum

import (
	"log"
	"pie-rum-sdk/common"
	"sync"
	"time"
)

type IWrite[In, Out any] struct {
	Profile ISequence[In]
	Report  []*IDispatchResult
}

// autoWrite performs the dispatching of the profile events as per the desc and writes the metrics
func (r *PieRum[In, Out]) autoWrite(profile ISequence[In]) *IWrite[In, Out] {
	ctx := r.ctx

	log.Println("writing...")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		prf := r.Store.Registry[profile.Profile]
		var profiles, kits, services, dispatchers, events []string
		profiles = append(profiles, profile.Profile)

		for _, k := range prf.GetKits() {
			log.Printf("workin on kit %s", k.Name)
			if !k.Config.getActivate() {
				continue
			}

			for _, ser := range k.GetServices() {
				log.Printf("workin on service %s", ser.Name)
				if !ser.Config.getActivate() {
					continue
				}
				for _, dispatcher := range ser.Registry {
					log.Printf("workin on dispatcher %s", dispatcher.Name)

					if !ser.Config.getActivate() {
						continue
					}

					mx := profile.Input

					kits = append(kits, k.Name)
					services = append(services, ser.Name)
					dispatchers = append(kits, dispatcher.Name)

					if !r.settings.EnableMetricReport {
						dispatcher.normalCall(ctx, *mx)
						for _, res := range dispatcher.result {
							r.Store.AddResult(res)
						}
					} else {
						var errs = dispatcher.metricCall(ctx, *mx)

						for eventName, e := range errs {
							errKey := profile.Profile + "." + k.GetName() + "." + ser.GetName() + "." + eventName
							captured := e
							r.cheetahDetector.Publish(errKey, &captured)
						}

						for eventName := range dispatcher.metric {
							log.Printf("workin on event %s", eventName)
							inp := dispatcher.GetResults(eventName)
							if inp == nil {
								continue
							}
							events = append(kits, eventName)
							inp.MetaInfo = IDispatchResultMetaData{
								Profile:     profiles,
								Kits:        kits,
								Services:    services,
								Dispatchers: dispatchers,
								Events:      events,
							}
							inp.CreatedAt = common.FormatDateForClient(time.Now())
							r.Store.AddResult(inp)
						}
					}

					//r.Store.UpdateProfileSlateUsage(profile.Profile)
					//prf.UpdateKitSlateUsage(k.GetName())
					//k.UpdateServiceSlateUsage(ser.GetName())
					//ser.UpdateDispatcherSlateUsage(dispatcher.GetName())

				}
			}
		}
	}()

	wg.Wait()

	finalResult := r.Store
	res := finalResult.result
	x := &IResults{
		Resuts:  res,
		IsReady: true,
	}

	r.cheetah.Publish(profile.Profile, x)
	log.Println("final res: ", finalResult)
	log.Println("done writing...")

	return &IWrite[In, Out]{
		Report:  finalResult.result,
		Profile: profile,
	}
}
