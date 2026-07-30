# Quality gates

## Required inputs

- exact target or queued question;
- ROM identity;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. Run gofmt.
2. Run go test ./....
3. Run go vet ./....
4. Validate JSON schemas.
5. Scan tracked files.
6. Check broken references.
7. Audit manifests.
8. Write report.

## Required outputs

- updated canonical record;
- raw evidence references and hashes;
- confidence and alternatives;
- relationship/index updates;
- exact next action;
- checkpoint when the unit stops.
