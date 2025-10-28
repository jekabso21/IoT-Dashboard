import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { 
  Activity, 
  Thermometer, 
  Droplets, 
  Wind, 
  AlertTriangle, 
  CheckCircle, 
  Clock,
  TrendingUp,
  TrendingDown,
  Minus,
  Eye,
  Bell,
  BellOff
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { mockDevices, mockAlerts, systemStats } from '../data/mockData';

export default function Dashboard() {
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(new Date()), 60000); // Update every minute instead of every second
    return () => clearInterval(timer);
  }, []);

  // Aggregate data for charts
  const last24Hours = mockDevices[0]?.historicalData.slice(-24) || [];
  const chartData = last24Hours.map(item => ({
    time: new Date(item.timestamp).getHours() + ':00',
    temperature: item.temperature,
    humidity: item.humidity,
    co2: item.co2
  }));

  // Device status distribution
  const statusData = [
    { name: 'Online', value: mockDevices.filter(d => d.status === 'online').length, color: '#4CAF50' },
    { name: 'Warning', value: mockDevices.filter(d => d.status === 'warning').length, color: '#FF9800' },
    { name: 'Offline', value: mockDevices.filter(d => d.status === 'offline').length, color: '#F44336' },
    { name: 'Calibration', value: mockDevices.filter(d => d.status === 'calibration').length, color: '#2196F3' }
  ].filter(item => item.value > 0); // Only show categories with values

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online': return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'warning': return <AlertTriangle className="w-4 h-4 text-yellow-500" />;
      case 'offline': return <AlertTriangle className="w-4 h-4 text-red-500" />;
      case 'calibration': return <Clock className="w-4 h-4 text-blue-500" />;
      default: return <Minus className="w-4 h-4 text-gray-500" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'text-green-500';
      case 'warning': return 'text-yellow-500';
      case 'offline': return 'text-red-500';
      case 'calibration': return 'text-blue-500';
      default: return 'text-gray-500';
    }
  };

  const getTrendIcon = (current: number, previous: number) => {
    if (current > previous) return <TrendingUp className="w-4 h-4 text-green-500" />;
    if (current < previous) return <TrendingDown className="w-4 h-4 text-red-500" />;
    return <Minus className="w-4 h-4 text-gray-500" />;
  };

  const formatTimeAgo = (timestamp: string) => {
    const now = new Date();
    const time = new Date(timestamp);
    const diffInMinutes = Math.floor((now.getTime() - time.getTime()) / (1000 * 60));
    
    if (diffInMinutes < 1) return 'Just now';
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`;
    return `${Math.floor(diffInMinutes / 1440)}d ago`;
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'text-red-500 bg-red-500 bg-opacity-10 border-red-500';
      case 'high': return 'text-orange-500 bg-orange-500 bg-opacity-10 border-orange-500';
      case 'medium': return 'text-yellow-500 bg-yellow-500 bg-opacity-10 border-yellow-500';
      case 'low': return 'text-blue-500 bg-blue-500 bg-opacity-10 border-blue-500';
      default: return 'text-gray-500 bg-gray-500 bg-opacity-10 border-gray-500';
    }
  };

  return (
    <div className="space-y-6 fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-primary">Dashboard</h1>
          <p className="text-secondary mt-1">
            {currentTime.toLocaleDateString('en-US', { 
              weekday: 'long', 
              year: 'numeric', 
              month: 'long', 
              day: 'numeric' 
            })} • {currentTime.toLocaleTimeString('en-US', { 
              hour: '2-digit', 
              minute: '2-digit' 
            })}
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Link to="/devices" className="btn btn-secondary">
            <Activity className="w-4 h-4" />
            Manage Devices
          </Link>
          <button className="btn btn-primary">
            <Eye className="w-4 h-4" />
            View Reports
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-secondary text-sm">Total Devices</p>
              <p className="text-2xl font-bold text-primary">{systemStats.totalDevices}</p>
            </div>
            <div className="w-12 h-12 bg-accent bg-opacity-20 rounded-lg flex items-center justify-center">
              <Activity className="w-6 h-6 text-accent" />
            </div>
          </div>
          <div className="flex items-center gap-2 mt-3">
            {getTrendIcon(systemStats.totalDevices, 5)}
            <span className="text-sm text-secondary">+1 this week</span>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-secondary text-sm">Online Devices</p>
              <p className="text-2xl font-bold text-primary">{systemStats.onlineDevices}</p>
            </div>
            <div className="w-12 h-12 bg-green-500 bg-opacity-20 rounded-lg flex items-center justify-center">
              <CheckCircle className="w-6 h-6 text-green-500" />
            </div>
          </div>
          <div className="flex items-center gap-2 mt-3">
            <span className="text-sm text-green-500">
              {Math.round((systemStats.onlineDevices / systemStats.totalDevices) * 100)}% uptime
            </span>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-secondary text-sm">Active Alerts</p>
              <p className="text-2xl font-bold text-primary">{systemStats.activeAlerts}</p>
            </div>
            <div className="w-12 h-12 bg-red-500 bg-opacity-20 rounded-lg flex items-center justify-center">
              <AlertTriangle className="w-6 h-6 text-red-500" />
            </div>
          </div>
          <div className="flex items-center gap-2 mt-3">
            {systemStats.activeAlerts > 0 ? (
              <span className="text-sm text-red-500">Requires attention</span>
            ) : (
              <span className="text-sm text-green-500">All clear</span>
            )}
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-secondary text-sm">System Health</p>
              <p className="text-2xl font-bold text-primary">{systemStats.systemHealth}%</p>
            </div>
            <div className="w-12 h-12 bg-blue-500 bg-opacity-20 rounded-lg flex items-center justify-center">
              <Activity className="w-6 h-6 text-blue-500" />
            </div>
          </div>
          <div className="flex items-center gap-2 mt-3">
            <div className="w-full bg-elevated rounded-full h-2">
              <div 
                className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                style={{ width: `${systemStats.systemHealth}%` }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Environmental Trends */}
        <div className="lg:col-span-2 card">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-semibold text-primary">24-Hour Trends</h3>
            <div className="flex items-center gap-4 text-sm">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-orange-500"></div>
                <span className="text-secondary">Temperature</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-blue-500"></div>
                <span className="text-secondary">Humidity</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                <span className="text-secondary">CO2</span>
              </div>
            </div>
          </div>
          <div className="h-96 -mx-2 -mb-2">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData} margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#404040" />
                <XAxis 
                  dataKey="time" 
                  stroke="#b3b3b3"
                  tick={{ fontSize: 12 }}
                  interval="preserveStartEnd"
                  axisLine={false}
                  tickLine={false}
                />
                <YAxis 
                  stroke="#b3b3b3"
                  tick={{ fontSize: 12 }}
                  axisLine={false}
                  tickLine={false}
                  width={40}
                />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#2d2d2d', 
                    border: '1px solid #404040',
                    borderRadius: '8px',
                    color: '#ffffff'
                  }}
                />
                <Line 
                  type="monotone" 
                  dataKey="temperature" 
                  stroke="#ff6b35" 
                  strokeWidth={2}
                  dot={false}
                  activeDot={{ r: 4, fill: '#ff6b35' }}
                />
                <Line 
                  type="monotone" 
                  dataKey="humidity" 
                  stroke="#2196f3" 
                  strokeWidth={2}
                  dot={false}
                  activeDot={{ r: 4, fill: '#2196f3' }}
                />
                <Line 
                  type="monotone" 
                  dataKey="co2" 
                  stroke="#4caf50" 
                  strokeWidth={2}
                  dot={false}
                  activeDot={{ r: 4, fill: '#4caf50' }}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Device Status Distribution */}
        <div className="card">
          <h3 className="text-xl font-semibold text-primary mb-6">Device Status</h3>
          <div className="space-y-4">
            {/* Online Devices */}
            <div className="p-4 bg-green-500 bg-opacity-10 border border-green-500 border-opacity-30 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  <CheckCircle className="w-6 h-6 text-green-500" />
                  <span className="font-semibold text-green-400">Online</span>
                </div>
                <span className="text-2xl font-bold text-green-400">
                  {mockDevices.filter(d => d.status === 'online').length}
                </span>
              </div>
              <div className="w-full bg-green-900 bg-opacity-30 rounded-full h-2">
                <div 
                  className="bg-green-500 h-2 rounded-full transition-all duration-500"
                  style={{ 
                    width: `${(mockDevices.filter(d => d.status === 'online').length / mockDevices.length) * 100}%` 
                  }}
                />
              </div>
              <p className="text-xs text-green-300 mt-2">
                {Math.round((mockDevices.filter(d => d.status === 'online').length / mockDevices.length) * 100)}% of total devices
              </p>
            </div>

            {/* Warning Devices */}
            {mockDevices.filter(d => d.status === 'warning').length > 0 && (
              <div className="p-4 bg-yellow-500 bg-opacity-10 border border-yellow-500 border-opacity-30 rounded-lg">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-3">
                    <AlertTriangle className="w-6 h-6 text-yellow-500" />
                    <span className="font-semibold text-yellow-400">Warning</span>
                  </div>
                  <span className="text-2xl font-bold text-yellow-400">
                    {mockDevices.filter(d => d.status === 'warning').length}
                  </span>
                </div>
                <div className="w-full bg-yellow-900 bg-opacity-30 rounded-full h-2">
                  <div 
                    className="bg-yellow-500 h-2 rounded-full transition-all duration-500"
                    style={{ 
                      width: `${(mockDevices.filter(d => d.status === 'warning').length / mockDevices.length) * 100}%` 
                    }}
                  />
                </div>
                <p className="text-xs text-yellow-300 mt-2">Requires attention</p>
              </div>
            )}

            {/* Offline Devices */}
            {mockDevices.filter(d => d.status === 'offline').length > 0 && (
              <div className="p-4 bg-red-500 bg-opacity-10 border border-red-500 border-opacity-30 rounded-lg">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-3">
                    <AlertTriangle className="w-6 h-6 text-red-500" />
                    <span className="font-semibold text-red-400">Offline</span>
                  </div>
                  <span className="text-2xl font-bold text-red-400">
                    {mockDevices.filter(d => d.status === 'offline').length}
                  </span>
                </div>
                <div className="w-full bg-red-900 bg-opacity-30 rounded-full h-2">
                  <div 
                    className="bg-red-500 h-2 rounded-full transition-all duration-500"
                    style={{ 
                      width: `${(mockDevices.filter(d => d.status === 'offline').length / mockDevices.length) * 100}%` 
                    }}
                  />
                </div>
                <p className="text-xs text-red-300 mt-2">Connection lost</p>
              </div>
            )}

            {/* Calibration Devices */}
            {mockDevices.filter(d => d.status === 'calibration').length > 0 && (
              <div className="p-4 bg-blue-500 bg-opacity-10 border border-blue-500 border-opacity-30 rounded-lg">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-3">
                    <Clock className="w-6 h-6 text-blue-500" />
                    <span className="font-semibold text-blue-400">Calibration</span>
                  </div>
                  <span className="text-2xl font-bold text-blue-400">
                    {mockDevices.filter(d => d.status === 'calibration').length}
                  </span>
                </div>
                <div className="w-full bg-blue-900 bg-opacity-30 rounded-full h-2">
                  <div 
                    className="bg-blue-500 h-2 rounded-full transition-all duration-500"
                    style={{ 
                      width: `${(mockDevices.filter(d => d.status === 'calibration').length / mockDevices.length) * 100}%` 
                    }}
                  />
                </div>
                <p className="text-xs text-blue-300 mt-2">Maintenance mode</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Device Grid and Alerts */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Quick Device Status */}
        <div className="lg:col-span-2 card">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-semibold text-primary">Device Overview</h3>
            <Link to="/devices" className="text-accent hover:text-blue-400 text-sm font-medium">
              View all →
            </Link>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {mockDevices.slice(0, 4).map((device) => (
              <Link
                key={device.id}
                to={`/device/${device.id}`}
                className="p-4 bg-elevated rounded-lg border border-gray-700 hover:border-accent transition-colors group"
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    {getStatusIcon(device.status)}
                    <span className="font-medium text-primary group-hover:text-accent transition-colors">
                      {device.name}
                    </span>
                  </div>
                  <span className="text-xs text-secondary">{formatTimeAgo(device.lastSync)}</span>
                </div>
                
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-2">
                      <Thermometer className="w-4 h-4 text-orange-500" />
                      <span className="text-secondary">Temperature</span>
                    </div>
                    <span className="text-primary font-medium">
                      {device.currentReading.temperature}°C
                    </span>
                  </div>
                  
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-2">
                      <Droplets className="w-4 h-4 text-blue-500" />
                      <span className="text-secondary">Humidity</span>
                    </div>
                    <span className="text-primary font-medium">
                      {device.currentReading.humidity}%
                    </span>
                  </div>
                  
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-2">
                      <Wind className="w-4 h-4 text-green-500" />
                      <span className="text-secondary">CO2</span>
                    </div>
                    <span className="text-primary font-medium">
                      {device.currentReading.co2} ppm
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>

        {/* Recent Alerts */}
        <div className="card">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-semibold text-primary">Recent Alerts</h3>
            <button className="text-accent hover:text-blue-400 text-sm font-medium">
              View all →
            </button>
          </div>
          <div className="space-y-3">
            {mockAlerts.slice(0, 5).map((alert) => (
              <div
                key={alert.id}
                className={`p-3 rounded-lg border ${getSeverityColor(alert.severity)}`}
              >
                <div className="flex items-start justify-between gap-2">
                  <div className="flex items-center gap-2 min-w-0">
                    {alert.acknowledged ? (
                      <BellOff className="w-4 h-4 flex-shrink-0" />
                    ) : (
                      <Bell className="w-4 h-4 flex-shrink-0" />
                    )}
                    <div className="min-w-0">
                      <p className="text-sm font-medium truncate">
                        {alert.deviceName}
                      </p>
                      <p className="text-xs opacity-80 line-clamp-2">
                        {alert.message}
                      </p>
                    </div>
                  </div>
                  <span className="text-xs opacity-60 flex-shrink-0">
                    {formatTimeAgo(alert.timestamp)}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}