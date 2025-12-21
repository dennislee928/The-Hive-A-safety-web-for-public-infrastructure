'use client';

import { useEffect, useState } from 'react';
import { auditApi, AuditLog, EvidenceRecord } from '@/lib/api/audit';
import { FileText, Shield, Search } from 'lucide-react';

export default function AuditPage() {
  const [activeTab, setActiveTab] = useState<'logs' | 'evidence'>('logs');
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [evidence, setEvidence] = useState<EvidenceRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState({
    operation_type: '',
    start_time: '',
    end_time: '',
    limit: 100,
  });

  useEffect(() => {
    if (activeTab === 'logs') {
      fetchLogs();
    } else {
      fetchEvidence();
    }
  }, [activeTab, filters]);

  const fetchLogs = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await auditApi.getAuditLogs(filters);
      setLogs(result.logs);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch audit logs');
    } finally {
      setLoading(false);
    }
  };

  const fetchEvidence = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await auditApi.listEvidence({
        start_time: filters.start_time,
        end_time: filters.end_time,
        limit: filters.limit,
      });
      setEvidence(result.evidence);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch evidence');
    } finally {
      setLoading(false);
    }
  };

  const getResultColor = (result: string) => {
    switch (result) {
      case 'success':
        return 'text-green-600 bg-green-100';
      case 'failure':
        return 'text-red-600 bg-red-100';
      case 'error':
        return 'text-red-800 bg-red-200';
      default:
        return 'text-gray-600 bg-gray-100';
    }
  };

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">審計日誌</h1>

      {/* Tabs */}
      <div className="mb-6 border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('logs')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'logs'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <FileText className="w-5 h-5 inline mr-2" />
            審計日誌
          </button>
          <button
            onClick={() => setActiveTab('evidence')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'evidence'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            <Shield className="w-5 h-5 inline mr-2" />
            證據封存
          </button>
        </nav>
      </div>

      {/* Filters */}
      <div className="mb-6 bg-white rounded-lg shadow p-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {activeTab === 'logs' && (
            <select
              value={filters.operation_type}
              onChange={(e) => setFilters({ ...filters, operation_type: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">所有操作類型</option>
              <option value="data_access">資料存取</option>
              <option value="decision_transition">決策轉換</option>
              <option value="system_config">系統配置</option>
            </select>
          )}
          <input
            type="datetime-local"
            value={filters.start_time}
            onChange={(e) => setFilters({ ...filters, start_time: e.target.value })}
            placeholder="開始時間"
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            type="datetime-local"
            value={filters.end_time}
            onChange={(e) => setFilters({ ...filters, end_time: e.target.value })}
            placeholder="結束時間"
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={() => activeTab === 'logs' ? fetchLogs() : fetchEvidence()}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 flex items-center justify-center gap-2"
          >
            <Search className="w-4 h-4" />
            搜尋
          </button>
        </div>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-red-800">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">載入中...</p>
        </div>
      ) : activeTab === 'logs' ? (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  時間
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  操作類型
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  動作
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  目標
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  結果
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {logs.map((log) => (
                <tr key={log.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {new Date(log.timestamp).toLocaleString('zh-TW')}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {log.operation_type}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {log.action}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {log.target_type && log.target_id
                      ? `${log.target_type}: ${log.target_id}`
                      : '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getResultColor(
                        log.result
                      )}`}
                    >
                      {log.result}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {logs.length === 0 && (
            <div className="text-center py-12 text-gray-500">無審計日誌記錄</div>
          )}
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  封存時間
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  證據類型
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  相關 ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  區域
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  狀態
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {evidence.map((record) => (
                <tr key={record.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {new Date(record.archived_at).toLocaleString('zh-TW')}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {record.evidence_type}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {record.related_id}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {record.zone_id || '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        record.sealed
                          ? 'text-green-800 bg-green-100'
                          : 'text-gray-800 bg-gray-100'
                      }`}
                    >
                      {record.sealed ? '已密封' : '未密封'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {evidence.length === 0 && (
            <div className="text-center py-12 text-gray-500">無證據記錄</div>
          )}
        </div>
      )}
    </div>
  );
}

