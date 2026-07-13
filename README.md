
# Prunning Radix Trie

A Go implementation of a pruning radix trie for prefix search and top-k term
suggestions. It is based on the original
[C# PruningRadixTrie](https://github.com/wolfgarbe/PruningRadixTrie) and keeps
the example small enough to use as a reference implementation.

## Requirements

- Go 1.26 or newer

## Usage

Run the test suite:

```sh
make test
```

Run the benchmark:

```sh
make bench
```

Check reachable vulnerabilities:

```sh
make vulncheck
```

Run an arbitrary command through the Makefile:

```sh
make run CMD='go test ./...'
```

## Maintenance Checklist

- [x] Update the Go runtime directive to Go 1.26.
- [x] Add a runnable example sanity test for prefix suggestions.
- [x] Fix split-node insertion so shared prefixes are not returned as fake terms.
- [x] Add Makefile targets for tests, benchmarks, vulnerability checks, and arbitrary commands.
- [x] Run `make test`.
- [x] Run `make vulncheck`; no reachable vulnerabilities were found.
