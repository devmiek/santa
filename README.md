# :santa: santa
A simple, fast and extensible structured logging implemented in Go.

## Getting Started
With Santa, you can easily implement a simple, fast, and easily expandable structured logging function in your application, without having to pay attention to log output, encoding, and storage.

Santa currently supports a variety of loggers, including but not limited to template loggers, structured loggers, etc. If your application needs to customize the message structure of log entries, you can also write customizable loggers and log message implementations based on standard loggers.

In the following, I will briefly introduce you how to use several common loggers provided by Santa to record logs in your application:

### Structured Logger
In a modern production environment, the most frequently used log entry structure may be JSON-encoded log entry data. This can be done easily in Santa using a structured logger:

```go
// Create a structured logger instance using default
// optional values.
logger, err := santa.NewStruct()

if err != nil {
    println(err)
    return
}

// Loggers should be explicitly closed when they
// are no longer in use.
defer logger.Close()

logger.Infos("Hello World!",
    santa.String("name", "santa"),
    santa.Int("age", 10),
)
```

In the above sample code, a structured logger with default options is used to print a structured log message with INFO log level to the standard output device (`os.Stdout`, `os.Stderr`). If you format the output manually, it should look like this:

```json
{
    "timestamp": 1597325688546993100,
    "sourceLocation": {
        "file": "main.go",
        "line": 18,
        "function": "main.main"
    },
    "logName": null,
    "level": "INFO",
    "message": "Hello World!",
    "payload": {
        "name": "santa",
        "age": 10
    }
}
```

Of course, the timestamp layout style and key names used in JSON encoding can be customized.

### Template Logger
If your application does not need to record one or more log fields, or does not need to record structured log entries, then a template logger may be a good choice.

Compared with the structured logger, the template logger provides an easier-to-use string formatting API for applications, just like using the `fmt.Sprintf` function. It is worth noting that because the template logger needs to format log messages according to template strings and parameters, its log entry output performance is lower than other types of loggers, because the template logger still uses the `fmt.Sprintf` function to Format the log message.

```go
// Create a template logger instance with default
// optional values.
logger, err := santa.NewTemplate()

if err != nil {
    println(err)
    return
}

// Loggers should be explicitly closed when they
// are no longer in use.
defer logger.Close()

logger.Infof("My name is %s and my age is %d.", "santa", 10)
```

Unlike the structured logger, the template logger uses a standard encoder to encode each log entry by default. If your application needs to output JSON-encoded log entries, you can specify the JSON encoder when building the template logger instance:

#### Use Builder Style

```go
logger, err := santa.NewTemplateOption().
    UseEncoding(
        santa.NewEncodingOption().
            UseJSON(),
    ).Build()
```

#### Use Option Style

```go
option := santa.NewTemplateOption()
option.Encoding.Kind = santa.EncoderJSON

logger, err := option.Build()
```

The sample code above uses the template logger built with default options to print out a template log entry with the log level INFO to the standard output device (`os.Stdout`, `os.Stderr`). If everything is normal, you should see something similar to the following:

```text
2020-08-13T21:56:30.0719939+08:00 main.go:18 [INFO] My name is santa and my age is 10.
```

### Standard Logger
If your application requires a custom log message structure or only string log messages, you can use the standard logger.

Unlike structured loggers and template loggers, standard loggers provide log message output APIs for applications that accept any log message values that have implemented the `santa.Message` interface, which means you can easily customize one or more log message structures.

```go
// Create a standard logger instance with default
// optional values.
logger, err := santa.NewStandard()

if err != nil {
    println(err)
    return
}

// Loggers should be explicitly closed when they
// are no longer in use.
defer logger.Close()

logger.Info(santa.StringMessage("Hello World!"))
```

In the above sample code, a standard logger instance is constructed using the default optional values, and then a log entry of only string messages with a log level of INFO is output to the standard output device (`os.Stdout`, `os.Stderr`). If everything is normal, you should see something similar to the following:

```text
2020-08-14T16:09:58.9404613+08:00 main.go:43 [INFO] Hello World!
```

In fact, structured loggers and template loggers are implemented based on standard loggers. The log message output API provided by the structured logger uses `santa.StructMessage` as the log message structure; the log message output API provided by the template logger uses `santa.TemplateMessage` as the log message structure.

If the application requires high log entry output performance, it may be a good choice to use the `santa.StringMessage` log message structure as the log message value of the log message output API provided by the standard logger. Under normal circumstances, the structured logger is also very fast, but requires some coding overhead for the structured fields.

It is worth noting that if your application needs to customize the log message structure, the customized log message structure also needs to implement the formatter interface of the corresponding encoder, otherwise the encoder does not know how to format the custom log message structure. Among them, the formatter interface of the standard encoder is `santa.StandardFormatter`, and the format interface of the JSON encoder is `santa.JSONFormatter`.

Similarly, your application can also customize one or more encoders, as long as these encoders have implemented the `santa.Encoder` interface. If the encoder needs to support multiple log message structures, you need to define a formatter interface and let all supported log message structures implement it.

### Others
Santa’s design focuses on scalability, and performance is second. This means that many of Santa’s functions can be easily customized, including but not limited to loggers, log messages, samplers, encoders, and synchronizers. For details, please refer to the comment section of each function in the Santa source code.

## Performance
It is worth noting that Santa focuses on scalability when designing, performance is secondary, but performance is still concerned.

Santa strives to be closer to the actual production environment in the benchmark performance test process, because the performance in the production environment is more meaningful.

The benchmark performance test was conducted on a VM instance, which is equipped with a 4-core AMD EPYC™ ROME processor and 16 GB of DDR4 memory. The processor clocked at 2.6GHz and adopts AMD64 architecture. The benchmark performance test is performed by the benchmark program running on the VM instance using all processor cores, and the benchmark program uses the Golang 1.15 runtime. The benchmark program is run 10 times in total, and the final result of each indicator is the average of all benchmark samples of the indicator. The test method is as follows:

### Structured Logger
For the structured logger, the benchmark program uses the `santa.NewStructBenchmark` function to build an instance of the structured logger for benchmark testing.

The benchmark program will continuously call the `santa.(*StructLogger).Infos` function to print out different log messages, each of which contains a different description text and 10 fields (including 5 complex fields). The benchmark test results are as follows:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 241 ns/op | 7 allocs/op |
| JSON | False | 681 ns/op | 8 allocs/op |
| Standard | True | 245 ns/op | 7 allocs/op |
| Standard | False | 749 ns/op | 8 allocs/op |

### Template Logger
For the template logger, the benchmark program uses the `santa.NewTemplateBenchmark` function to build an instance of the template logger for benchmark testing.

The benchmark program will continuously call the `santa.(*TemplateLogger).Infof` function to print out different log messages, including a different template string and 10 commonly used template parameters. The benchmark test results are as follows:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 82.6 ns/op | 1 allocs/op |
| JSON | False | 375 ns/op | 3 allocs/op |
| Standard | True | 84.5 ns/op | 1 allocs/op |
| Standard | False | 448 ns/op | 3 allocs/op |

### Standard Logger
For the standard logger, the benchmark program uses `santa.NewStandardBenchmark` to build an instance of the standard logger for benchmark testing.

The benchmark program will continuously call the `santa.(*StandardLogger).Info` function to print out different `santa.StringMessage` log messages. The benchmark test results are as follows:

| Encoder | Sampling | Time | Objects Allocated |
| :------ | :------: | :--: | :---------------: |
| JSON | True | 35.9 ns/op | 0 allocs/op |
| JSON | False | 69.2 ns/op | 1 allocs/op |
| Standard | True | 35.8 ns/op | 0 allocs/op |
| Standard | False | 118 ns/op | 1 allocs/op |
