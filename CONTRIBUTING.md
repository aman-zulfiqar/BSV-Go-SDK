# Contributing

Thanks for contributing to BSV-Go!

## Development setup
- Go 1.21+
- golangci-lint (optional): see Makefile `lint` target

## Workflow
1. Fork the repo and create a feature branch.
2. Run:
   - `make deps`
   - `make fmt`
   - `make lint` (if you have golangci-lint installed)
   - `make test`
3. Commit with clear messages.
4. Open a PR describing:
   - What changed and why
   - Any breaking changes
   - Tests and docs updated

## Tests
- Network-dependent examples may fail without funded wallets; tests should not require internet or should be gated.
- Prefer deterministic unit tests; avoid flakiness.

## Code style
- Keep exported symbols documented.
- Prefer clear names over abbreviations.
- Avoid catching errors without handling.

## Security
- The current sharding implementation is a demo. Do not market as production SSS.
- Report vulnerabilities privately (open a security issue or email the maintainer).


