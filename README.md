The Snowplow Tracking CLI is a command-line app (written in Golang) to make it easy to send an event to Snowplow from the command line. Use this to embed Snowplow tracking into your shell scripts.

### Usage

The app is called `snowplowtrk`.

The command line interface is as follows:

```bash
snowplowtrk {{COLLECTOR_DOMAIN}} --appid={{APP_ID}} --method=[POST|GET] --sdjson={{SELF_DESC_JSON}}
```
    
or:

```bash
snowplowtrk {{COLLECTOR_DOMAIN}} --appid={{APP_ID}} --method=[POST|GET] --schema={{SCHEMA_URI}} --json={{JSON}}
```

where:

* `{{COLLECTOR_DOMAIN}}` is the domain for your Snowplow collector, e.g. `snowplow-collector.acme.com`
* `--appid` is optional (not sent if not set)
* `--method` is optional. Defaults to `GET`
* `--sdjson` is a self-describing JSON of the standard form `{ "schema": "iglu:xxx", "data": { ... } }`
* `--schema` is a schema URI, most likely of the form `iglu:xxx`
* `--json` is a (non-self-describing) JSON, of the form `{ ... }`

The idea here is that you can either send in a self-describing JSON, or pass in the constituent parts (i.e. a regular JSON plus a schema URI) and the Snowplow Tracking CLI will construct the final self-describing JSON for you.

Return codes:

* 0 if the Snowplow collector responded with an OK status (2xx or 3xx)
* 4 if the Snowplow collector responded with a 4xxx status
* 5 if the Snowplow collector responded with a 5xx status
* 1 for any other error

Under the hood:

* There is no buffering - each event is sent as an individual payload whether `GET` or `POST`
* The Snowplow Tracking CLI will exit once the Snowplow collector has responded
* The app uses the [Snowplow Golang Tracker](https://github.com/snowplow/snowplow-golang-tracker)

## Copyright and license

The Snowplow Tracking CLI is copyright 2016 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0] [license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[license]: http://www.apache.org/licenses/LICENSE-2.0
