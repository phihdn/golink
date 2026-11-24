# Lessons learned from the PR #2 review

Source: instructor review on [PR #2 Lesson 2](https://github.com/phihdn/golink/pull/2), commit `79e88ae`, 2026-11-23. Six comments, which boil down to four lessons, plus a follow-up from the instructor at the end.

| Lesson | Comments it comes from | Applied in |
| --- | --- | --- |
| L1. Keep scratch programs out of the repo | `cmd/test/main.go` | [`01a1466`](https://github.com/phihdn/golink/commit/01a1466) |
| L2. Pass one context through every layer | `handler/health.go:29`, `handler/shortenurl.go:62` | [`7ffc6a2`](https://github.com/phihdn/golink/commit/7ffc6a2) |
| L3. Tests must prove that exact context arrives | `handler/health_test.go:47`, `handler/shortenurl_test.go:37` | [`2448647`](https://github.com/phihdn/golink/commit/2448647) |
| L4. Create the context in the test and expect that same value | `service/health_test.go:45` | [`fb920e4`](https://github.com/phihdn/golink/commit/fb920e4) |

## L1. Keep scratch programs out of the repo

Original comment: "Mình ko commit file test lên repo" (we don't commit the test file to the repo).

`cmd/test/main.go` was a throwaway program from the video to poke at Redis by hand. It is not part of the product, it has no tests, it ignores errors, and it becomes a second binary that readers wonder about.

Rule: anything you write only to try something out stays local. Real checks belong in `_test.go` files.

Applied in `01a1466`: deleted `cmd/test/main.go` and its README row. Nothing was lost, because the one fact the program showed (a `CACHE` prefix picks its own settings) is already covered by `TestNewConfig_PrefixWinsOverPlainName` in `pkg/redis/config_test.go`.

```bash
# keep experiments outside the module, or in a folder git ignores
mkdir -p ~/scratch/redis-try && cd ~/scratch/redis-try
# or add a local-only ignore that is never committed
echo "cmd/test/" >> "$(git rev-parse --git-common-dir)/info/exclude"
```

If the experiment proves something worth keeping, turn it into a test:

```go
func TestNewClient_PrefixUsesSeparateDB(t *testing.T) {
    t.Setenv("CACHE_REDIS_DB", "1")
    client, err := redis.NewClient("CACHE")
    require.NoError(t, err)
    assert.Equal(t, 1, client.Options().DB)
}
```

## L2. Pass one context through every layer

Original comments: "Mình ko xài `c.Request().Context()` ở đây nha, truyền cả `c` vào" (don't use `c.Request().Context()` here, pass `c` itself), and "đây cũng vậy" (same here).

The idea: one request has one context. It carries the deadline, the cancel signal when the client disconnects, and request-scoped values such as a trace id. Handler → service → repository → Redis must all receive that same context, so cancelling the request stops the Redis call too.

How you pass it depends on the framework:

| Framework | Does the framework context implement `context.Context`? | What to pass to the service |
| --- | --- | --- |
| Gin (the video) | Yes, `*gin.Context` has `Deadline`, `Done`, `Err` and `Value` | `c` |
| Echo v5 (this project) | No, `*echo.Context` does not (checked with a type assertion) | `c.Request().Context()` |

Gin, as in the video:

```go
func (h *shortenLink) CreateShortenLink(c *gin.Context) {
    code, err := h.svc.CreateShortenLink(c, input.URL, exp) // c is a context.Context
    // ...
}
```

Echo, the equivalent:

```go
func (h *shortenLink) CreateShortenLink(c *echo.Context) error {
    code, err := h.svc.CreateShortenLink(c.Request().Context(), input.URL, exp)
    // ...
}
```

Do not change the service to take `*echo.Context` or `*gin.Context`. The service would then depend on the web framework, which breaks clean architecture: services must only know `context.Context`.

Applied in `7ffc6a2`. The code already followed the rule: both handlers pass `c.Request().Context()`, and the services and repositories pass `ctx` on to Redis. The commit adds a comment at both call sites so a reader coming from Gin knows why the code does not pass `c`. The one layer that broke the chain is the health service's `context.WithTimeout`, fixed in L4. Proving that the context really arrives is L3's job.

## L3. Tests must prove that exact context arrives

Original comments: "Vì mình truyền c.Request().Context() nên ở đây mình mới phải truyền `mock.Anything`. Như thế này sẽ ko đảm bảo 1 context đc pass xuyên suốt các layer" (because we pass `c.Request().Context()` we had to use `mock.Anything`, which does not guarantee one context is passed through all layers), and "Không xài mock.Anything cho context" (don't use `mock.Anything` for context).

`mock.Anything` accepts any value. If the handler passed `context.Background()` by mistake, the test would still pass, and in production the Redis call would no longer stop when the client disconnects.

Rule: build a context the test owns, put it on the request, and expect that exact value.

Weak test, passes even when the context is lost:

```go
svc.EXPECT().Check(mock.Anything).Return(health, nil)

req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
```

Strong test, fails if any layer swaps the context:

```go
ctx := t.Context()
svc.EXPECT().Check(ctx).Return(health, nil)

req := httptest.NewRequest(http.MethodGet, "/health-check", nil).WithContext(ctx)
rec := httptest.NewRecorder()
c := echo.New().NewContext(req, rec)

require.NoError(t, NewHealth(svc).Check(c))
```

`t.Context()` (Go 1.24+) gives every test and subtest its own context, so it is never equal to `context.Background()`. A handler that replaced the request context with `context.Background()` makes the mock fail. Do not use `context.Background()` as the test context: every call returns the same value, so that mistake would still pass.

The same applies one layer down. In service tests, expect the context you passed in, not `mock.Anything`:

```go
ctx := t.Context()
storage.EXPECT().StoreURL(ctx, "abc1234", url, time.Hour).Return(nil).Once()

code, err := NewShortenUrl(storage, gen).CreateShortenLink(ctx, url, time.Hour)
```

Applied in `2448647` (the first version used a custom helper, replaced later, see the follow-up section):

- `context_test.go` in the `handler` and `service` packages adds `newRequestContext()`, which builds the `context.WithValue` context above. Each test case calls it, so every case owns its context.
- `handler/health_test.go` and `handler/shortenurl_test.go` put that context on the request with `req.WithContext(ctx)` and expect it in the service mock.
- `service/shortenurl_test.go` passes it to `CreateShortenLink` and expects it in every `StoreURL` call, including the retries.
- No `mock.Anything` is left for a context in these tests. The health service test still uses `mock.MatchedBy`, which is L4.

Proof that the tests now guard the rule: changing the health handler to pass `context.Background()` makes `TestHealth_Check` fail with "mock: Unexpected Method Call".

## L4. Create the context in the test and expect that same value

Original comment: "Mình tạo context trước, sau đó expect đúng context đó thay vì xài mock.MatchedBy như thế này" (create the context first, then expect that exact context instead of using `mock.MatchedBy` like this).

`mock.MatchedBy(func(ctx) bool { _, ok := ctx.Deadline(); return ok })` only checks a property, so any context with a deadline passes. It is a weaker version of `mock.Anything`.

It was needed here because the health service wrapped the context:

```go
func (h *health) Check(ctx context.Context) (model.Health, error) {
    ctx, cancel := context.WithTimeout(ctx, pingTimeout) // a new context, not the caller's
    defer cancel()
    return ..., h.r.Ping(ctx)
}
```

A layer that silently derives a new context makes "same context everywhere" impossible to assert. Keep the service a straight pass-through, and the test becomes exact:

```go
func (h *health) Check(ctx context.Context) (model.Health, error) {
    if err := h.r.Ping(ctx); err != nil {
        return model.Health{}, err
    }
    return model.Health{Message: "OK", ServiceName: h.serviceName, InstanceID: h.instanceID}, nil
}
```

```go
ctx := t.Context()
repo.EXPECT().Ping(ctx).Return(nil).Once()

got, err := NewHealth("golink", "abc-123", repo).Check(ctx)
```

Trade-off to remember: without the wrap, an unreachable Redis host makes the health check wait for go-redis's own timeouts (about 5 seconds to dial). If a short limit is needed, set it once at the edge, for example in `redis.Options` (`DialTimeout`, `ReadTimeout`) or with an HTTP timeout middleware, not hidden inside a service.

Applied in `fb920e4`:

- `service/health.go` drops `pingTimeout` and the `context.WithTimeout` wrap. `Check` passes `ctx` straight to `Ping`.
- `service/health_test.go` creates the context with `newRequestContext()`, expects `Ping(ctx)`, and calls `Check(ctx)`. The `mock.MatchedBy` deadline check is gone.
- The README no longer promises a 1 second timeout for `GET /health-check`.

Since this commit, no test in the project uses `mock.Anything` or `mock.MatchedBy` for a context.

## Follow-up: use `t.Context()` for the test context

After seeing the fixes, the instructor approved the PR and added a tip. He had not noticed the project uses Echo, so he accepted `c.Request().Context()`. For tests he suggested `t.Context()` instead of building a context by hand:

```go
func TestEcho(t *testing.T) {
    ctx := t.Context()

    req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
    rec := httptest.NewRecorder()

    c := echo.New().NewContext(req, rec)

    assert.Equal(t, ctx, c.Request().Context()) // true
}
```

The snippet shows that Echo hands back the exact context put on the request. So a test only needs to create `t.Context()`, attach it with `WithContext(ctx)`, and expect it in the mock. That keeps one context through every layer without `mock.Anything`.

Why `t.Context()` beats the hand-made helper:

- It is built into `testing`, so there is no `ctxKey` type or helper file to maintain.
- Each test and subtest gets its own value, which is all the mock comparison needs.
- It is cancelled when the test ends, so code that keeps working after the test finishes sees a cancelled context.

One detail in his message: he wrote `mocks.NewHealth(t, ctx)`. The mockery constructor only takes `t`. The context goes into the expectation instead: `svc.EXPECT().Check(ctx)`.

Applied in the "use t.Context for the test context" commit: the handler and service tests now use `ctx := t.Context()`, and both `context_test.go` helper files are deleted. The proof still holds: changing the health handler to pass `context.Background()` makes `TestHealth_Check` fail.

## Checklist for the next PR

- [ ] No scratch programs or debug files in the diff.
- [ ] Every service and repository method takes `ctx context.Context` as its first argument.
- [ ] Handlers pass the request context (`c` in Gin, `c.Request().Context()` in Echo), never `context.Background()`.
- [ ] Services and repositories pass `ctx` on unchanged.
- [ ] Tests use `ctx := t.Context()`, put it on the request with `WithContext(ctx)`, and expect that exact value. No `mock.Anything` or `mock.MatchedBy` for contexts.
