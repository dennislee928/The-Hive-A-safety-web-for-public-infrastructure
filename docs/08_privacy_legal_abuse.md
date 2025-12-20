# Privacy, Legal, and Abuse Controls

## Privacy-by-design (PoC)

- Prefer zone-level aggregation.
- Avoid continuous tracking; use time-window summaries.
- Minimize collection; limit retention; log access.

## Legal considerations (Taiwan PDPA-oriented)

- Purpose limitation and minimization
- Notice/consent or lawful basis (policy-defined)
- Vendor and processor controls

## Abuse scenarios to cover

- Spam/flooding crowd reports
- Coordinated brigading / false-flagging
- Panic induction via unauthorized broadcasts

## Mitigations

- Rate limiting and throttling
- Corroboration thresholds
- Dual control + dead-man keepalive for high-impact actions
- CAP canonical messages with consistency checks
