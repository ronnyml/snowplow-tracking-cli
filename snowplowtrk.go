package main 

import(
	"fmt"
	"os"
	"github.com/urfave/cli"
	//"gopkg.in/snowplow/snowplow-golang-tracker.v1/tracker"
)

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
    	appid := c.String("appid")
    	method := c.String("method")
    	sdjson := c.String("sdjson")
    	schema := c.String("schema")
    	json := c.String("json")

        fmt.Println("APP ID:", appid)
        fmt.Println("Method:", method)
        fmt.Println("Self-Describing JSON:", sdjson)
        fmt.Println("Schema:", schema)
        fmt.Println("JSON:", json)
        return nil
    }

    app.Run(os.Args)
}

