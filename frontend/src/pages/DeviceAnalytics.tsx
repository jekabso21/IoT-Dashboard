import React, { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { 
  ArrowLeft, 
  Download, 
  Calendar, 
  Thermometer, 
  Droplets, 
  Wind,
  Settings,
  AlertTriangle,
  CheckCircle,
  Clock,
  Battery,
  Wifi,
  WifiOff
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';
import { mockDevices } from '../data/mockData';

export default function DeviceAnalytics() {
  const { id } = useParams();
  const [timeRange, setTimeRange] = useState('24H');
  const [selectedMetric, setSelectedMetric] = useState('all');
  const [showThresholds, setShowThresholds] = useState(true);

  const device = mockDevices.find(d => d.id === id);

  if (!device) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <p className="text-xl text-secondary mb-4">Device not found</p>
          <Link to="/devices" className="btn btn-primary">
            <ArrowLeft className="w-4 h-4" />
            Back to Devices
          </Link>
        </div>
      </div>
    );
  }

  const getTimeRangeData = () => {
    const hours = timeRange === '1H' ? 1 : timeRange === '24H' ? 24 : timeRange === '7D' ? 168 : 720;
    return device.historicalData.slice(-hours);
  };

  const chartData = getTimeRangeData().map(item => ({
    timestamp: new Date(item.timestamp).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: timeRange === '1H' ? '2-digit' : undefined,
      minute: timeRange === '1H' ? '2-digit' : undefined
    }),
    temperature: item.temperature,
    humidity: item.humidity,
    co2: item.co2
  }));

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online': return <CheckCircle className="w-5 h-5 text-green-500" />;
      case 'warning': return <AlertTriangle className="w-5 h-5 text-yellow-500" />;
      case 'offline': return <AlertTriangle className="w-5 h-5 text-red-500" />;
      case 'calibration': return <Clock className="w-5 h-5 text-blue-500" />;
      default: return null;
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

  const getBatteryColor = (level: number) => {
    if (level > 50) return 'text-green-500';
    if (level > 20) return 'text-yellow-500';
    return 'text-red-500';
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

  const handleExport = (format: 'csv' | 'pdf') => {
    // Simulate export
    console.log(`Exporting data as ${format.toUpperCase()}`);
  };

  const timeRanges = [
    { value: '1H', label: '1 Hour' },
    { value: '24H', label: '24 Hours' },
    { value: '7D', label: '7 Days' },
    { value: '30D', label: '30 Days' }
  ];

  const metrics = [
    { value: 'all', label: 'All Metrics' },
    { value: 'temperature', label: 'Temperature' },
    { value: 'humidity', label: 'Humidity' },
    { value: 'co2', label: 'CO2' }
  ];

  return (
    <div className="space-y-6 fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div className="flex items-center gap-4">
          <Link
            to="/devices"
            className="p-2 hover:bg-elevated rounded-lg transition-colors"
          >
            <ArrowLeft className="w-5 h-5 text-secondary hover:text-primary" />
          </Link>
          <div>
            <h1 className="text-3xl font-bold text-primary">{device.name}</h1>
            <p className="text-secondary mt-1">{device.location}</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => handleExport('csv')}
            className="btn btn-secondary"
          >
            <Download className="w-4 h-4" />
            Export CSV
          </button>
          <button
            onClick={() => handleExport('pdf')}
            className="btn btn-secondary"
          >
            <Download className="w-4 h-4" />
            Export PDF
          </button>
          <button className="btn btn-primary">
            <Settings className="w-4 h-4" />
            Configure
          </button>
        </div>
      </div>

      {/* Device Status Card */}
      <div className="card">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="flex items-center gap-3">
            {getStatusIcon(device.status)}
            <div>
              <p className="text-sm text-secondary">Status</p>
              <p className={`font-medium capitalize ${getStatusColor(device.status)}`}>
                {device.status}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            {device.status === 'online' ? (
              <Wifi className="w-5 h-5 text-green-500" />
            ) : (
              <WifiOff className="w-5 h-5 text-red-500" />
            )}
            <div>
              <p className="text-sm text-secondary">Connection</p>
              <p className="font-medium text-primary">
                {device.status === 'online' ? 'Connected' : 'Disconnected'}
              </p>
            </div>
          </div>

          {device.battery && (
            <div className="flex items-center gap-3">
              <Battery className={`w-5 h-5 ${getBatteryColor(device.battery)}`} />
              <div>
                <p className="text-sm text-secondary">Battery</p>
                <p className={`font-medium ${getBatteryColor(device.battery)}`}>
                  {device.battery}%
                </p>
              </div>
            </div>
          )}

          <div className="flex items-center gap-3">
            <Clock className="w-5 h-5 text-blue-500" />
            <div>
              <p className="text-sm text-secondary">Last Sync</p>
              <p className="font-medium text-primary">
                {formatTimeAgo(device.lastSync)}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Current Readings */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 gradient-temperature rounded-lg flex items-center justify-center">
              <Thermometer className="w-6 h-6 text-white" />
            </div>
            <div>
              <p className="text-sm text-secondary">Temperature</p>
              <p className="text-2xl font-bold text-primary">
                {device.currentReading.temperature}°C
              </p>
            </div>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-secondary">Range: 18-26°C</span>
            <span className="text-green-500">Normal</span>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 gradient-humidity rounded-lg flex items-center justify-center">
              <Droplets className="w-6 h-6 text-white" />
            </div>
            <div>
              <p className="text-sm text-secondary">Humidity</p>
              <p className="text-2xl font-bold text-primary">
                {device.currentReading.humidity}%
              </p>
            </div>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-secondary">Range: 30-60%</span>
            <span className="text-green-500">Normal</span>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 gradient-co2 rounded-lg flex items-center justify-center">
              <Wind className="w-6 h-6 text-white" />
            </div>
            <div>
              <p className="text-sm text-secondary">CO2 Level</p>
              <p className="text-2xl font-bold text-primary">
                {device.currentReading.co2} ppm
              </p>
            </div>
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="text-secondary">Range: 400-1000 ppm</span>
            <span className={device.currentReading.co2! > 1000 ? 'text-yellow-500' : 'text-green-500'}>
              {device.currentReading.co2! > 1000 ? 'Elevated' : 'Normal'}
            </span>
          </div>
        </div>
      </div>

      {/* Chart Controls */}
      <div className="card">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Calendar className="w-4 h-4 text-secondary" />
              <span className="text-sm text-secondary">Time Range:</span>
            </div>
            <div className="flex bg-elevated rounded-lg p-1">
              {timeRanges.map((range) => (
                <button
                  key={range.value}
                  onClick={() => setTimeRange(range.value)}
                  className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
                    timeRange === range.value
                      ? 'bg-accent text-primary'
                      : 'text-secondary hover:text-primary'
                  }`}
                >
                  {range.label}
                </button>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-4">
            <select
              value={selectedMetric}
              onChange={(e) => setSelectedMetric(e.target.value)}
              className="form-input min-w-0"
            >
              {metrics.map((metric) => (
                <option key={metric.value} value={metric.value}>
                  {metric.label}
                </option>
              ))}
            </select>

            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={showThresholds}
                onChange={(e) => setShowThresholds(e.target.checked)}
                className="w-4 h-4 text-accent bg-elevated border-gray-600 rounded focus:ring-accent focus:ring-2"
              />
              <span className="text-sm text-secondary">Show Thresholds</span>
            </label>
          </div>
        </div>
      </div>

      {/* Main Chart */}
      <div className="card">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-xl font-semibold text-primary">Environmental Data</h3>
          <div className="flex items-center gap-4 text-sm">
            {(selectedMetric === 'all' || selectedMetric === 'temperature') && (
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-orange-500"></div>
                <span className="text-secondary">Temperature</span>
              </div>
            )}
            {(selectedMetric === 'all' || selectedMetric === 'humidity') && (
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-blue-500"></div>
                <span className="text-secondary">Humidity</span>
              </div>
            )}
            {(selectedMetric === 'all' || selectedMetric === 'co2') && (
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                <span className="text-secondary">CO2</span>
              </div>
            )}
          </div>
        </div>

        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            {selectedMetric === 'all' ? (
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#404040" />
                <XAxis dataKey="timestamp" stroke="#b3b3b3" />
                <YAxis stroke="#b3b3b3" />
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
            ) : (
              <AreaChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#404040" />
                <XAxis dataKey="timestamp" stroke="#b3b3b3" />
                <YAxis stroke="#b3b3b3" />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#2d2d2d', 
                    border: '1px solid #404040',
                    borderRadius: '8px',
                    color: '#ffffff'
                  }}
                />
                <Area
                  type="monotone"
                  dataKey={selectedMetric}
                  stroke={
                    selectedMetric === 'temperature' ? '#ff6b35' :
                    selectedMetric === 'humidity' ? '#2196f3' : '#4caf50'
                  }
                  fill={
                    selectedMetric === 'temperature' ? '#ff6b35' :
                    selectedMetric === 'humidity' ? '#2196f3' : '#4caf50'
                  }
                  fillOpacity={0.2}
                  strokeWidth={2}
                />
              </AreaChart>
            )}
          </ResponsiveContainer>
        </div>
      </div>

      {/* Alert Thresholds */}
      <div className="card">
        <h3 className="text-xl font-semibold text-primary mb-6">Alert Thresholds</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <Thermometer className="w-5 h-5 text-orange-500" />
              <span className="font-medium text-primary">Temperature</span>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm text-secondary">Min Threshold</span>
                <input
                  type="number"
                  defaultValue="18"
                  className="w-20 px-2 py-1 bg-elevated border border-gray-600 rounded text-sm"
                />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-secondary">Max Threshold</span>
                <input
                  type="number"
                  defaultValue="26"
                  className="w-20 px-2 py-1 bg-elevated border border-gray-600 rounded text-sm"
                />
              </div>
            </div>
          </div>

          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <Droplets className="w-5 h-5 text-blue-500" />
              <span className="font-medium text-primary">Humidity</span>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm text-secondary">Min Threshold</span>
                <input
                  type="number"
                  defaultValue="30"
                  className="w-20 px-2 py-1 bg-elevated border border-gray-600 rounded text-sm"
                />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-secondary">Max Threshold</span>
                <input
                  type="number"
                  defaultValue="60"
                  className="w-20 px-2 py-1 bg-elevated border border-gray-600 rounded text-sm"
                />
              </div>
            </div>
          </div>

          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <Wind className="w-5 h-5 text-green-500" />
              <span className="font-medium text-primary">CO2</span>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm text-secondary">Warning Level</span>
                <input
                  type="number"
                  defaultValue="1000"
                  className="w-20 px-2 py-1 bg-elevated border border-gray-600 rounded text-sm"
                />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-secondary">Critical Level</span>
                <input
                  type="number"
                  defaultValue="1500"
                  className="w-20 px-2 py-1 bg-elevated border border-gray-600 rounded text-sm"
                />
              </div>
            </div>
          </div>
        </div>
        <div className="flex justify-end mt-6">
          <button className="btn btn-primary">
            Save Thresholds
          </button>
        </div>
      </div>
    </div>
  );
}