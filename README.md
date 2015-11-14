# LockAPI

*Minimalist Go cross-platform API for locking*


## Introduction

LockAPI is a small but useful and efficient Go library dedicated to locking files via OS-specific system calls:

* **fcntl** on *Linux* and *Mac OS*
* **LockFileEx** and **UnlockFileEx** in *kernel32.dll* on *Windows*


These details are masked by a simplified, cross-platform facade provided by the library.


## Installation

The installation requires, as usual, *go get*:

> go get github.com/giancosta86/LockAPI/lockapi


## Requirements

LockAPI has been created with Go 1.5.

It has been tested on Windows (8.1) and Linux (Xubuntu 15.10); support for Mac OS (Darwin) has been introduced, but is currently untested.

To test LockAPI on your current system, run:

> go test github.com/giancosta86/LockAPI/lockapitest -v

after installing the the library.


## Usage

First of all, import the library package into a source file:

```go
import "github.com/giancosta86/LockAPI/lockapi"
```

From now on, you can use the functions available in the **lockapi** package.


## Reference

The full GoDoc reference is available at [this page](https://godoc.org/github.com/giancosta86/LockAPI/lockapi).


## Example programs

Simple example programs can be found in the test package:

* [Blocking file locking](lockapitest/lockapi_blocking_file_test/main.go)
* [Non-blocking file locking](lockapitest/lockapi_non_blocking_file_test/main.go)
