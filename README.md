# easygoer 🐎

[![Go Reference](https://pkg.go.dev/badge/github.com/mktbsh/easygoer.svg)](https://pkg.go.dev/github.com/mktbsh/easygoer)
[![Go Report Card](https://goreportcard.com/badge/github.com/mktbsh/easygoer)](https://goreportcard.com/report/github.com/mktbsh/easygoer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**easygoer** is a robust and handy utility library for Go.
Inspired by the legendary racehorse [_Easy Goer_](https://en.wikipedia.org/wiki/Easy_Goer), this library aims to make your Go development run smoother and faster.

It provides a collection of helper functions for slices, strings, maps, and system operations that are frequently used but not included in the standard library.

## 🚀 Installation

```bash
go get github.com/mktbsh/easygoer
```

## 📦 Packages

### crypto/envelope

Provides envelope encryption and decryption functionality using AES-256-GCM.

```go
import "github.com/mktbsh/easygoer/crypto/envelope"

// Generate a KEK (Key Encryption Key)
kek, err := envelope.GenerateKEK()

// Encrypt data
data := []byte("Secret message")
env, err := envelope.Encrypt(data, kek)

// Decrypt data
decrypted, err := envelope.Decrypt(env, kek)
```

### dirs

XDG Base Directory specification compliant directory management.

```go
import "github.com/mktbsh/easygoer/dirs"

paths, err := dirs.Resolve("myapp")
// Access paths.ConfigDir, paths.DataDir, etc.
```

## 💡 Usage

Import the package into your code:

```go
import "github.com/mktbsh/easygoer"
```
