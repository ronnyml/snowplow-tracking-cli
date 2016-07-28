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
        "github.com/stretchr/testify/assert"
        "github.com/urfave/cli"
        gt "gopkg.in/snowplow/snowplow-golang-tracker.v1/tracker"
        "testing"
)

func TestInitTracker(t *testing.T) {
        assert := assert.New(t)
        callback := func(s []gt.CallbackResult, f []gt.CallbackResult) {}
        tracker := initTracker("com.acme", "myapp", "POST", callback)
        assert.NotNil(tracker)
}

func TestGetReturnCode(t *testing.T) {
        assert := assert.New(t)
        result := getReturnCode(200)
        assert.NotNil(result)
        assert.Equal(0, result)
}

func TestApp(t *testing.T) {
        var collector, appid, method, sdjson, schema, jsonData string

        app := cli.NewApp()
        app.Flags = []cli.Flag{
                cli.StringFlag{Name: "appid, id"},
                cli.StringFlag{Name: "method, m", Value: "GET"},
                cli.StringFlag{Name: "sdjson, sdj"},
                cli.StringFlag{Name: "schema, s"},
                cli.StringFlag{Name: "json, j"},
        }

        app.Action = func(c *cli.Context) error {
                collector = c.Args().Get(0)
                appid = c.String("appid")
                method = c.String("method")
                sdjson = c.String("sdjson")
                schema = c.String("schema")
                jsonData = c.String("json")
                return nil
        }

        err := app.Run([]string{"", "--appid", "myappid", "--method", "POST", "--schema", "iglu:com.snowplowanalytics.snowplow/event/jsonschema/1-0-0", "--json", "{\"name\":\"foo\"}", "snowplow-collector.acme.com"})
        if err != nil {
                t.Fatalf("Run error: %s", err)
        }

        assert.Equal(t, collector, "snowplow-collector.acme.com")
        assert.Equal(t, appid, "myappid")
        assert.Equal(t, method, "POST")
        assert.Equal(t, schema, "iglu:com.snowplowanalytics.snowplow/event/jsonschema/1-0-0")
        assert.Equal(t, jsonData, "{\"name\":\"foo\"}")
}
