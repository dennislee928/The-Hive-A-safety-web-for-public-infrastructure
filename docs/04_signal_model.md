# Signal Model

## Signal sources
- Infrastructure (station telemetry, coarse people-flow)
- Staff reports
- Emergency calls (where device-assisted location is supported)
- Crowd reports (app)

## Aggregation rule
Crowd and infrastructure signals should be aggregated into **time-windowed, zone-level summaries** before reaching the decision layer.

## Trust scoring (concept)
- Rate limiting
- Reputation / device integrity signals (optional)
- Cross-corroboration in a time window

