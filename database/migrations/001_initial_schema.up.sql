-- Initial database schema for ERH Safety System
-- Based on docs/04_signal_model.md and docs/03_decision_points.md

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Signals table
CREATE TABLE IF NOT EXISTS signals (
    id VARCHAR(255) PRIMARY KEY,
    source_type VARCHAR(50) NOT NULL,
    source_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    zone_id VARCHAR(10) NOT NULL,
    sub_zone VARCHAR(100),
    signal_type VARCHAR(50),
    value JSONB,
    metadata JSONB,
    quality_score DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_source_type CHECK (source_type IN ('infrastructure', 'staff', 'crowd', 'emergency')),
    CONSTRAINT chk_zone_id CHECK (zone_id IN ('Z1', 'Z2', 'Z3', 'Z4'))
);

CREATE INDEX idx_signals_zone_time ON signals(zone_id, timestamp);
CREATE INDEX idx_signals_source_type ON signals(source_type);
CREATE INDEX idx_signals_source_id ON signals(source_id);
CREATE INDEX idx_signals_timestamp ON signals(timestamp);

-- Aggregated summaries table
CREATE TABLE IF NOT EXISTS aggregated_summaries (
    id VARCHAR(255) PRIMARY KEY,
    zone_id VARCHAR(10) NOT NULL,
    sub_zone VARCHAR(100),
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    source_count JSONB,
    weighted_value DECIMAL(10,4),
    confidence DECIMAL(3,2),
    signal_ids TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_agg_zone_id CHECK (zone_id IN ('Z1', 'Z2', 'Z3', 'Z4'))
);

CREATE INDEX idx_agg_summaries_zone_window ON aggregated_summaries(zone_id, window_start, window_end);
CREATE INDEX idx_agg_summaries_sub_zone ON aggregated_summaries(sub_zone);

-- Decision states table
CREATE TABLE IF NOT EXISTS decision_states (
    id VARCHAR(255) PRIMARY KEY,
    zone_id VARCHAR(10) NOT NULL,
    current_state VARCHAR(10) NOT NULL,
    aggregated_summary_id VARCHAR(255),
    signal_count INT,
    decision_depth INT,
    context_states INT,
    complexity_total DECIMAL(5,4),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_dec_zone_id CHECK (zone_id IN ('Z1', 'Z2', 'Z3', 'Z4')),
    CONSTRAINT chk_decision_state CHECK (current_state IN ('inactive', 'D0', 'D1', 'D2', 'D3', 'D4', 'D5', 'D6'))
);

CREATE INDEX idx_decision_states_zone_state ON decision_states(zone_id, current_state);
CREATE INDEX idx_decision_states_updated_at ON decision_states(updated_at);

-- Approval requests table (for dual control)
CREATE TABLE IF NOT EXISTS approval_requests (
    id VARCHAR(255) PRIMARY KEY,
    action_type VARCHAR(10) NOT NULL,
    zone_id VARCHAR(10) NOT NULL,
    proposal JSONB,
    requester_id VARCHAR(255) NOT NULL,
    approver1_id VARCHAR(255),
    approver2_id VARCHAR(255),
    approver3_id VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    CONSTRAINT chk_approval_action_type CHECK (action_type IN ('D3', 'D4', 'D5')),
    CONSTRAINT chk_approval_zone_id CHECK (zone_id IN ('Z1', 'Z2', 'Z3', 'Z4')),
    CONSTRAINT chk_approval_status CHECK (status IN ('pending', 'approved', 'rejected', 'expired'))
);

CREATE INDEX idx_approval_requests_status ON approval_requests(status, expires_at);
CREATE INDEX idx_approval_requests_zone ON approval_requests(zone_id);

-- Keepalive sessions table
CREATE TABLE IF NOT EXISTS keepalive_sessions (
    action_id VARCHAR(255) PRIMARY KEY,
    approver1_last_keepalive TIMESTAMP,
    approver2_last_keepalive TIMESTAMP,
    approver3_last_keepalive TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_keepalive_last_keepalive ON keepalive_sessions(approver1_last_keepalive, approver2_last_keepalive);

-- Device trust scores table
CREATE TABLE IF NOT EXISTS device_trust_scores (
    device_id_hash VARCHAR(255) PRIMARY KEY,
    accuracy_score DECIMAL(3,2) DEFAULT 0.5,
    frequency_score DECIMAL(3,2),
    integrity_score DECIMAL(3,2),
    last_corroboration_score DECIMAL(3,2),
    trust_score DECIMAL(3,2),
    report_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Device report history table
CREATE TABLE IF NOT EXISTS device_report_history (
    id VARCHAR(255) PRIMARY KEY,
    device_id_hash VARCHAR(255) NOT NULL,
    report_id VARCHAR(255) NOT NULL,
    actual_outcome VARCHAR(20),
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_outcome CHECK (actual_outcome IN ('true_positive', 'true_negative', 'false_positive', 'false_negative'))
);

CREATE INDEX idx_device_report_history_device_outcome ON device_report_history(device_id_hash, actual_outcome);

-- CAP messages table
CREATE TABLE IF NOT EXISTS cap_messages (
    id VARCHAR(255) PRIMARY KEY,
    identifier VARCHAR(255) UNIQUE NOT NULL,
    sender VARCHAR(255) NOT NULL,
    sent TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL,
    msg_type VARCHAR(20) NOT NULL,
    scope VARCHAR(20) NOT NULL,
    info JSONB NOT NULL,
    area JSONB NOT NULL,
    expires TIMESTAMP NOT NULL,
    signature TEXT,
    published_channels TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_cap_status CHECK (status IN ('Actual', 'Test', 'Exercise')),
    CONSTRAINT chk_cap_msg_type CHECK (msg_type IN ('Alert', 'Update', 'Cancel')),
    CONSTRAINT chk_cap_scope CHECK (scope IN ('Public', 'Restricted'))
);

CREATE INDEX idx_cap_messages_expires ON cap_messages(expires);
CREATE INDEX idx_cap_messages_identifier ON cap_messages(identifier);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(255) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    operator_id_hash VARCHAR(255) NOT NULL,
    target VARCHAR(255),
    result VARCHAR(20) NOT NULL,
    reason TEXT,
    hash VARCHAR(255) NOT NULL,
    sealed BOOLEAN DEFAULT FALSE,
    sealed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_audit_result CHECK (result IN ('success', 'failure'))
);

CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_operator ON audit_logs(operator_id_hash);
CREATE INDEX idx_audit_logs_action_type ON audit_logs(action_type);

