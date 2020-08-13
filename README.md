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

In the above sample code, a structured logger with default options is used to print a structured log message with INFO log level to the standard output device (os.Stdout, os.Stderr). If you format the output manually, it should look like this:

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

The sample code above uses the template logger built with default options to print out a template log entry with the log level INFO to the standard output device (os.Stdout, os.Stderr). If everything is normal, you should see something similar to the following:

```json
2020-08-13T21:56:30.0719939+08:00 main.go:43 [INFO] My name is santa and my age is 10.
```
