import { apiClient, ApiResponse } from '../api';

export interface DecisionState {
  id: string;
  zone_id: string;
  current_state: string;
  previous_state?: string;
  reason?: string;
  created_at: string;
  updated_at: string;
}

export interface DecisionTransitionRequest {
  target_state: string;
  reason: string;
}

export const decisionApi = {
  // Get latest decision state for a zone
  getLatestState: async (zoneId: string): Promise<DecisionState> => {
    const response = await apiClient.get<ApiResponse<DecisionState>>(`/operator/zones/${zoneId}/state`);
    return response.data.data!;
  },

  // Create pre-alert (D0)
  createPreAlert: async (zoneId: string, reason: string): Promise<DecisionState> => {
    const response = await apiClient.post<ApiResponse<DecisionState>>(
      `/operator/decisions/${zoneId}/d0`,
      { reason }
    );
    return response.data.data!;
  },

  // Transition decision state
  transitionState: async (decisionId: string, data: DecisionTransitionRequest): Promise<DecisionState> => {
    const response = await apiClient.post<ApiResponse<DecisionState>>(
      `/operator/decisions/${decisionId}/transition`,
      data
    );
    return response.data.data!;
  },
};

