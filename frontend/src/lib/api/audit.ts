import { apiClient, ApiResponse } from '../api';

export interface AuditLog {
  id: string;
  timestamp: string;
  operation_type: string;
  operator_id: string;
  target_type?: string;
  target_id?: string;
  action: string;
  result: string;
  reason?: string;
  metadata?: string;
}

export interface EvidenceRecord {
  id: string;
  evidence_type: string;
  related_id: string;
  zone_id?: string;
  archived_at: string;
  archived_by: string;
  retention_until: string;
  sealed: boolean;
}

export const auditApi = {
  // Get audit logs
  getAuditLogs: async (filters?: {
    operation_type?: string;
    operator_id?: string;
    target_type?: string;
    target_id?: string;
    action?: string;
    result?: string;
    start_time?: string;
    end_time?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ logs: AuditLog[]; count: number }> => {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, String(value));
        }
      });
    }
    
    const response = await apiClient.get<ApiResponse<{ logs: AuditLog[]; count: number }>>(
      `/audit/logs?${params.toString()}`
    );
    return response.data.data || { logs: [], count: 0 };
  },

  // Verify integrity
  verifyIntegrity: async (startTime?: string, endTime?: string): Promise<any> => {
    const params = new URLSearchParams();
    if (startTime) params.append('start_time', startTime);
    if (endTime) params.append('end_time', endTime);
    
    const response = await apiClient.get<ApiResponse<any>>(
      `/audit/verify-integrity?${params.toString()}`
    );
    return response.data.data;
  },

  // Get evidence
  getEvidence: async (evidenceId: string): Promise<EvidenceRecord> => {
    const response = await apiClient.get<ApiResponse<EvidenceRecord>>(`/audit/evidence/${evidenceId}`);
    return response.data.data!;
  },

  // List evidence
  listEvidence: async (filters?: {
    evidence_type?: string;
    related_id?: string;
    zone_id?: string;
    start_time?: string;
    end_time?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ evidence: EvidenceRecord[]; count: number }> => {
    const params = new URLSearchParams();
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, String(value));
        }
      });
    }
    
    const response = await apiClient.get<ApiResponse<{ evidence: EvidenceRecord[]; count: number }>>(
      `/audit/evidence?${params.toString()}`
    );
    return response.data.data || { evidence: [], count: 0 };
  },

  // Archive evidence
  archiveEvidence: async (data: {
    evidence_type: string;
    related_id: string;
    zone_id?: string;
    snapshot: string;
    retention_period?: number;
  }): Promise<EvidenceRecord> => {
    const response = await apiClient.post<ApiResponse<EvidenceRecord>>('/audit/evidence/archive', data);
    return response.data.data!;
  },
};

