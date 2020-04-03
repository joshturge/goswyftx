# Goswyftx

A go wrapper around the swyftx API

## Getting Started

The following instructions will help in using this package

### Usage

In order to use this library, you need to setup an api key which can be done
within the profile section under API KEYS on the [Swyftx website](https://swyftx.com.au).
Once you have the key, you're all set to use the library.

#### Installing

```bash
go get github.com/joshturge/goswyftx
```

After installing the package you need to create a new client which will
authenticate you with the api.

```go
// if you already have an authentication token it can be used in place of the
// empty string
client, err := goswyftx.NewClient("apiKey", "")
if err != nil {
    // handle error
}

version, err := client.Version()
if err != nil {
    // handle error
}

fmt.Println("Swyftx API version:", version)
```

With this client you can then access all the API endpoints available at the time
of writing.

### Testing

In order to run the unit test for this package you need the `API_KEY`
environment variable set to a swyftx api key. Then to run the test:

```bash
go test -v
```

### License

This package is licensed under the BSD License - see the [LICENSE](LICENSE) file
for details.
