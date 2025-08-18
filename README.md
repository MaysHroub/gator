# Gator 🐊 - The Blog Aggregator

A simple command-line tool for managing RSS feeds.  
Multiple users can log in and use it. Each user can add, follow, and unfollow feeds, and check their posts.

It uses a PostgreSQL database to store users, feeds, follow relationships, and posts. I used the CLI tool `goose` to manage database migrations and `sqlc` to generate Go code from SQL queries. And used the `testify` & `httptest` packages to write the the unit tests.

This project was built using a TDD approach (although it wasn’t required) where tests are written before the implementation. It was challenging and sometimes frustrating, but ultimately rewarding since I learned more about unit testing in Go and how to use mocking to build unit tests—especially for database and external server interactions.

And shout out to ChatGPT for its coding tips. I don’t have the experience to judge the *cleanliness* of the code, but I can say it helped me organize this project in a cleaner way than I could on my own.


## Commands

To see all available commands, run:
```bash
gator cmnds
```

For help on a specific command:
```bash
gator man <command-name>
```

I won’t spoil all the supported commands here—check them yourself :)
<br>
But here's a thing:

### The `agg` command runs in the foreground

To stop it, use `Ctrl + c` keys.

To make it run in the background, there are two ways:

#### 1. Quick and Dirty (temporary background run)
Just append `&` at the end:
```bash
gator agg 10s &
```

To kill it:

1. Get its job number from the output of `jobs` command
2. Bring it to the foreground `fg %<job-number>`
3. Kill it with `kill %<job-number>` 

Or just get its process ID and run `kill <PID>`.

Nevertheless, it stops if you close the termainal.


#### 2. More reliable (survives terminal close)

Use `nohup`:
```bash
nohup gator agg > agg.log 2>&1 &
```

This detaches the process from the terminal and the output will go into `agg.log`.

The process keeps running even if you log out or close the terminal.

You can kill it later with `kill <PID>`.

## Requirements

- Go 1.20+  
- PostgreSQL (running locally or remotely)  
- (Optional) `goose` for database migrations


## Installation

You have two easy options:

### 1. Install via `go install`

You need to have Go installed in your machine.

```bash
go install github.com/MaysHroub/gator@latest
```

### 2. Build from source

You need to have Git installed in your machine.

```bash
git clone https://github.com/MaysHroub/gator.git
cd gator
go build -o gator
```

### Then you can run

```bash
./gator <command-name> [command-args]
```


## Usage

This is an example of how you might use Gator:

```bash
gator register mays-alreem
gator addfeed https://techcrunch.com/feed/
gator follow https://news.ycombinator.com/rss
gator agg
```


## Resources

This project was built as part of the Backend-Go track on [boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang).

The following were good resources that helped me learn how to write solid unit tests in Go:

[How to write unit tests in Go](https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package)

[Learn Go using TDD - Mocking](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/mocking) (This one is super good, but check the dependency injection part first)

[Testing in Go with httptest](https://speedscale.com/blog/testing-golang-with-httptest/)

<br>

I didn't read these ones before, but they seem solid:

[Comprehensive Guide to Testing in Go](https://blog.jetbrains.com/go/2022/11/22/comprehensive-guide-to-testing-in-go/)

[How to use mock testing in Go](https://www.jetbrains.com/guide/go/tutorials/mock_testing_with_go/) (This one is also good)

I don't understand why good resources only show up when I finish.

