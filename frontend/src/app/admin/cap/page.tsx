'use client';

import { useEffect, useState } from 'react';
import { capApi, CAPMessage } from '@/lib/api/cap';
import { Plus, AlertCircle, AlertTriangle, Info } from 'lucide-react';

export default function CAPPage() {
  const [zoneId, setZoneId] = useState('Z1');
  const [capMessages, setCapMessages] = useState<CAPMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState({
    zone_id: 'Z1',
    message_type: 'Alert',
    severity: 'Severe',
    urgency: 'Immediate',
    headline: '',
    description: '',
    instruction: '',
  });

  const zones = ['Z1', 'Z2', 'Z3', 'Z4'];

  useEffect(() => {
    fetchCAPMessages();
  }, [zoneId]);

  const fetchCAPMessages = async () => {
    setLoading(true);
    setError(null);
    try {
      const messages = await capApi.getCAPMessagesByZone(zoneId);
      setCapMessages(messages);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch CAP messages');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCAP = async () => {
    if (!formData.headline.trim() || !formData.description.trim()) {
      alert('請填寫標題和描述');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      await capApi.generateAndPublish(formData);
      setShowCreateForm(false);
      setFormData({
        zone_id: 'Z1',
        message_type: 'Alert',
        severity: 'Severe',
        urgency: 'Immediate',
        headline: '',
        description: '',
        instruction: '',
      });
      fetchCAPMessages();
    } catch (err: any) {
      setError(err.message || 'Failed to create CAP message');
    } finally {
      setLoading(false);
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case 'Extreme':
        return <AlertCircle className="w-6 h-6 text-red-600" />;
      case 'Severe':
        return <AlertTriangle className="w-6 h-6 text-orange-600" />;
      default:
        return <Info className="w-6 h-6 text-blue-600" />;
    }
  };

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">CAP 訊息管理</h1>
        <div className="flex items-center gap-4">
          <select
            value={zoneId}
            onChange={(e) => setZoneId(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {zones.map((zone) => (
              <option key={zone} value={zone}>
                區域 {zone}
              </option>
            ))}
          </select>
          <button
            onClick={() => setShowCreateForm(true)}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <Plus className="w-5 h-5" />
            建立 CAP 訊息
          </button>
        </div>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-red-800">
          {error}
        </div>
      )}

      {showCreateForm && (
        <div className="mb-6 bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">建立 CAP 訊息</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">區域</label>
              <select
                value={formData.zone_id}
                onChange={(e) => setFormData({ ...formData, zone_id: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {zones.map((zone) => (
                  <option key={zone} value={zone}>
                    {zone}
                  </option>
                ))}
              </select>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">嚴重程度</label>
                <select
                  value={formData.severity}
                  onChange={(e) => setFormData({ ...formData, severity: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="Extreme">Extreme</option>
                  <option value="Severe">Severe</option>
                  <option value="Moderate">Moderate</option>
                  <option value="Minor">Minor</option>
                  <option value="Unknown">Unknown</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">緊急程度</label>
                <select
                  value={formData.urgency}
                  onChange={(e) => setFormData({ ...formData, urgency: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="Immediate">Immediate</option>
                  <option value="Expected">Expected</option>
                  <option value="Future">Future</option>
                  <option value="Past">Past</option>
                  <option value="Unknown">Unknown</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">標題 *</label>
              <input
                type="text"
                value={formData.headline}
                onChange={(e) => setFormData({ ...formData, headline: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">描述 *</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">指示（選填）</label>
              <textarea
                value={formData.instruction}
                onChange={(e) => setFormData({ ...formData, instruction: e.target.value })}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleCreateCAP}
                disabled={loading}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? '建立中...' : '建立並發布'}
              </button>
              <button
                onClick={() => {
                  setShowCreateForm(false);
                  setFormData({
                    zone_id: 'Z1',
                    message_type: 'Alert',
                    severity: 'Severe',
                    urgency: 'Immediate',
                    headline: '',
                    description: '',
                    instruction: '',
                  });
                }}
                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
              >
                取消
              </button>
            </div>
          </div>
        </div>
      )}

      {loading && capMessages.length === 0 ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">載入中...</p>
        </div>
      ) : capMessages.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-lg shadow">
          <p className="text-gray-600">區域 {zoneId} 目前沒有 CAP 訊息</p>
        </div>
      ) : (
        <div className="space-y-4">
          {capMessages.map((message) => {
            const info = message.info[0];
            return (
              <div key={message.identifier} className="bg-white rounded-lg shadow p-6">
                <div className="flex items-start gap-4">
                  {getSeverityIcon(info.severity)}
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h2 className="text-xl font-bold text-gray-900">{info.headline}</h2>
                      <span className="px-3 py-1 bg-gray-100 rounded-full text-sm font-semibold text-gray-700">
                        {info.severity} / {info.urgency}
                      </span>
                    </div>
                    <p className="text-gray-700 mb-4">{info.description}</p>
                    {info.instruction && (
                      <div className="bg-blue-50 rounded p-3 mb-4">
                        <p className="font-semibold text-gray-900 mb-1">指示：</p>
                        <p className="text-gray-700">{info.instruction}</p>
                      </div>
                    )}
                    <div className="flex items-center gap-4 text-sm text-gray-600">
                      <span>訊息 ID: {message.identifier}</span>
                      <span>發布時間: {new Date(message.sent).toLocaleString('zh-TW')}</span>
                      {info.expires && (
                        <span>有效期至: {new Date(info.expires).toLocaleString('zh-TW')}</span>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

