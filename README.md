# go-log4g

`go-log4g` provides Log4j-style configuration and pattern layouts for
Go's standard `log/slog` logging facade.

It configures `slog`; it does not replace it. Applications can use
standard `slog` directly or the optional `log4g` facade for `{}`
parameterized messages.

## Configuration

Create `config/log4g.yaml`:

``` yaml
appenders:
  console:
    type: console
    target: stdout
    layout:
      type: pattern
      pattern: "%d{yyyy-MM-dd HH:mm:ss.SSS}{UTC} %-5p %c.%M:%L - %m%n"

root:
  level: warn
  appenders:
    - console

loggers:
  playground/internal/app:
    level: debug

  github.com/go-beans/go/ioc:
    level: debug
```

This configuration writes to `stdout`, sets the default level to `WARN`,
and enables `DEBUG` for the two configured logger hierarchies.

Example output:

``` text
2026-08-13 12:34:56.789 INFO  playground/internal/app/Service1.AfterPropertiesSet:23 - Service initialized
```

Logger configuration is hierarchical. For example:

``` text
playground/internal/app
```

also matches:

``` text
playground/internal/app/Service1
playground/internal/app/service/UserService
```

## Initialization

Import `core` early in the application so logging is configured before
other application packages initialize:

``` go
import (
    _ "github.com/go-log4g/core"

    "log/slog"
)
```

Use standard `slog` normally:

``` go
slog.Info("Application started")
slog.Debug("Loading user", "userId", 123)
```

Or use the optional `log4g` facade:

``` go
import "github.com/go-log4g/core/log4g"

log4g.Info("User {} authenticated", 123)
log4g.Debug("Loaded {} records in {}", count, elapsed)
```

## Pattern layout

The configured pattern:

``` text
%d{yyyy-MM-dd HH:mm:ss.SSS}{UTC} %-5p %c.%M:%L - %m%n
```

produces output such as:

``` text
2026-08-13 12:34:56.789 INFO  playground/internal/app/Service1.AfterPropertiesSet:23 - Service initialized
```

### Supported patterns

``` text
%d{pattern}          Date/time
%d{pattern}{UTC}     Date/time in UTC

%p                   Log level
%level               Same as %p

%c                   Logger name
%logger              Same as %c

%M                   Method/function name
%method              Same as %M

%F                   Source file
%file                Same as %F

%L                   Source line
%line                Same as %L

%m                   Message
%msg                 Same as %m
%message             Same as %m

%X{key}              MDC value
%X                   All MDC values

%n                   Platform newline
%%                   Literal %
```

Examples:

``` text
%d{yyyy-MM-dd HH:mm:ss.SSS}{UTC}
→ 2026-08-13 12:34:56.789

%p
→ INFO

%c
→ playground/internal/app/Service1

%M
→ AfterPropertiesSet

%F:%L
→ Service1.go:23

%m
→ Service initialized

%X{requestId}
→ 8f24...

%X
→ {requestId=8f24..., userId=123}
```

## Width and alignment

A minimum width can be specified before a converter:

``` text
%5p       minimum width 5, right aligned
%-5p      minimum width 5, left aligned
```

Example:

``` text
%-5p
```

produces aligned levels:

``` text
DEBUG
INFO 
WARN 
ERROR
```

Width is also useful for optional MDC values:

``` text
[%36X{requestId}]
```

With a request ID:

``` text
[550e8400-e29b-41d4-a716-446655440000]
```

Without a request ID:

``` text
[                                    ]
```

Use `-` for left alignment:

``` text
[%-36X{requestId}]
```

## Logger name precision

Given:

``` text
playground/internal/app/Service1
```

the following patterns produce:

``` text
%c        → playground/internal/app/Service1
%c{1}     → Service1
%c{2}     → app/Service1
%c{3}     → internal/app/Service1

%c{-1}    → internal/app/Service1
%c{-2}    → app/Service1
```

Positive precision keeps rightmost components. Negative precision
removes components from the left.

## Logger name abbreviation

Given:

``` text
org/apache/commons/test/Foo
```

simple abbreviation:

``` text
%c{1.}    → o/a/c/t/Foo
%c{2.}    → or/ap/co/te/Foo
```

Explicit component rules:

``` text
%c{1.1.1.*}
→ o/a/c/test/Foo
```

Here the first three components are shortened to one character and `*`
leaves the remaining components unchanged.

Dynamic abbreviation:

``` text
%c{1.2.*}
→ o/a/c/test/Foo
```

Here leading components are shortened to one character while the last
two components are preserved.

Another example:

``` text
%c{1.3.*}
→ o/a/commons/test/Foo
```

This style is useful for Go logger names because a long module prefix
can be abbreviated while the local package and type remain readable.

## MDC

MDC associates logging values with a `context.Context`.

``` go
ctx = mdc.Put(ctx, "requestId", requestId)
ctx = mdc.Put(ctx, "userId", userId)

log4g.InfoContext(ctx, "Processing request")
```

Use one MDC value in a pattern:

``` text
%X{requestId}
```

or all values:

``` text
%X
```

`mdc.Put` follows Go context semantics and returns a derived context:

``` go
ctx = mdc.Put(ctx, "requestId", requestId)
```

Each derived MDC context contains the complete MDC snapshot, so
subsequent `Put` operations retain previously added values.
