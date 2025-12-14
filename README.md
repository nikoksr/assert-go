<div align="center">

&nbsp;
<h1>assert-go</h1>
<p><i>Tiny (~100 LoC) Go assertion library focused on crystal-clear failure messages and thoughtful source context.</i></p>

&nbsp;

[![CI](https://github.com/nikoksr/assert-go/actions/workflows/ci.yml/badge.svg)](https://github.com/nikoksr/assert-go/actions/workflows/ci.yml)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/assert-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikoksr/assert-go)](https://goreportcard.com/report/github.com/nikoksr/assert-go)
</div>

&nbsp;

## About

- 🔍 Crystal-clear failure messages with contextual values
- 📚 Rich source context showing the exact failure location
- 🛠 Tiny and free of dependencies (~100 lines of Go)
- 💡 Elegant, idiomatic Go API
- 🎯 Two-tier assertion system with build tag support
- ⚙️ Configurable source context behavior
- ⚡ Zero-allocation hot path (~3ns per assertion)

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
  [timestamp]: "2025-12-12T15:04:05Z"

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
// Configure assertion behavior (call during initialization)
assert.SetConfig(assert.Config{
    // Enable/disable source context in error messages
    IncludeSource: true,
    // Number of context lines to show before and after the failing line
    ContextLines:  5,
})
```

**Note:** `SetConfig` should be called during program initialization before any assertions are made. It is not thread-safe and should not be called concurrently with assertions.

## Performance

Assertions are designed to have minimal performance impact:

```
BenchmarkAssert_Success              ~3.0 ns/op    0 B/op    0 allocs/op
BenchmarkAssert_SuccessWithValues    ~6.1 ns/op    0 B/op    0 allocs/op
BenchmarkDebug_Success (disabled)    ~0.3 ns/op    0 B/op    0 allocs/op
```

**Key Takeaways:**
- Successful assertions are extremely cheap (~3 nanoseconds)
- Zero allocations on the hot path
- Debug assertions when disabled are essentially free (compiler optimizes them away)
- Even with contextual values, overhead remains minimal

Run benchmarks yourself: `go test -bench=. -benchmem`

## A Personal Perspective on Assertions in Go

Like many Go developers, I initially dismissed assertions as incompatible with Go's philosophy of explicit error handling. That changed when I read [TigerStyle's take on assertions](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety) and decided to experiment.

Here's the problem I'd been living with: I've always felt the urge to validate internal invariants—checking that a logger I just initialized isn't nil, verifying state I just set up is correct, confirming assumptions about data flow within my own code. These aren't checks for user input or network failures. They're checks for *my* mistakes.

But the traditional Go approach felt wrong. Returning an error means telling my users: "Hey, handle this case where I might have screwed up." It pollutes APIs with impossible error cases, forcing callers to handle conditions that can only occur if my code is broken. I'd write these defensive checks anyway, feeling uncomfortable doing so, knowing I was treating my own bugs the same as legitimate operational failures.

Assertions solved this. They let me validate what *must* be true without burdening my API consumers. When an invariant is violated, there's no graceful recovery; continuing would only mask the bug and spread corrupted state. Better to fail immediately with rich context pointing directly at the problem.

I use assertions in **application code** where I control the full context and can make strong guarantees about internal state. I don't use them in libraries (where I can't control callers) or for validating external input (that's proper error handling territory). But for checking preconditions, postconditions, and invariants I own? They're essential.

This approach has also opened doors to patterns like negative-space testing, where assertions help verify not just what should happen, but what shouldn't. Worth exploring if you go down this path.

What started as a skeptical experiment has become fundamental to how I write Go. Assertions make my code more reliable and bugs dramatically easier to catch and diagnose.

## More Projects

If you find this library useful, you might also be interested in:

- **[notify](https://github.com/nikoksr/notify)** - Dead simple Go library for sending notifications to various messaging services (3,500+ ⭐)
- **[typeid-zig](https://github.com/nikoksr/typeid-zig)** - Complete Zig implementation of the TypeID specification, recognized as an official community implementation

---

<div align="center">
<sub>Built with ❤️ by <a href="https://github.com/nikoksr">@nikoksr</a></sub>
</div>
