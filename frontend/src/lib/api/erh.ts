import { apiClient, ApiResponse } from '../api';

export interface ERHStatus {
  zone_id: string;
  complexity: {
    signal_sources: number;
    decision_depth: number;
    context_states: number;
    complexity_total: number;
  };
  ethical_primes: {
    fn_prime: number;
    fp_prime: number;
    bias_prime: number;
    integrity_prime: number;
  };
  breakpoints: any[];
  active_mitigations: any[];
}

export interface ERHMetricsHistory {
  id: string;
  zone_id: string;
  x_total: number;
  x_signal: number;
  x_depth: number;
  x_context: number;
  fn_prime: number;
  fp_prime: number;
  bias_prime: number;
  integrity_prime: number;
  timestamp: string;
}

export interface ERHReport {
  report_id: string;
  report_type: string;
  generated_at: string;
  summary: any;
  complexity: any;
  ethical_primes: any;
  breakpoints: any[];
  trends: any;
  recommendations: string[];
}

export const erhApi = {
  // Get ERH status for a zone
  getERHStatus: async (zoneId: string): Promise<ERHStatus> => {
    const response = await apiClient.get<ApiResponse<ERHStatus>>(`/erh/status/${zoneId}`);
    return response.data.data!;
  },

  // Get metrics history
  getMetricsHistory: async (
    zoneId: string,
    startTime?: string,
    endTime?: string
  ): Promise<ERHMetricsHistory[]> => {
    const params = new URLSearchParams();
    if (startTime) params.append('start_time', startTime);
    if (endTime) params.append('end_time', endTime);
    
    const response = await apiClient.get<ApiResponse<ERHMetricsHistory[]>>(
      `/erh/metrics/${zoneId}/history?${params.toString()}`
    );
    return response.data.data || [];
  },

  // Get metrics trends
  getMetricsTrends: async (zoneId: string, duration: string = '24h'): Promise<any> => {
    const response = await apiClient.get<ApiResponse<any>>(
      `/erh/metrics/${zoneId}/trends?duration=${duration}`
    );
    return response.data.data;
  },

  // Generate report
  generateReport: async (zoneId: string, reportType: 'daily' | 'weekly' | 'monthly'): Promise<ERHReport> => {
    const response = await apiClient.get<ApiResponse<ERHReport>>(
      `/erh/reports/${zoneId}/${reportType}`
    );
    return response.data.data!;
  },

  // Activate mitigation
  activateMitigation: async (data: {
    measure_type: string;
    trigger_type: string;
    trigger_condition: string;
    reason: string;
  }): Promise<any> => {
    const response = await apiClient.post<ApiResponse<any>>('/erh/mitigations', data);
    return response.data.data;
  },
};

