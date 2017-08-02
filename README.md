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
