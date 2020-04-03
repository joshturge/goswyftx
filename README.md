# Goswyftx

A go wrapper around the swyftx API

## Getting Started

The following instructions will help in using this package

### Usage

In order to use this library, you need to setup an api key which can be done
within the profile section under API KEYS. Once you have the key, you're all set
to use the library.

#### Installing

```bash
go get github.com/joshturge/goswyftx
```

After installing the package you need to create a new client which will ask for
an authentication token from the API:

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
