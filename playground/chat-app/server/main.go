// Package main implements the entry point of this package
package main

import (
	"app/server/server"
	"log"
)

// Fairly Speaking except pie-rum all codes was wrote in
// 10-02-2026 and at point of time I was just learning to
// typically understand what's goingin now
// I had no idea about the pg-vectro so i wrote that using teh claude
// So luckily I had this code
// Now the time being we mature today is 20-07-206 I know its a 5 month
// difference but one thing is that I never visited this 10-02-2026 repo i mean this one
// was so lucky for me to find that something to be wrote purely clear to
// explain how pie-rum works
// what changes did i made?
// the changes made was created a resuable flow while back it was only coded for gemini
// removed pg-vector & added marlin the milvus typically best to monitor things i beleive 😅 i mean like its claude generated but it works as i dont want
// to spend much more time to write my own milvus serach engine as i have been working on other projects too 😎
// -----------------
// added the pie-rum while rest is pure naked as before
// this will tell you how powerful is pie-rum casue
// it typically need nothing apart from some extra stuff whcih is not the
// demand of the pie-rum only the demand to acutally make the code more clean
// so yeah like I started working on PIE RUM since 15-04-2026
func main() {
	log.SetFlags(log.Lshortfile)
	server.Serve()
}
