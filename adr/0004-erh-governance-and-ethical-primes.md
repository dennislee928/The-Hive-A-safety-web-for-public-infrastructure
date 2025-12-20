# ADR 0004 — ERH Governance and Ethical Primes

## Decision
Use ERH governance to manage misjudgment growth with increasing complexity.

- Complexity vector: (x_s, x_d, x_c)
- Ethical primes: FN, FP, Bias, Integrity

## Rationale
Crowd-in-the-loop systems increase scale and heterogeneity; without governance, misjudgments can grow rapidly.

## Consequences
- Every feature must declare how it changes x_s/x_d/x_c.
- Evaluation must report E(x_total) per prime.

