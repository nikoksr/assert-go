<div align="center">

&nbsp;
<h1>assert-go</h1>
<p><i>Zero-dependency, idiomatic Go assertion library focused on crystal-clear failure messages and thoughtful source context.</i></p>

&nbsp;

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/assert-go)
</div>

&nbsp;

## About

- 🔍 Crystal-clear failure messages with contextual values
- 📚 Rich source context showing the exact failure location
- 🛠 Tiny and free of dependencies (~100 lines of Go)
- 💡 Elegant, idiomatic Go API
- 🎯 Two-tier assertion system with build tag support
- ⚙️ Configurable source context behavior

Inspired by [Tiger Style](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety).

## Installation

```bash
go get github.com/nikoksr/assert-go
```

## Usage

### Basic Usage

```go
import "github.com/nikoksr/assert-go"

func PaymentProcessing() {
    payment := processPayment(PaymentRequest{
        Amount:      99.99,
        CustomerID: "cust_123",
        Currency:   "USD",
    })
    
    // Assert payment was processed successfully
    assert.Assert(payment.Status == "completed", "payment should be completed",
        // Optionally, add context to the panic
        "payment_id", payment.ID,
        "status", payment.Status,
        "amount", payment.Amount,
        "error", payment.Error,
        "timestamp", payment.ProcessedAt,
    )
}
```

On failure, you get:

```
Assertion failed at payment_test.go:43
Message: payment should be completed

Relevant values:
  [payment_id]: "pmt_789"
  [status]: "failed"
  [amount]: 99.99
  [error]: "insufficient_funds"
  [timestamp]: "2024-12-06T15:04:05Z"

Source context:
   37 |     payment := processPayment(PaymentRequest{
   38 |         Amount:      99.99,
   39 |         CustomerID: "cust_123",
   40 |         Currency:   "USD",
   41 |     })
   42 |
→  43 |     assert.Assert(payment.Status == "completed", "payment should be completed",
   44 |         "payment_id", payment.ID,
   45 |         "status", payment.Status,
   46 |         "amount", payment.Amount,
   47 |         "error", payment.Error,
   48 |         "timestamp", payment.ProcessedAt,
   49 |     )

goroutine 1 [running]:
github.com/nikoksr/assert-go.PaymentProcessing(0xc00011c000)
    /app/payment.go:43 +0x1b4
# ... regular Go stacktrace continues
```

### Two-Tier Assertion System

The library provides two types of assertions:

1. `Assert()` - Always active, meant for critical checks that should run in all environments
2. `Debug()` - Development-time assertions that can be disabled in production

#### Using Debug Assertions

Debug assertions are disabled by default. To enable them, use the `assertdebug` build tag:

```bash
go test -tags assertdebug ./...
go run -tags assertdebug main.go
```

Example usage:

```go
// This will only be evaluated when built with -tags assertdebug
assert.Debug(len(items) < 1000, "items list too large",
    "current_length", len(items),
    "max_allowed", 1000,
)

// This will always be evaluated regardless of build tags
assert.Assert(response != nil, "HTTP response cannot be nil",
    "status_code", response.StatusCode,
)
```

### Configuration

You can configure the assertion behavior:

```go
// Configure assertion behavior
assert.SetConfig(assert.Config{
    // Enable/disable source context in error messages
    IncludeSource: true,
    // Number of context lines to show before and after the failing line
    ContextLines:  5,
})
```

## Philosophy

- **Minimal**: Single-purpose library that does one thing well
- **Context over complexity**: Rich debugging information without complex APIs
- **Clear failures**: Source context shows exactly where and why things went wrong
- **Idiomatic Go**: Feels natural in your Go codebase
