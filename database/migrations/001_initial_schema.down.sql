-- Rollback initial schema

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS cap_messages;
DROP TABLE IF EXISTS device_report_history;
DROP TABLE IF EXISTS device_trust_scores;
DROP TABLE IF EXISTS keepalive_sessions;
DROP TABLE IF EXISTS approval_requests;
DROP TABLE IF EXISTS decision_states;
DROP TABLE IF EXISTS aggregated_summaries;
DROP TABLE IF EXISTS signals;

