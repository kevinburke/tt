# tt

`tt` is a faster command runner for Javascript tests. It assembles useful
environment variables, and useful arguments that you always want to pass to
Mocha, and then invokes Mocha with those arguments.

### Why should I use this instead of (gulp/grunt/whatever)

Those tools are often the bottleneck in your test suite. `tt` starts instantly,
but those tools can take 200ms or longer to print the version string. This can
slow you down.

### Install

Add the following to your Makefile:

```
TT := node_modules/.bin/tt

$(TT):
ifeq ($(shell uname -s), Darwin)
	curl --location --silent https://github.com/kevinburke/tt/releases/download/0.3/tt-darwin-amd64 > $(TT)
else
	curl --location --silent https://github.com/kevinburke/tt/releases/download/0.3/tt-linux-amd64 > $(TT)
endif
	chmod +x $(TT)

test: $(TT)
	$(TT)
```

If a user runs `make test`, that will download the appropriate binary and place
it in `node_modules/.bin`, then run tests. (If a user already has it, it'll just
run `make test`).
