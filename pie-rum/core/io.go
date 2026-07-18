package pierum

import (
	"log"
	"sync"
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

		prf := r.store.registry[profile.Profile]

		for _, k := range prf.GetKits() {
			for _, ser := range k.GetServices() {
				for _, dispatcher := range ser.registry {
					mx := profile.Input

					if !r.settings.EnableMetricReport {
						dispatcher.normalCall(ctx, *mx)
						for _, res := range dispatcher.result {
							r.store.AddResult(res)
						}
					} else {
						var errs = dispatcher.metricCall(ctx, *mx)

						for eventName, e := range errs {
							errKey := profile.Profile + "." + k.GetName() + "." + ser.GetName() + "." + eventName
							captured := e
							r.cheetahDetector.Publish(errKey, &captured)
						}

						for eventName := range dispatcher.metric {
							inp := dispatcher.GetResults(eventName)
							r.store.AddResult(inp)
						}
					}

					r.store.UpdateProfileSlateUsage(profile.Profile)
					prf.UpdateKitSlateUsage(k.GetName())
					k.UpdateServiceSlateUsage(ser.GetName())
					ser.UpdateDispatcherSlateUsage(dispatcher.GetName())
				}
			}
		}
	}()

	wg.Wait()

	finalResult := r.store
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
