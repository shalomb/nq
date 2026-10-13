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

If jobs are submitted at a rate greater than the processing time - two
unwanted things arise:

- (1) the queue length increases
- (2) new commands cannot be submitted as the previous commands block

`nq` allows the non-blocking intake of multiple (identical) jobs and will
begin processing the last job in the job queue, ignoring and discarding all
the others that precede it.

Additionally, the `nq` client does not block, it returns immediately after pushing the
job command to the server. This frees up the shell allowing `nq` to be
invoked again. In this mode `nq` does not wait to check the exit status
of the command(s) submitted to the queue.

## How It Works: Latest-Only Execution (Conflation)

Instead of a standard First-In, First-Out (FIFO) queue or a timer-based debounce, `nq` implements a **coalescing/conflating queue**:

```text
Client 1 (Save 1) ──┐
Client 2 (Save 2) ──┼──> [ FIFO Named Pipe ] ──> Server Enqueue
Client 3 (Save 3) ──┘                                   │
                                                        ▼
                                                  [ Job Queue ]
                                                        │
                      ┌─────────────────────────────────┴─────────────────────────────────┐
                      ▼                                                                   ▼
           Takes LATEST item (Job 3)                                     DISCARDS all previous items
              q.queue.Back()                                                    q.queue.Init()
                      │
                      ▼
             Execute Job 3 (Worker)
```

1. **Fire-and-Forget Submission**: The client serializes command arguments into JSON, writes them to the UNIX named pipe (`$TMP/nq.fifo`), and exits immediately.
2. **Immediate Start on Idle**: When the worker is idle, an incoming command starts processing immediately without waiting for a debounce timer.
3. **Automatic Discard of Intermediate Jobs**: When multiple commands arrive while a task is running, they accumulate in the queue. When the worker finishes and pulls the next task, it takes **only the newest job** (`q.queue.Back()`) and immediately clears the rest (`q.queue.Init()`).

### How It Compares

| Strategy | Behavior Under Bursts | Trade-off |
| :--- | :--- | :--- |
| **Standard FIFO Queue** | Executes every job (1, 2, 3...) sequentially | Backlog grows; wastes CPU and delays feedback by running stale code versions. |
| **Debounce Timer** | Waits for $N$ ms of quiet before starting | Adds an artificial delay to the first job even when idle. |
| **`nq` Conflation** | Runs first job immediately; runs latest job next and drops intermediate ones | Immediate initial feedback, zero backlog accumulation, always tests the latest code. |

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
