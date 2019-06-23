# gomniauth-traq

[gomniauth](https://github.com/stretchr/gomniauth)のtraQ用Provider

## 例

```go
gomniauth.WithProviders(
  github.New("key", "secret", "callback"),
  google.New("key", "secret", "callback"),
  gomniauth_traq.New("key", "secret", "callback"),
)
```
