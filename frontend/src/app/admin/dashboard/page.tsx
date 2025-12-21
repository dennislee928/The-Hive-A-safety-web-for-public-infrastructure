'use client';

import { useEffect, useState } from 'react';
import { dashboardApi, DashboardData } from '@/lib/api/dashboard';
import { Activity, AlertTriangle, TrendingUp, Shield } from 'lucide-react';

export default function AdminDashboard() {
  const [zoneId, setZoneId] = useState('Z1');
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const zones = ['Z1', 'Z2', 'Z3', 'Z4'];

  useEffect(() => {
    fetchDashboardData();
  }, [zoneId]);

  const fetchDashboardData = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await dashboardApi.getDashboardData(zoneId);
      setDashboardData(data);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch dashboard data');
    } finally {
      setLoading(false);
    }
  };

  const getComplexityLevel = (xTotal: number) => {
    if (xTotal < 0.3) return { label: '低', color: 'text-green-600 bg-green-100' };
    if (xTotal < 0.6) return { label: '中', color: 'text-yellow-600 bg-yellow-100' };
    if (xTotal < 0.8) return { label: '高', color: 'text-orange-600 bg-orange-100' };
    return { label: '極高', color: 'text-red-600 bg-red-100' };
  };

  if (loading) {
    return (
      <div className="p-8 text-center">
        <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p className="mt-4 text-gray-600">載入中...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8">
        <div className="bg-red-50 border border-red-200 rounded-md p-4 text-red-800">
          {error}
        </div>
      </div>
    );
  }

  if (!dashboardData) {
    return null;
  }

  const complexityLevel = dashboardData.complexity_metrics
    ? getComplexityLevel(dashboardData.complexity_metrics.complexity_total)
    : null;

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">儀表板</h1>
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
      </div>

      {/* Current State */}
      <div className="mb-6 bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">當前狀態</h2>
        <div className="flex items-center gap-4">
          {dashboardData.decision_state ? (
            <div className="px-4 py-2 bg-blue-100 text-blue-800 rounded-md font-semibold">
              {dashboardData.decision_state.current_state}
            </div>
          ) : (
            <div className="px-4 py-2 bg-gray-100 text-gray-600 rounded-md font-semibold">
              無活動狀態
            </div>
          )}
        </div>
      </div>

      {/* Complexity Metrics */}
      {dashboardData.complexity_metrics && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">信號來源 (x_s)</h3>
              <Activity className="w-5 h-5 text-blue-600" />
            </div>
            <p className="text-2xl font-bold text-gray-900">{dashboardData.complexity_metrics.signal_sources}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">決策深度 (x_d)</h3>
              <TrendingUp className="w-5 h-5 text-blue-600" />
            </div>
            <p className="text-2xl font-bold text-gray-900">{dashboardData.complexity_metrics.decision_depth}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">情境狀態 (x_c)</h3>
              <Activity className="w-5 h-5 text-blue-600" />
            </div>
            <p className="text-2xl font-bold text-gray-900">{dashboardData.complexity_metrics.context_states}</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm font-medium text-gray-600">總複雜度 (x_total)</h3>
              <AlertTriangle className="w-5 h-5 text-orange-600" />
            </div>
            <p className="text-2xl font-bold text-gray-900">
              {dashboardData.complexity_metrics.complexity_total.toFixed(3)}
            </p>
            {complexityLevel && (
              <span className={`inline-block mt-2 px-2 py-1 rounded text-xs font-semibold ${complexityLevel.color}`}>
                {complexityLevel.label}
              </span>
            )}
          </div>
        </div>
      )}

      {/* Ethical Primes */}
      {dashboardData.ethical_primes && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-600 mb-2">FN-prime</h3>
            <p className="text-2xl font-bold text-gray-900">
              {dashboardData.ethical_primes.fn_prime.toFixed(3)}
            </p>
            <p className="text-xs text-gray-500 mt-1">目標: &lt; 0.2</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-600 mb-2">FP-prime</h3>
            <p className="text-2xl font-bold text-gray-900">
              {dashboardData.ethical_primes.fp_prime.toFixed(3)}
            </p>
            <p className="text-xs text-gray-500 mt-1">目標: &lt; 0.15</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-600 mb-2">Bias-prime</h3>
            <p className="text-2xl font-bold text-gray-900">
              {dashboardData.ethical_primes.bias_prime.toFixed(3)}
            </p>
            <p className="text-xs text-gray-500 mt-1">目標: &lt; 0.1</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-sm font-medium text-gray-600 mb-2">Integrity-prime</h3>
            <p className="text-2xl font-bold text-gray-900">
              {dashboardData.ethical_primes.integrity_prime.toFixed(3)}
            </p>
            <p className="text-xs text-gray-500 mt-1">目標: &lt; 0.05</p>
          </div>
        </div>
      )}
    </div>
  );
}

