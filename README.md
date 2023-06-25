# Go Fluent Bit GitHub Status Plugin

This is an example plugin used to showcase how to create a simple Golang input plugin for [Fluent Bit](https://fluentbit.io).

The article for this experiment can be found [here](https://niedbalski.dev/posts/writing-golang-fluent-bit-input-plugins/)

## Configuration

The following is an example Fluent Bit configuration (`fluent-bit.yaml`) that uses the Go Fluent Bit GitHub Status plugin:

```yaml
service:
  flush: 1
  log_level: debug
  plugins_file: /fluent-bit/etc/plugins.conf
  Parsers_file: /fluent-bit/etc/parsers.conf

pipeline:
  inputs:
    - Name: go-fluentbit-github-status

  filters:
    - Name: lua
      Match: '*'
      call: filter
      code: |
        function filter(tag, timestamp, record)
          if record.status and record.status.status and record.status.status.description == "All Systems Operational" then
            record = { status = "✅ All GitHub systems are operational" }
          else
            record = { status = "❌ Issues with GitHub (check: https://www.githubstatus.com/)" }
          end
          return 2, timestamp, record
        end

  outputs:
    - Name: slack
      Match: '*'
      webhook: https://hooks.slack.com/xxxx
```

Make sure to replace `https://hooks.slack.com/xxxx` with the actual webhook URL for your Slack integration.

Additionally, you need to create a `plugins.conf` file with the following content:

```ini
[PLUGINS]
    Path /fluent-bit/etc/go-fluentbit-github-status.so
```

To build the plugin, run the following command:

```shell
go build -trimpath -buildmode c-shared -o github_status.so .
```

Note for users of M1, compile as follows:

```shell
CGO_ENABLED=1 \
GOOS=linux \
GOARCH=amd64 \
CC="zig cc -target x86_64-linux-gnu -isystem /usr/include -L/usr/lib/x86_64-linux-gnu" \
CXX="zig c++ -target x86_64-linux-gnu -isystem /usr/include -L/usr/lib/x86_64-linux-gnu" \
go build -trimpath -buildmode c-shared -o github_status.so .
```

The resulting `github_status.so` file should be placed in the `/fluent-bit/etc` directory.

Finally, you can run the latest Fluent Bit release with the plugin loaded and executed using the following command:

```shell
docker run -v $(pwd)/github_status.so:/fluent-bit/etc/go-fluentbit-github-status.so -v $(pwd)/fluent-bit.yaml:/fluent-bit/etc/fluent-bit.yaml:ro -v $(pwd)/plugins.conf:/fluent-bit/etc/plugins.conf:ro fluent/fluent-bit:2.1.2 -c /fluent-bit/etc/fluent-bit.yaml
```

Make sure to adjust the volume mounts (`$(pwd)` represents the current directory) and the Fluent Bit image tag according to your environment.

With this configuration, the Go Fluent Bit GitHub Status plugin will check the GitHub status and rewrite the record's "status" field to indicate whether all systems are operational or if there are issues. The output can be sent to a Slack channel using the `slack` output plugin or modified to fit your specific use case.

Note: Ensure that you have the appropriate permissions and access to the necessary resources (e.g., GitHub API, Slack webhook) when using this plugin and configuring the output.
