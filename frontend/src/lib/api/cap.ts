import { apiClient, ApiResponse } from '../api';

export interface CAPMessage {
  identifier: string;
  sender: string;
  sent: string;
  status: string;
  msgType: string;
  scope: string;
  info: CAPInfo[];
  areas: CAPArea[];
}

export interface CAPInfo {
  language: string;
  category: string[];
  event: string;
  urgency: string;
  severity: string;
  certainty: string;
  headline: string;
  description: string;
  instruction?: string;
  contact?: string;
  expires?: string;
}

export interface CAPArea {
  zone_id: string;
  zone_type: string;
  time_window?: {
    start: string;
    end: string;
  };
}

export const capApi = {
  // Get CAP message by identifier
  getCAPMessage: async (identifier: string): Promise<CAPMessage> => {
    const response = await apiClient.get<ApiResponse<CAPMessage>>(`/cap/${identifier}`);
    return response.data.data!;
  },

  // Get CAP messages by zone
  getCAPMessagesByZone: async (zoneId: string): Promise<CAPMessage[]> => {
    const response = await apiClient.get<ApiResponse<CAPMessage[]>>(`/cap/zone/${zoneId}`);
    return response.data.data || [];
  },

  // Generate and publish CAP message (admin only)
  generateAndPublish: async (data: {
    zone_id: string;
    message_type: string;
    severity: string;
    urgency: string;
    headline: string;
    description: string;
    instruction?: string;
  }): Promise<CAPMessage> => {
    const response = await apiClient.post<ApiResponse<CAPMessage>>('/cap/generate', data);
    return response.data.data!;
  },
};

