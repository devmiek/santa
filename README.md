# :santa: santa
A simple, fast and extensible structured logging implemented in Go.

## Installation
Installing the Santa package is very simple, usually you only need to execute:

```shell
go get -u github.com/nobody-night/santa
```

## Getting Started
Santa is a simple, fast and extensible structured logging package implemented in Go. With Santa, you can quickly implement the printing function of structured log entries with multiple selectable log levels (severities) in the application, and output these log entries to one or more compatible specific storage devices.

Santa provides multiple types of loggers for applications, each of which provides efficient and easy-to-use APIs. If the provided loggers and APIs do not meet the workload requirements of the application, you can easily build custom loggers and APIs on top of these loggers.

Next, I will show you how to use the various loggers and commonly used APIs provided by Santa.

### Structured Logger
The first thing to show you is the structured logger. As the name implies, the structured logger provides an API for printing structured log entries. The payload encoding of a common structured log entry is JSON, which contains one or more fields. The following shows the JSON payload of a structured log entry:

```json
{
    "status": 500,
    "message": "Internal server error.",
    "requestId": "d4b02c71-7ca0-46c3-b18c-9e79f55df3c9"
}
```

Different from common single-line multi-field string log entries, structured log entries are easier to automatically analyze and query. The following example shows how to use the API provided by the structured logger to print out a simple structured log entry:

```go
// Create a default structured logger instance.
logger, _ := santa.NewStruct()

// The created logger instance must be explicitly closed.
defer logger.Close()

// Print out a structured log entry to standard output.
logger.Infos("Internal server error.",
    santa.Int("status", 500),
    santa.String("message", "Internal server error."),
    santa.String("requestId", "d4b02c71-7ca0-46c3-b18c-9e79f55df3c9"),
)
```

The following shows the actual output structured log entries, which have been formatted for human reading:

```json
{
    "timestamp": 1602324290419642100,
    "sourceLocation": {
        "file": "main.go",
        "line": 12,
        "function": "main.main"
    },
    "labels": null,
    "name": null,
    "level": "INFO",
    "message": {
        "text": "Internal server error.",
        "payload": {
            "status": 500,
            "message": "Internal server error.",
            "requestId": "d4b02c71-7ca0-46c3-b18c-9e79f55df3c9"
        }
    }
}
```

It is worth noting that the timestamp layout, key names and some empty key values shown above can all be customized.

### Template logger
The next thing to show you is the template logger. Unlike the structured logger, the template logger provides an API for printing template log entries, and its style is similar to the logger provided by the standard library. The payload encoding of common template log entries is a single-line string containing one or more field values separated by specific characters. The following shows the single-line string payload of the template log entry:

```text
500 "Internal server error." "d4b02c71-7ca0-46c3-b18c-9e79f55df3c9"
```

Generally, log entries encoded by single-line strings are easier for humans to read, but they are not conducive to automated analysis and query. The following example shows how to use the API provided by the template logger to print out a simple template log entry:

```go
// Create a default template logger instance.
logger, _ := santa.NewTemplate()

// The created logger instance must be explicitly closed.
defer logger.Close()

// Print out a template log entry to standard output.
logger.Infof("Status: %d, Message: %s, Request ID: %s",
    500, "Internal server error.", "d4b02c71-7ca0-46c3-b18c-9e79f55df3c9")
```

The following shows the actual output template log entries:

```text
2020-10-10T18:41:04.5541337+08:00 main.go:17 no-labels [INFO] "Status: 500,
Message: Internal server error., Request ID: d4b02c71-7ca0-46c3-b18c-9e79f55df3c9"
```

It is worth noting that the `NewTemplate` function will create a template logger instance with the default option values. Among them, the default log entry encoder is `EncoderStandard`. If you need a template logger to output structured log entries, you should use the `EncoderJSON` encoder. The following example shows how to create a template logger instance using the JSON encoder:

```go
// Create an option instance with default option values.
option := santa.NewTemplateOption()

// Change the encoder type to `EncoderJSON`.
option.Encoding.UseJSON()
// Optional: option.Encoding.UseJSONOption(...)
// Optional: option.UseEncoding(...)

// Use custom options to build a template logger instance.
logger, _ := option.Build()
```

Other types of loggers also support similar APIs. For details, please refer to the comment section of the `StandardOption` structure.

### Standard Logger
The last thing to show you is the standard logger. The standard logger provides an API for printing custom log entry message types, which means you can use the standard logger to print custom log entry message types, or you can build a custom logger based on the standard logger. It is worth noting that both the structured logger and the template logger are built on the standard logger.

The API provided by the standard logger accepts any log entry messages that have implemented the `Message` interface and the corresponding serialization interface, such as `StructMessage` and `TemplateMessage` structures. Santa provides a built-in pure string log entry message type named `StringMessage`. The following shows the single-line string payload of the string log entry:

```text
Internal server error.
```

The following example shows how to use the API provided by the standard logger to print out a simple string log entry:

```go
// Create a default template logger instance.
logger, _ := santa.NewStandard()

// The created logger instance must be explicitly closed.
defer logger.Close()

// Print out a string log entry to standard output.
logger.Info(santa.StringMessage("Internal server error."))
```

The following shows the actual output string log entries:

```text
2020-10-10T19:17:41.7185663+08:00 main.go:17 no-labels [INFO] "Internal server error."
```

### Outputting
Normally, the logger will output the log entries of `DEBUG`, `INFO` and `WARNING` levels to the standard output device (`os.Stdout`), and output the log entries of `ERROR` and `FATAL` levels to The standard error device (`os.Stderr`), which is controlled by the default value of the option.

Santa uses a synchronizer to output the log entries encoded by the encoder to a specific storage device (for example: local hard disk). Currently, the following types of synchronizers are provided:

- Standard Synchronizer
- File Synchronizer
- Network Synchronizer
- Discard Synchronizer

Among them, the standard synchronizer allows any structure that has implemented the `io.Writer` interface to be used as a specific storage device. For details, please refer to the comment section of the `StandardSyncer` structure.

#### Local File
The first thing to show you is how to use the file synchronizer to output log entries to a file located on the local hard disk:

```go
// Create an option instance with default option values.
option := santa.NewStructOption()

// Change the synchronizer type to `SyncerFile`.
option.Outputting.UseFile("./testing.log")
option.ErrorOutputting.UseFile("./testing_error.log")

// Use custom options to build a structured logger instance.
logger, _ := option.Build()
```

#### Network
The next thing I want to show you is how to use the network synchronizer to output log entries to TCP/IP or Unix Domain Socket streams:

```go
// Create an option instance with default option values.
option := santa.NewStructOption()

// Change the synchronizer type to `SyncerNetwork`.
option.Outputting.UseNetwork(santa.ProtocolTCP, "127.0.0.1:5000")
option.ErrorOutputting.UseNetwork(santa.ProtocolTCP, "127.0.0.1:5000")
// Optional: option.Outputting.UseNetwork(santa.ProtocolUnix, "/var/run/santa.sock")
// Optional: option.ErrorOutputting.UseNetwork(santa.ProtocolUnix, "/var/run/santa.sock")

// Use custom options to build a structured logger instance.
logger, _ := option.Build()
```

#### Discard
The last thing to show you is how to use the discard synchronizer to output log entries to the black hole:

```go
// Create an option instance with default option values.
option := santa.NewStructOption()

// Change the synchronizer type to `SyncerDiscard`.
option.Outputting.UseDiscard()
option.ErrorOutputting.UseDiscard()

// Use custom options to build a structured logger instance.
logger, _ := option.Build()
```

As the name implies, the discard synchronizer discards all output log entries, and no log entries are written to any specific storage device.

### Others
The logger also has many customizable options, including but not limited to: samplers, hooks, encoders, etc. For details, please refer to the comment section of the `StandardOption` structure.

## Performance
Santa provides efficient loggers and APIs, and uses many features to improve API performance, which means your application will not waste a lot of CPU time on printing out log entries. However, Santa pays more attention to the ease of use and extensible API, which requires the use of runtime features and maintaining some state, which requires some CPU time overhead.

I believe that every logging package has its value and advantages, so I will not compare Santa with other logging packages for the time being. The following data shows the Benchmark situation of the laboratory environment, for reference only, it is recommended that you evaluate the actual performance before releasing the application to the production environment.

It is worth noting that Benchmark was performed on a virtual machine instance in a laboratory environment. The following table shows the specifications of the virtual machine instance:

| Component | Parameters |
| :-------- | :--------- |
| Processor | AMD EPYC™ ROME, 4 Cores, 2.6GHz, AMD64 |
| Memory | 16 GBytes, DDR4 |
| Hard Disk | 1 TBytes, NVMe, SSD |
| System | CentOS 8.0, AMD64 |
| Runtime | Golang 1.15, AMD64 |

Each Benchmark will run 10 times and take the smallest value among all samples.

### Structured Logger
The first thing to show you is the Benchmark data of the structured logger. Benchmark uses the `santa.NewStructBenchmark` function to create a structured logger instance for testing, and then uses the `santa.(*StructLogger).Infos` function to print out structured log entries.

#### Complex Nested Fields
The following table shows the performance when using the API provided by the structured logger to output a structured log containing 10 complex fields:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 241 ns/op | 7 allocs/op |
| JSON | False | 681 ns/op | 7 allocs/op |
| Standard | True | 245 ns/op | 7 allocs/op |
| Standard | False | 749 ns/op | 7 allocs/op |

#### Simple Fields
The following table shows the performance of using the API provided by the structured logger to output a structured log containing 2 simple fields:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 88.8 ns/op | 1 allocs/op |
| JSON | False | 102 ns/op | 1 allocs/op |

#### Local File
The following table shows the performance when using the API provided by the structured logger to output a structured log containing 2 simple fields to a local file:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 89.7 ns/op | 1 allocs/op |
| JSON | False | 453 ns/op | 1 allocs/op |

Please note that the local hard disk may be a performance bottleneck, and the actual performance is affected by the performance of the local hard disk.

### Template Logger
The next thing to show you is the Benchmark data of the template logger. Benchmark uses the `santa.NewTemplateBenchmark` function to create a template logger instance for testing, and then uses the `santa.(*TemplateLogger).Infos` function to print out template log entries.

The following table shows the performance of using the API provided by the template logger to output a structured log containing parameters of 10 commonly used data types:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 82.6 ns/op | 1 allocs/op |
| JSON | False | 375 ns/op | 2 allocs/op |
| Standard | True | 84.5 ns/op | 1 allocs/op |
| Standard | False | 448 ns/op | 2 allocs/op |

### Standard Logger
The last thing I want to show you is the Benchmark data of the standard logger. Benchmark uses the `santa.NewStandardBenchmark` function to create a standard logger instance for testing, and then uses the `santa.(*StandardLogger).Infos` function to print out string log entries (`santa.StringMessage`).

The following table shows the performance when using the API provided by the standard logger to output pure string logs:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | ---: | :---------------: |
| JSON | True | 32.2 ns/op | 0 allocs/op |
| JSON | False | 51.4 ns/op | 0 allocs/op |
| Standard | True | 35.8 ns/op | 0 allocs/op |
| Standard | False | 118 ns/op | 0 allocs/op |

### Other
The standard encoder uses `time.RFC3339Nano` as the layout style of the timestamp by default, which means that additional time string formatting is required, causing additional CPU time and heap memory allocation overhead.

## Development Status: Alpha
At present, the preliminary development work has been basically completed, but some functions still need to be further tested and verified. When there is enough telemetry data to prove the stability of these functions, I will release the Beta version. It is worth noting that the API provided by the Alpha version may have errors, and one or more APIs may be changed or removed in the future. 

It is recommended that you only use the Alpha version in a test environment to avoid accidents.

## Contribute
I welcome contributions from any developers interested in Santa, but there is no detailed contribution guide for the time being. If you encounter a problem during use, please feel free to open a new issue on the issue tracker.

<hr>

Released under the [MIT License](LICENSE).
