# plan.md — PoC Framework Plan (Crowd-in-the-Loop Safety Decision)

## 1. Objective

Design a PoC specification for a **crowd-in-the-loop safety decision system** that reduces harm in **high-density public environments** by improving:

- Detection-to-response latency
- Decision reliability under escalating complexity
- Public communication consistency
- Abuse resistance and privacy compliance

The PoC explicitly supports **two parallel delivery routes**:

1. **Route 1 (Baseline / No Download):** cell broadcast / location-based SMS and other broadcast channels, anchored by a CAP message.
2. **Route 2 (App / Optional Download):** a citizen-facing app providing structured reporting and individualized safety guidance.

## 2. System scope: four Zones

### Z1 — Station Interior
- Concourse, gates, platforms, transfer corridors inside fare control.

### Z2 — Train Car
- Inside carriage; operational options rely on next-station handling and onboard guidance.

### Z3 — Station Perimeter
- Entrances/exits, surface crossings, adjacent sidewalks, bus interchanges.

### Z4 — Other High-Density Areas
- Events, plazas, festivals, dispersal flows (same modeling interface; different topology/resources).

## 3. Decision points (revised for 4 zones)

### D0 — Enter/Exit Pre-Alert
- Operators acknowledge a candidate incident and enable "pre-alert" monitoring.

### D1 — Resource Dispatch Recommendation
- Suggest dispatch level (station staff / security / first responders) based on corroborated signals.

### D2 — Local Zone Guidance Activation
- Turn on in-station / in-car guidance (signage/PA templates; app guidance) without causing panic.

### D3 — Escalate Zone Alert Level (High impact)
- Escalate a single Zone (Z1/Z2/Z3/Z4) to an elevated state.
- Requires: dual control + dead-man keepalive + TTL + rollback plan.

### D4 — Multi-Zone / Network-Level Coordination (High impact)
- Coordinate multiple zones or multi-station posture.
- Requires stricter governance: dual control + additional approval (policy-defined).

### D5 — Public Warning Broadcast (High impact)
- Release a CAP-based public warning (baseline channels + app push).
- Requires consistency checks, human confirmation, and TTL.

### D6 — De-escalation and Evidence Sealing
- Roll back elevated states; seal logs and artifacts for audit and improvement.

## 4. Signals model

### 4.1 Signal categories
- Fixed infrastructure signals (e.g., station sensors, operations telemetry)
- Staff reports (station personnel)
- Emergency calls (where available, device-assisted location sharing)
- Crowd reports (app-based, structured)

### 4.2 Crowd reporting rules (Route 2)
- Structured report with coarse zone, time window, and confidence.
- Rate-limited per device/account.
- Weighted by trust score (without requiring identity disclosure in the PoC).
- Never a direct trigger for high-impact actions.

## 5. ERH governance: complexity definition

The PoC defines complexity as a **vector**:

- **x_s (signals):** number of effective signal sources (after aggregation)
- **x_d (decision depth):** decision graph depth / number of high-impact gates
- **x_c (context states):** number of scenario/context states the policy must cover

We track a derived scalar for reporting:

- **x_total:** a weighted function of (x_s, x_d, x_c)

### 5.1 Ethical primes (critical misjudgments)
- **FN-prime:** failure to escalate or dispatch when needed
- **FP-prime:** unnecessary escalation / overreaction producing panic or rights violations
- **Bias-prime:** systematic error affecting protected or vulnerable groups
- **Integrity-prime:** false/forged signals or commands that change outcomes

### 5.2 ERH objective (PoC-level)
As x_total increases (more signals, deeper decisions, richer contexts), the PoC must demonstrate:

- Misjudgment growth remains **bounded** and measurable.
- High-impact actions remain governed by human confirmation and dual control.

## 6. Route 1 (Baseline / No download)

### 6.1 Canonical alert format
- CAP as the canonical message structure.

### 6.2 Delivery channels (examples)
- Cell broadcast
- Location-based SMS
- Public signage, PA, TV/radio, social channels (policy-defined)

### 6.3 Principles
- Route 1 is the minimum guarantee.
- Route 2 must never be a substitute for Route 1.

## 7. Route 2 (App / Optional download)

### 7.1 App capabilities
- Structured crowd reports
- Personalized guidance and routing suggestions (coarse)
- Safety check-in and assistance requests
- Post-incident feedback (for model/process improvement)

### 7.2 App governance constraints
- Minimal data collection
- Coarse location handling (zone-level)
- Abuse controls: rate limiting, trust scoring, corroboration

## 8. Privacy / legal / abuse risk register (PoC deliverable)

The PoC must include a risk register covering:

- Purpose limitation and minimization
- Retention limits
- Access control and audit
- Misuse scenarios (spam, brigading, panic induction)

## 9. Evaluation plan

### 9.1 Primary metrics
- Time-to-acknowledge (TTA)
- Time-to-dispatch recommendation (TTDR)
- False negative rate (FN) for high-severity scenarios
- False positive rate (FP) for low-severity scenarios
- Bias indicators (group-agnostic in PoC; defined as distributional parity on proxies)

### 9.2 ERH reporting
- Plot E(x_total) for each prime (FN/FP/Bias/Integrity)
- Identify breakpoints where misjudgment growth accelerates
- Demonstrate mitigation via aggregation, gating, and dual control

## 10. Milestones (spec-only)

- **M1 — Requirements & threat model** (docs complete)
- **M2 — Zone models & scenario library** (docs + templates complete)
- **M3 — ERH governance + evaluation dashboard spec** (metrics + reporting spec)
- **M4 — Route 1/2 message and workflow specs** (CAP templates, gating, rollback)

## 11. Acceptance criteria

The PoC specification is acceptable if it:

- Supports all four zones with a unified modeling interface
- Implements both Route 1 and Route 2 with CAP-centered consistency
- Defines x_s, x_d, x_c and how to compute x_total
- Defines ethical primes and evaluation metrics
- Includes a privacy/legal/abuse risk register with mitigations

