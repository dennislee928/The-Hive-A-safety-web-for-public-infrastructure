# ADR 0003 — Route 1 + Route 2 Must Coexist

## Decision
Implement two parallel public communication routes:

- Route 1: baseline (no download)
- Route 2: app (optional download)

Route 2 must not substitute for Route 1.

## Rationale
Voluntary app downloads rarely achieve sufficient penetration; baseline channels are required for equitable reach.

## Consequences
- CAP becomes the canonical message format.
- Channel adapters render consistent instructions.

