//
// Copyright (c) 2016 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package main 

import(
	"fmt"
	"os"
	"github.com/urfave/cli"
	gt "gopkg.in/snowplow/snowplow-golang-tracker.v1/tracker"
)

var tracker gt.Tracker
var sdjson *gt.SelfDescribingJson
var sdj []gt.SelfDescribingJson

func main() {
    app := cli.NewApp()
    app.Name = "snowplowtrk"
    app.Usage = "The Snowplow Tracking CLI"
    app.Author = "Snowplow Analytics"
    app.Version = "0.1"
    app.Flags = []cli.Flag {
	    cli.StringFlag{
	      Name: "appid, id",
	      Usage: "APP ID (Optional)",
	    },
	    cli.StringFlag{
	      Name: "method, m",
	      Usage: "Method[POST|GET] (Optional)",
	      Value: "GET",
	    },
	    cli.StringFlag{
	      Name: "sdjson, sdj",
	      Usage: "self-describing JSON of the standard form { 'schema': 'iglu:xxx', 'data': { ... } }",
	    },
	    cli.StringFlag{
	      Name: "schema, s",
	      Usage: "schema URI, most likely of the form iglu:xxx",
	    },
	    cli.StringFlag{
	      Name: "json, j",
	      Usage: "(non-self-describing) JSON, of the form { ... }",
	    },
	}

    app.Action = func(c *cli.Context) error {
    	collector := c.Args().Get(0)
        appid := c.String("appid")
    	method := c.String("method")
    	sdjson := c.String("sdjson")
    	schema := c.String("schema")
    	json := c.String("json")

        fmt.Println("Collector:", collector)
        fmt.Println("APP ID:", appid)
        fmt.Println("Method:", method)
        fmt.Println("Self-Describing JSON:", sdjson)
        fmt.Println("Schema:", schema)
        fmt.Println("JSON:", json)

        initTracker(collector, appid)
        return nil
    }

    app.Run(os.Args)
}

func initTracker(collector string, appid string) {
	subject := gt.InitSubject()
    emitter := gt.InitEmitter(gt.RequireCollectorUri(collector))
    tracker := gt.InitTracker(
	    gt.RequireEmitter(emitter),
		gt.OptionSubject(subject),
		gt.OptionAppId(appid),
	)
}

func trackPageView(pageurl string) {
	tracker.TrackPageView(
        gt.PageViewEvent{ 
        PageUrl: gt.NewString(pageurl),
        Contexts: sdj,
      },
    )
}

func trackScreenView(screen_id string) {
	tracker.TrackScreenView(gt.ScreenViewEvent{ 
        Id: gt.NewString(screen_id),
    })
}

func trackSelfDescribingEvent() {
	tracker.TrackSelfDescribingEvent(gt.SelfDescribingEvent{ 
        Event: sdjson,
    })
}

func trackStructEvent(category string, action string, label string) {
	tracker.TrackStructEvent(gt.StructuredEvent{ 
	    Category: gt.NewString(category),
		Action: gt.NewString(action),
		Label: gt.NewString(label),
	})
}

func trackTiming(category string, variable string, timing int64) {
    tracker.TrackTiming(gt.TimingEvent{ 
	    Category: gt.NewString(category),
	    Variable: gt.NewString(variable),
	    Timing: gt.NewInt64(timing),
	})
}

