nq(1) - A simple last-in-queue command executor

```
# Start the server

$ nq -s
INFO[21:06:15.9488] Server requested. Starting server
...
INFO[21:06:15.952] Setting up worker
INFO[21:06:15.9521] Worker started


# In another shell - push commands to the server

$ nq -- sh -c 'make build test'  # Queued, starts running, returns immediately, does not block
$ nq -- sh -c 'make build test'  # Queued, likely never runs
$ nq -- sh -c 'make build test'  # Queued, likely never runs
$ nq -- sh -c 'make build test'  # Queued, runs
```

## Why?

In most build pipelines, processing build tasks/jobs is a blocking activity
and this causes a problem in managing queues where multiple _fire-and-forget_
jobs can be submitted.

If jobs are submitted at a rate  greater than the processing time - two
unwanted things arise

- (1) the queue length increases
- (2) new commands cannot be submitted as the previous commands block

`nq` allows the non-blocking intake of multiple (identical) jobs and will
begin processing the last job in the job queue, ignoring and discarding all
the others that precede it.

Additionally, the `nq` client does not block, it returns immediately after pushing the
job command to the server. This frees up the shell allowing `nq` to be
invoked again. In this mode `nq` does not wait to check the exit status
of the command(s) submitted to the queue.

### Use case - Triggering builds from file changes in the editor

e.g. in vim + tmux/tslime, you may wish to run tests on every buffer/file write

```
:au BufWritePost,FileWritePost * silent :Tmux make build test
```

Here `make` blocks until it has finished processing.

With `nq` wrapping the command, flow is non-blocking

```
:au BufWritePost,FileWritePost * silent :Tmux nq -- make build
```

### Use case - Running builds from file changes detected by inotify, etc

`nq` can listen in on STDIN being a pipe and react only when the input
matches a regular expression.

```
inotifywait -q -e close_write -m . | nq -p CREATE -- make build test
```

## Integrations with custom clients

Simply write an array of strings as a JSON string to the named pipe.

```
{ echo -e '["sh", "-c", "'echo foo'"]';  } > "$TMP/nq.fifo"
```

## Building nq

```
$ make build
```
