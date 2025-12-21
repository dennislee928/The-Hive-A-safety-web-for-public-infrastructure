'use client';

import { useEffect, useState } from 'react';
import { erhApi, ERHStatus } from '@/lib/api/erh';
import { Activity, TrendingUp, AlertTriangle, Shield } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

export default function ERHPage() {
  const [zoneId, setZoneId] = useState('Z1');
  const [erhStatus, setErhStatus] = useState<ERHStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const zones = ['Z1', 'Z2', 'Z3', 'Z4'];

  useEffect(() => {
    fetchERHStatus();
    const interval = setInterval(fetchERHStatus, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, [zoneId]);

  const fetchERHStatus = async () => {
    setLoading(true);
    setError(null);
    try {
      const status = await erhApi.getERHStatus(zoneId);
      setErhStatus(status);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch ERH status');
    } finally {
      setLoading(false);
    }
  };

  if (loading && !erhStatus) {
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

  if (!erhStatus) {
    return null;
  }

  const complexityData = [
    { name: '信號來源', value: erhStatus.complexity.signal_sources },
    { name: '決策深度', value: erhStatus.complexity.decision_depth },
    { name: '情境狀態', value: erhStatus.complexity.context_states },
  ];

  const primesData = [
    { name: 'FN-prime', value: erhStatus.ethical_primes.fn_prime, target: 0.2 },
    { name: 'FP-prime', value: erhStatus.ethical_primes.fp_prime, target: 0.15 },
    { name: 'Bias-prime', value: erhStatus.ethical_primes.bias_prime, target: 0.1 },
    { name: 'Integrity-prime', value: erhStatus.ethical_primes.integrity_prime, target: 0.05 },
  ];

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">ERH 監控</h1>
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

      {/* Complexity Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600">信號來源</h3>
            <Activity className="w-5 h-5 text-blue-600" />
          </div>
          <p className="text-2xl font-bold text-gray-900">{erhStatus.complexity.signal_sources}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600">決策深度</h3>
            <TrendingUp className="w-5 h-5 text-blue-600" />
          </div>
          <p className="text-2xl font-bold text-gray-900">{erhStatus.complexity.decision_depth}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600">情境狀態</h3>
            <Activity className="w-5 h-5 text-blue-600" />
          </div>
          <p className="text-2xl font-bold text-gray-900">{erhStatus.complexity.context_states}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600">總複雜度</h3>
            <AlertTriangle className="w-5 h-5 text-orange-600" />
          </div>
          <p className="text-2xl font-bold text-gray-900">
            {erhStatus.complexity.complexity_total.toFixed(3)}
          </p>
        </div>
      </div>

      {/* Ethical Primes */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">倫理質數</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {primesData.map((prime) => {
            const isOverTarget = prime.value > prime.target;
            return (
              <div key={prime.name} className="p-4 bg-gray-50 rounded-lg">
                <h3 className="text-sm font-medium text-gray-600 mb-2">{prime.name}</h3>
                <p className={`text-2xl font-bold ${isOverTarget ? 'text-red-600' : 'text-green-600'}`}>
                  {prime.value.toFixed(3)}
                </p>
                <p className="text-xs text-gray-500 mt-1">目標: &lt; {prime.target}</p>
              </div>
            );
          })}
        </div>
      </div>

      {/* Breakpoints */}
      {erhStatus.breakpoints && erhStatus.breakpoints.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">檢測到的斷點</h2>
          <div className="space-y-2">
            {erhStatus.breakpoints.map((breakpoint: any, index: number) => (
              <div key={index} className="p-3 bg-red-50 border border-red-200 rounded-md">
                <p className="font-semibold text-red-800">{breakpoint.type}: {breakpoint.value}</p>
                <p className="text-sm text-red-600">{breakpoint.description}</p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Active Mitigations */}
      {erhStatus.active_mitigations && erhStatus.active_mitigations.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">活動中的緩解措施</h2>
          <div className="space-y-2">
            {erhStatus.active_mitigations.map((mitigation: any, index: number) => (
              <div key={index} className="p-3 bg-blue-50 border border-blue-200 rounded-md">
                <p className="font-semibold text-blue-800">{mitigation.measure_type}</p>
                <p className="text-sm text-blue-600">{mitigation.reason}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

