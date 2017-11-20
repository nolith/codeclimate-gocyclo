# Code Climate Gocyclo Engine

`codeclimate-gocyclo` is a Code Climate engine that wraps [Gocyclo](https://github.com/fzipp/gocyclo). You can run it on your command line using the Code Climate CLI, or on our hosted analysis platform.

Gocyclo calculates cyclomatic complexities of functions in Go source code.

### Installation

1. If you haven't already, [install the Code Climate CLI](https://github.com/codeclimate/codeclimate).
2. Run `codeclimate engines:enable gocyclo`. This command both installs the engine and enables it in your `.codeclimate.yml` file.
3. You're ready to analyze! Browse into your project's folder and run `codeclimate analyze`.

### Configuration

Like the `gocyclo` binary, you can configure the minimum allowed cyclomatic complexity
The default value is `9`: you can configure your own threshold in your `.codeclimate.yml`:

```yaml
version: "2"
plugins:
  gocyclo:
    enabled: true
    config:
      over: 5
```

### Building

```console
docker build -t codeclimate/codeclimate-gocyclo .
```
