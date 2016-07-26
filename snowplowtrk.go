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

import (
        "encoding/json"
        "fmt"
        "github.com/urfave/cli"
        gt "gopkg.in/snowplow/snowplow-golang-tracker.v1/tracker"
        "os"
        "strconv"                                                                                                                                 
)                                                                                                                                                 
                                                                                                                                                  
type SelfDescJson struct {                                                                                                                        
        Schema string                 `json:"schema"`                                                                                             
        Data   map[string]interface{} `json:"data"`                                                                                               
}                                                                                                                                                 

var sdj *gt.SelfDescribingJson

func main() {
        app := cli.NewApp()
        app.Name = "snowplowtrk"
        app.Usage = "The Snowplow Tracking CLI"
        app.Author = "Snowplow Analytics"
        app.Version = "0.1.0"
        app.Flags = []cli.Flag{
                cli.StringFlag{
                        Name:  "appid, id",
                        Usage: "APP ID (Optional)",
                },
                cli.StringFlag{
                        Name:  "method, m",
                        Usage: "Method[POST|GET] (Optional)",
                        Value: "GET",
                },
                cli.StringFlag{
                        Name:  "sdjson, sdj",
                        Usage: "self-describing JSON of the standard form { 'schema': 'iglu:xxx', 'data': { ... } }",
                },
                cli.StringFlag{
                        Name:  "schema, s",
                        Usage: "schema URI, most likely of the form iglu:xxx",
                },
                cli.StringFlag{
                        Name:  "json, j",
                        Usage: "(non-self-describing) JSON, of the form { ... }",
                },
        }

        app.Action = func(c *cli.Context) error {
                collector := c.Args().Get(0)
                appid := c.String("appid")
                method := c.String("method")
                sdjson := c.String("sdjson")
                schema := c.String("schema")
                jsonData := c.String("json")

                if sdjson == "" && schema == "" && jsonData == "" {
                        panic("FATAL: A --sdjson or a --schema URI plus a --json need to be specified.")
                } else if sdjson == "" && schema != "" && jsonData == "" {
                        panic("FATAL: A --json need to be specified.")
                } else if sdjson == "" && schema == "" && jsonData != "" {
                        panic("FATAL: A --schema URI need to be specified.")
                } else {
                        if sdjson != "" {
                                res := SelfDescJson{}
                                json.Unmarshal([]byte(sdjson), &res)
                                data, err := json.Marshal(res.Data)
                                if err != nil {
                                        panic(err)
                                }

                                fmt.Println("Schema:", res.Schema)
                                fmt.Println("Data:", string(data))
                                sdj = gt.InitSelfDescribingJson(res.Schema, string(data))
                        } else if schema != "" && jsonData != "" {
                                b := []byte(jsonData)
                                var jsonDataMap map[string]interface{}
                                err := json.Unmarshal(b, &jsonDataMap)

                                fmt.Println(jsonDataMap)
                                fmt.Println(err)
                                sdj = gt.InitSelfDescribingJson(schema, jsonDataMap)
                        }

                        fmt.Println("Collector:", collector)
                        fmt.Println("APP ID:", appid)
                        fmt.Println("Method:", method)
                        fmt.Println("Self-Describing JSON:", sdjson)
                        fmt.Println("Schema:", schema)
                        fmt.Println("JSON:", jsonData)

                        trackerChan := make(chan int, 1)
                        callback := func(s []gt.CallbackResult, f []gt.CallbackResult) {
                                fmt.Println("Callback executing")
                                status := 0

                                if len(s) == 1 {
                                        status = s[0].Status
                                }

                                if len(f) == 1 {
                                        status = f[0].Status
                                }
                                trackerChan <- status
                        }

                        tracker := initTracker(collector, appid, method, callback)
                        trackSelfDescribingEvent(tracker)

                        statusCode := <-trackerChan
                        fmt.Println("StatusCode: " + strconv.Itoa(statusCode))

                }

                return nil
        }

        app.Run(os.Args)
}

func initTracker(collector string, appid string, requestType string, callback func(successCount []gt.CallbackResult, failureCount []gt.CallbackResult)) *gt.Tracker {
        subject := gt.InitSubject()
        emitter := gt.InitEmitter(gt.RequireCollectorUri(collector),
                gt.OptionCallback(callback),
                gt.OptionRequestType(requestType),
                gt.OptionDbName("/home/vagrant/test.db"),
        )
        tracker := gt.InitTracker(
                gt.RequireEmitter(emitter),
                gt.OptionSubject(subject),
                gt.OptionAppId(appid),
        )

        return tracker
}

func getReturnCode(statusCode int) int {
        var returnCode int

        switch statusCode {
        case 2, 3:
                returnCode = 0
        case 4:
                returnCode = 4
        case 5:
                returnCode = 5
        default:
                returnCode = 1

        }

        fmt.Println("ReturnCode:", strconv.Itoa(returnCode))
        return returnCode
}

func trackSelfDescribingEvent(tracker *gt.Tracker) {
        tracker.TrackSelfDescribingEvent(gt.SelfDescribingEvent{
                Event: sdj,
        })
}
