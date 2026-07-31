// Package pierum implements the continous flow of concurrent funcs manager which can be regisreted via server can be called via clinet
// flow:
//
//												                Create Profile:
//											                            ↓
//											                    "profile" includes:
//											         kit: provides the descirption of the profile that includes:
//											             "services": contains descirption of the service which include:
//										                             "time-format" ->  deactivate the service, remove the service,activation time, retry call if the dispatch failed, on what time to invoke, duration of delay per dispatch,
//									                                 "dispatcher" -> controls the number of registered funcs
//								                                     "budget"     -> controls the mode budget
//													     "time-format": same as services but for the profile
//							                                  Call The Profile:
//							                                            ↓
//					                                        grpc accepts the call & triggers hub through channel
//				                                                     ↓
//			                                               hub performs as per the call:
//						                                 onPost-> fetches the service -> reads the format -> performs the write -> write publishes the work -> tickFetch fetches the result -> Paper publishes the result -> result is passed to the client
//	                                                  onDeactivate-> read desc ->  temporarly remove the service or profile
//	                                                  onActivate-> read desc ->  find the deactivate serivce or profile -> activates the service or profile
//	                                                  onRemove-> read desc ->  remove the service or profile
//
// [new] profiles, services and dispatchers now has their own configurations
// each configuration based on the toggle system
// toggle system: swap, stop , start
// how config works:
// rum
// |---- store
// ------------- profile
// --------------------- kit
// -------------------------- service
// ---------------------------------- dispatchers
// ----------------------------------------------- events
// store controls the profile config
// profile controls the kit config
// kit controls the service config
// service controls the dispatcher config
// dispatcher controls the events config
// documentation:-
// 1. active, inactive & swapped numbers and names of each generic configurations
package pierum

// events are maintain by the dispatcher
// dispatchers are maintain by the service
// services are maintain by the kit
// kits are maintain by the profile
// profiles are maintain by the store
// each has thier own configs

// GetMetadata returns the metadata of the stored
// on each request the metadata first rebuild & return new data
