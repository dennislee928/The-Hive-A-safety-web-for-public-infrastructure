import { apiClient, ApiResponse } from '../api';

export interface DashboardData {
  zone_id: string;
  decision_state: {
    id: string;
    zone_id: string;
    current_state: string;
    previous_state?: string;
    reason?: string;
    created_at: string;
    updated_at: string;
    signal_count?: number;
    context_states?: number;
  } | null;
  complexity_metrics: {
    signal_sources: number;
    decision_depth: number;
    context_states: number;
    complexity_total: number;
    complexity_level?: string;
  } | null;
  ethical_primes: {
    fn_prime: number;
    fp_prime: number;
    bias_prime: number;
    integrity_prime: number;
  } | null;
}

export const dashboardApi = {
  // Get dashboard data for a zone
  getDashboardData: async (zoneId: string): Promise<DashboardData> => {
    const response = await apiClient.get<ApiResponse<DashboardData>>(`/dashboard/zones/${zoneId}`);
    return response.data.data!;
  },
};

