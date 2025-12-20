# agent.md — Working Agreement for Agents (Spec Repo)

This file defines a **multi-agent working agreement** for producing and maintaining the PoC specification.

## 1. Scope and non-goals

### In scope
- Architecture and governance for a crowd-in-the-loop safety decision PoC.
- Data-flow and trust boundaries for: sensors, operations center, station edge, and citizen mobile devices.
- ERH-based measurement of misjudgment growth as complexity increases.
- Privacy, legal, and abuse-mitigation constraints (Taiwan PDPA-aligned).

### Out of scope
- Detailed instructions that could enable violence.
- Detailed tactics to bypass security, evade detection, or exploit public infrastructure.
- Collection of unnecessary personal data.

## 2. Roles

### 2.1 Systems Architect Agent
- Owns end-to-end component boundaries, interfaces, and failure modes.
- Ensures "fail-safe" defaults for high-impact actions.

### 2.2 Governance / ERH Agent
- Defines ethical primes (critical misjudgments) and ERH measurement plan.
- Owns complexity vector definition and reporting.

### 2.3 Privacy & Legal Agent
- Maps data items to lawful basis, purpose limitation, minimization, retention, and access control.
- Produces DPIA-style risk register and mitigations.

### 2.4 Abuse & Safety Red-Team Agent
- Enumerates plausible abuse (spam reports, coordinated brigading, false-flagging, panic induction).
- Defines rate limits, trust scoring, and "human-in-the-loop" gates.

### 2.5 Public Warning / CAP Agent
- Owns CAP message templates, channel mapping (baseline vs app), and consistency rules.

### 2.6 Simulation & Evaluation Agent
- Defines zone models, scenario library, and evaluation metrics.
- Ensures results are reproducible from documented assumptions.

## 3. Cross-cutting constraints

### 3.1 High-impact action gating
All actions that can significantly affect public movement or rights (e.g., station-wide controls, wide-area alerts) must be guarded by:
- **Two-person integrity / dual control** (at least two authorized operators)
- **Dead-man / keepalive hold-to-maintain** semantics for elevated states
- Explicit **time-to-live (TTL)** and automatic rollback

### 3.2 Safety communications consistency
- A single canonical alert (CAP) should be rendered across channels to avoid conflicting instructions.

### 3.3 Privacy-by-design
- Prefer **aggregation** over raw collection.
- Prefer **coarse zones** over fine-grained tracking.
- Use retention limits and access logging.

### 3.4 Abuse resistance
- Never treat crowd reports as a direct trigger for high-impact actions.
- Require corroboration, time-window validation, and operator confirmation.

## 4. Deliverable quality bar

Each doc change must include:
- Assumptions
- Attack/abuse considerations
- Privacy impacts
- How it affects ERH complexity and ethical primes

