package rcrypt

import "os"

var (
	ENV           = ""                          // your api key
	MODEL         = "googleai/gemini-2.5-flash" // or anyother model
	DBURL         = ""
	PIERUMADRESS  = "localhost:1200"
	PIERUMNETWORK = "tcp"
)

func init() {
	if e := os.Getenv("ENV"); e != "" {
		ENV = e
	}
	if m := os.Getenv("MODEL"); m != "" {
		MODEL = m
	}
	if d := os.Getenv("DBURL"); d != "" {
		DBURL = d
	}
}
