'use client';

import { useEffect, useState } from 'react';
import { capApi, CAPMessage } from '@/lib/api/cap';
import { AlertCircle, CheckCircle, Info, AlertTriangle } from 'lucide-react';

export default function ClientPage() {
  const [zoneId, setZoneId] = useState('Z1');
  const [capMessages, setCapMessages] = useState<CAPMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const zones = ['Z1', 'Z2', 'Z3', 'Z4'];

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

  useEffect(() => {
    fetchCAPMessages();
    // Refresh every 30 seconds
    const interval = setInterval(fetchCAPMessages, 30000);
    return () => clearInterval(interval);
  }, [zoneId]);

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

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'Extreme':
        return 'border-red-500 bg-red-50';
      case 'Severe':
        return 'border-orange-500 bg-orange-50';
      default:
        return 'border-blue-500 bg-blue-50';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">ERH Safety - 客戶端</h1>
            <div className="flex items-center gap-4">
              <label className="text-sm font-medium text-gray-700">區域：</label>
              <select
                value={zoneId}
                onChange={(e) => setZoneId(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {zones.map((zone) => (
                  <option key={zone} value={zone}>
                    {zone}
                  </option>
                ))}
              </select>
              <button
                onClick={fetchCAPMessages}
                disabled={loading}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? '載入中...' : '重新載入'}
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-red-800">
            {error}
          </div>
        )}

        {loading && capMessages.length === 0 ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">載入中...</p>
          </div>
        ) : capMessages.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg shadow">
            <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
            <h2 className="text-2xl font-semibold text-gray-900 mb-2">目前無活動警示</h2>
            <p className="text-gray-600">您所在的區域 {zoneId} 目前是安全的</p>
          </div>
        ) : (
          <div className="space-y-4">
            {capMessages.map((message, index) => {
              const info = message.info[0];
              return (
                <div
                  key={message.identifier}
                  className={`border-2 rounded-lg p-6 bg-white shadow-md ${getSeverityColor(info.severity)}`}
                >
                  <div className="flex items-start gap-4">
                    {getSeverityIcon(info.severity)}
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-2">
                        <h2 className="text-xl font-bold text-gray-900">{info.headline}</h2>
                        <span className="px-3 py-1 bg-white rounded-full text-sm font-semibold text-gray-700">
                          {info.severity}
                        </span>
                      </div>
                      <p className="text-gray-700 mb-4">{info.description}</p>
                      {info.instruction && (
                        <div className="bg-white/80 rounded p-3 mb-4">
                          <p className="font-semibold text-gray-900 mb-1">指示：</p>
                          <p className="text-gray-700">{info.instruction}</p>
                        </div>
                      )}
                      <div className="flex items-center gap-4 text-sm text-gray-600">
                        <span>發布時間：{new Date(message.sent).toLocaleString('zh-TW')}</span>
                        {info.expires && (
                          <span>有效期至：{new Date(info.expires).toLocaleString('zh-TW')}</span>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </main>
    </div>
  );
}

