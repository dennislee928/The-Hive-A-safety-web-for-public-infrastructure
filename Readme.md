# ERH Crowd-in-the-Loop Safety Decision PoC (Spec-Only)

This repository is a **specification-only** (Markdown-only) framework for a PoC that explores a **crowd-in-the-loop safety decision system** across four high-density zones:

1. **Station Interior** (concourse / gates / platforms)
2. **Train Car** (inside carriage)
3. **Station Perimeter** (station surroundings, entrances/exits, transfer corridors)
4. **Other High-Density Areas** (events, plazas, festivals, dispersal flows)

The PoC assumes two public-communication routes must coexist:

- **Route 1 (Baseline / No Download):** Standards-based public warning delivery (e.g., cell broadcast / location-based SMS) with a CAP-centered alert format.
- **Route 2 (App / Optional Download):** A public app that enables bidirectional interactions (structured crowd reports, personalized guidance, check-in), while adding abuse and privacy controls.

A core constraint is **ERH governance**: as system complexity increases, critical misjudgments ("ethical primes") must remain **measurably bounded** rather than exploding with scale.

## What you get

- A PoC blueprint that you can copy into a real engineering repo later.
- A governance-first approach that treats **false positives, false negatives, and systematic bias** as first-class risks.

## What you do not get

- No production-ready safety system.
- No surveillance instructions.
- No operational playbook for law-enforcement tactics.

## Start here

- [`agent.md`](./agent.md) — agent roles and operating constraints
- [`plan.md`](./plan.md) — PoC plan, milestones, and acceptance criteria
- [`structure.md`](./structure.md) — repository layout and how docs fit together

