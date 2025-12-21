'use client';

import { useEffect, useState } from 'react';
import { decisionApi, DecisionState } from '@/lib/api/decision';
import { Plus, ArrowRight } from 'lucide-react';

export default function DecisionsPage() {
  const [zoneId, setZoneId] = useState('Z1');
  const [decisionState, setDecisionState] = useState<DecisionState | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showCreateAlert, setShowCreateAlert] = useState(false);
  const [reason, setReason] = useState('');

  const zones = ['Z1', 'Z2', 'Z3', 'Z4'];

  useEffect(() => {
    fetchDecisionState();
  }, [zoneId]);

  const fetchDecisionState = async () => {
    setLoading(true);
    setError(null);
    try {
      const state = await decisionApi.getLatestState(zoneId);
      setDecisionState(state);
    } catch (err: any) {
      if (err.response?.status !== 404) {
        setError(err.message || 'Failed to fetch decision state');
      } else {
        setDecisionState(null);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCreatePreAlert = async () => {
    if (!reason.trim()) {
      alert('請輸入原因');
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const newState = await decisionApi.createPreAlert(zoneId, reason);
      setDecisionState(newState);
      setShowCreateAlert(false);
      setReason('');
    } catch (err: any) {
      setError(err.message || 'Failed to create pre-alert');
    } finally {
      setLoading(false);
    }
  };

  const handleTransition = async (targetState: string) => {
    if (!reason.trim()) {
      alert('請輸入原因');
      return;
    }

    if (!decisionState) return;

    setLoading(true);
    setError(null);
    try {
      const newState = await decisionApi.transitionState(decisionState.id, {
        target_state: targetState,
        reason: reason,
      });
      setDecisionState(newState);
      setReason('');
    } catch (err: any) {
      setError(err.message || 'Failed to transition state');
    } finally {
      setLoading(false);
    }
  };

  if (loading && !decisionState) {
    return (
      <div className="p-8 text-center">
        <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p className="mt-4 text-gray-600">載入中...</p>
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">決策管理</h1>
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

      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-red-800">
          {error}
        </div>
      )}

      {decisionState ? (
        <div className="bg-white rounded-lg shadow p-6">
          <div className="mb-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">當前決策狀態</h2>
            <div className="flex items-center gap-4">
              <div className="px-4 py-2 bg-blue-100 text-blue-800 rounded-md font-semibold text-lg">
                {decisionState.current_state}
              </div>
              {decisionState.previous_state && (
                <>
                  <ArrowRight className="w-5 h-5 text-gray-400" />
                  <span className="text-gray-600">從 {decisionState.previous_state}</span>
                </>
              )}
            </div>
            {decisionState.reason && (
              <p className="mt-2 text-gray-600">原因：{decisionState.reason}</p>
            )}
            <p className="mt-2 text-sm text-gray-500">
              更新時間：{new Date(decisionState.updated_at).toLocaleString('zh-TW')}
            </p>
          </div>

          <div className="mt-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">狀態轉換</h3>
            <div className="mb-4">
              <textarea
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder="輸入轉換原因..."
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={3}
              />
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => handleTransition('D1')}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                轉換至 D1
              </button>
              <button
                onClick={() => handleTransition('D2')}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                轉換至 D2
              </button>
              <button
                onClick={() => handleTransition('D3')}
                className="px-4 py-2 bg-orange-600 text-white rounded-md hover:bg-orange-700"
              >
                轉換至 D3
              </button>
              <button
                onClick={() => handleTransition('D6')}
                className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
              >
                降級 (D6)
              </button>
            </div>
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">無活動決策狀態</h2>
          <button
            onClick={() => setShowCreateAlert(true)}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <Plus className="w-5 h-5" />
            建立預警 (D0)
          </button>

          {showCreateAlert && (
            <div className="mt-4 p-4 bg-gray-50 rounded-md">
              <textarea
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder="輸入預警原因..."
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 mb-4"
                rows={3}
              />
              <div className="flex gap-2">
                <button
                  onClick={handleCreatePreAlert}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                  確認建立
                </button>
                <button
                  onClick={() => {
                    setShowCreateAlert(false);
                    setReason('');
                  }}
                  className="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400"
                >
                  取消
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

