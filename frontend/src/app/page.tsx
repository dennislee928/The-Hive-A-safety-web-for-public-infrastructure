import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="text-center space-y-8 p-8">
        <h1 className="text-5xl font-bold text-gray-900">ERH Safety System</h1>
        <p className="text-xl text-gray-600">選擇您的入口</p>
        <div className="flex gap-6 justify-center mt-8">
          <Link
            href="/client"
            className="px-8 py-4 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-lg font-semibold"
          >
            客戶端
          </Link>
          <Link
            href="/admin"
            className="px-8 py-4 bg-gray-800 text-white rounded-lg hover:bg-gray-900 transition-colors text-lg font-semibold"
          >
            管理員後台
          </Link>
        </div>
      </div>
    </div>
  );
}

