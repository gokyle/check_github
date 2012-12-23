## check_github

This program is a simple check to see if Github is up. If Github is down, it
will keep checking. If it can find the "say" program on the path, it will
speak the status.

I have a festival-based say program: https://bitbucket.org/kisom/say

## Installing
```
$ go get bitbucket.org/kisom/check_github
$ go install bitbucket.org/kisom/check_github
$ check_github -t 5m
```

See `check_github -h` for command line flags.
