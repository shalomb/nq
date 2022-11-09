nq(1) - A simple last-job-in-queue command executor

```
# Start the server

$ nq -s
...
INFO[20:36:29.4928] Starting server


# Push commands to the server

$ nq -- sh -c 'make build test'  # Returns immediately, does not block
$ nq -- sh -c 'make build test'
```

## Why?

In most CI pipelines, build tasks/jobs are blocking and this causes queue
management problems.

If the system  triggering the CI commands  does so at a rate  greater than the
rate the CI system can process jobs  - two unwanted things arise (1) the queue
length increases (2) new commands cannot be submitted as the previous commands
block.

`nq` allows the non-blocking intake of multiple (identical) jobs and
will begin processing the last job taken, disregarding all previous jobs.

As the `nq` client does not block, it cannot wait to check the exit status
of the command(s) submitted to the queue. In a future version, the ability
to query the health of the "builds" will be possible.

### Use case - Triggering builds from file changes in the editor

e.g. in vim + tmux/tslime, you may wish to run tests on every buffer/file write

```
:au BufWritePost,FileWritePost * silent :Tmux make build
```

While this works for quick test suites (desirable), `make` blocks until
the tests have completed and so a subsequent `make ...` command cannot
be accepted over an already running command.

With `nq` wrapping the command, flow is non-blocking and you can make
as many edits and saves as you like - only the last issued command
is triggered.

```
:au BufWritePost,FileWritePost * silent :Tmux nq -- make build
```

### Use case - Running builds from file changes detected by inotify, etc

`nq` can listen in on STDIN being a pipe and react only when the input
matches a regular expression.

```
inotifywait -q -e close_write -m . | nq -p CREATE -- make build test
```

## Using a custom client

Simply write an array of strings as a JSON string to the named pipe.

```
{ echo -e '["sh", "-c", "'echo foo'"]';  } > "$TMP/nq.fifo"
```

## Building nq

```
$ make build
```
