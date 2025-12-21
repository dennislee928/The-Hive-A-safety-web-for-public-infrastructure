import { apiClient, ApiResponse } from '../api';

export interface DashboardData {
  zone_id: string;
  current_state: string;
  complexity: {
    x_signal: number;
    x_depth: number;
    x_context: number;
    x_total: number;
  };
  ethical_primes: {
    fn_prime: number;
    fp_prime: number;
    bias_prime: number;
    integrity_prime: number;
  };
  recent_signals: any[];
  active_mitigations: any[];
}

export const dashboardApi = {
  // Get dashboard data for a zone
  getDashboardData: async (zoneId: string): Promise<DashboardData> => {
    const response = await apiClient.get<ApiResponse<DashboardData>>(`/dashboard/zones/${zoneId}`);
    return response.data.data!;
  },
};

