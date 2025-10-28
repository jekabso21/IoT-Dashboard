import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { 
  Plus, 
  Search, 
  Filter, 
  Grid3X3, 
  List, 
  Thermometer, 
  Droplets, 
  Wind,
  CheckCircle,
  AlertTriangle,
  Clock,
  Battery,
  QrCode,
  X,
  Wifi,
  WifiOff
} from 'lucide-react';
import { mockDevices } from '../data/mockData';

export default function DeviceManagement() {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [locationFilter, setLocationFilter] = useState('all');
  const [showAddModal, setShowAddModal] = useState(false);
  const [addMode, setAddMode] = useState<'manual' | 'qr'>('manual');
  const [newDeviceId, setNewDeviceId] = useState('');

  const filteredDevices = mockDevices.filter(device => {
    const matchesSearch = device.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         device.location.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || device.status === statusFilter;
    const matchesLocation = locationFilter === 'all' || device.location === locationFilter;
    
    return matchesSearch && matchesStatus && matchesLocation;
  });

  const locations = [...new Set(mockDevices.map(device => device.location))];

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online': return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'warning': return <AlertTriangle className="w-4 h-4 text-yellow-500" />;
      case 'offline': return <AlertTriangle className="w-4 h-4 text-red-500" />;
      case 'calibration': return <Clock className="w-4 h-4 text-blue-500" />;
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

  const formatTimeAgo = (timestamp: string) => {
    const now = new Date();
    const time = new Date(timestamp);
    const diffInMinutes = Math.floor((now.getTime() - time.getTime()) / (1000 * 60));
    
    if (diffInMinutes < 1) return 'Just now';
    if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}h ago`;
    return `${Math.floor(diffInMinutes / 1440)}d ago`;
  };

  const getBatteryColor = (level: number) => {
    if (level > 50) return 'text-green-500';
    if (level > 20) return 'text-yellow-500';
    return 'text-red-500';
  };

  const handleAddDevice = () => {
    // Simulate adding device
    console.log('Adding device:', newDeviceId);
    setShowAddModal(false);
    setNewDeviceId('');
    setAddMode('manual');
  };

  return (
    <div className="space-y-6 fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-primary">Device Management</h1>
          <p className="text-secondary mt-1">
            Monitor and manage your IoT devices
          </p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="btn btn-primary"
        >
          <Plus className="w-4 h-4" />
          Add Device
        </button>
      </div>

      {/* Filters and Search */}
      <div className="card">
        <div className="flex flex-col lg:flex-row gap-4">
          {/* Search */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-secondary" />
              <input
                type="text"
                placeholder="Search devices..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="form-input pl-10"
              />
            </div>
          </div>

          {/* Filters */}
          <div className="flex gap-4">
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="form-input min-w-0"
            >
              <option value="all">All Status</option>
              <option value="online">Online</option>
              <option value="warning">Warning</option>
              <option value="offline">Offline</option>
              <option value="calibration">Calibration</option>
            </select>

            <select
              value={locationFilter}
              onChange={(e) => setLocationFilter(e.target.value)}
              className="form-input min-w-0"
            >
              <option value="all">All Locations</option>
              {locations.map(location => (
                <option key={location} value={location}>{location}</option>
              ))}
            </select>

            {/* View Mode Toggle */}
            <div className="flex bg-elevated rounded-lg p-1">
              <button
                onClick={() => setViewMode('grid')}
                className={`p-2 rounded ${viewMode === 'grid' ? 'bg-accent text-primary' : 'text-secondary hover:text-primary'}`}
              >
                <Grid3X3 className="w-4 h-4" />
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`p-2 rounded ${viewMode === 'list' ? 'bg-accent text-primary' : 'text-secondary hover:text-primary'}`}
              >
                <List className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Results Summary */}
      <div className="flex items-center justify-between">
        <p className="text-secondary">
          Showing {filteredDevices.length} of {mockDevices.length} devices
        </p>
        <div className="flex items-center gap-4 text-sm">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-green-500"></div>
            <span className="text-secondary">
              {mockDevices.filter(d => d.status === 'online').length} Online
            </span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-red-500"></div>
            <span className="text-secondary">
              {mockDevices.filter(d => d.status === 'offline').length} Offline
            </span>
          </div>
        </div>
      </div>

      {/* Device Grid/List */}
      {viewMode === 'grid' ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredDevices.map((device) => (
            <Link
              key={device.id}
              to={`/device/${device.id}`}
              className="card hover:shadow-lg hover:border-accent transition-all group"
            >
              {/* Device Header */}
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  {getStatusIcon(device.status)}
                  <div>
                    <h3 className="font-semibold text-primary group-hover:text-accent transition-colors">
                      {device.name}
                    </h3>
                    <p className="text-sm text-secondary">{device.location}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {device.status === 'online' ? (
                    <Wifi className="w-4 h-4 text-green-500" />
                  ) : (
                    <WifiOff className="w-4 h-4 text-red-500" />
                  )}
                  {device.battery && (
                    <div className="flex items-center gap-1">
                      <Battery className={`w-4 h-4 ${getBatteryColor(device.battery)}`} />
                      <span className={`text-xs ${getBatteryColor(device.battery)}`}>
                        {device.battery}%
                      </span>
                    </div>
                  )}
                </div>
              </div>

              {/* Current Readings */}
              <div className="space-y-3 mb-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Thermometer className="w-4 h-4 text-orange-500" />
                    <span className="text-sm text-secondary">Temperature</span>
                  </div>
                  <span className="text-sm font-medium text-primary">
                    {device.currentReading.temperature}°C
                  </span>
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Droplets className="w-4 h-4 text-blue-500" />
                    <span className="text-sm text-secondary">Humidity</span>
                  </div>
                  <span className="text-sm font-medium text-primary">
                    {device.currentReading.humidity}%
                  </span>
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Wind className="w-4 h-4 text-green-500" />
                    <span className="text-sm text-secondary">CO2</span>
                  </div>
                  <span className="text-sm font-medium text-primary">
                    {device.currentReading.co2} ppm
                  </span>
                </div>
              </div>

              {/* Footer */}
              <div className="flex items-center justify-between pt-3 border-t border-gray-700">
                <span className={`text-sm capitalize ${getStatusColor(device.status)}`}>
                  {device.status}
                </span>
                <span className="text-xs text-secondary">
                  {formatTimeAgo(device.lastSync)}
                </span>
              </div>
            </Link>
          ))}
        </div>
      ) : (
        <div className="card">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-700">
                  <th className="text-left py-3 px-4 font-medium text-secondary">Device</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary">Status</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary">Temperature</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary">Humidity</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary">CO2</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary">Battery</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary">Last Sync</th>
                </tr>
              </thead>
              <tbody>
                {filteredDevices.map((device) => (
                  <tr
                    key={device.id}
                    className="border-b border-gray-800 hover:bg-elevated transition-colors"
                  >
                    <td className="py-4 px-4">
                      <Link
                        to={`/device/${device.id}`}
                        className="flex items-center gap-3 hover:text-accent transition-colors"
                      >
                        {device.status === 'online' ? (
                          <Wifi className="w-4 h-4 text-green-500" />
                        ) : (
                          <WifiOff className="w-4 h-4 text-red-500" />
                        )}
                        <div>
                          <p className="font-medium text-primary">{device.name}</p>
                          <p className="text-sm text-secondary">{device.location}</p>
                        </div>
                      </Link>
                    </td>
                    <td className="py-4 px-4">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(device.status)}
                        <span className={`text-sm capitalize ${getStatusColor(device.status)}`}>
                          {device.status}
                        </span>
                      </div>
                    </td>
                    <td className="py-4 px-4 text-primary">
                      {device.currentReading.temperature}°C
                    </td>
                    <td className="py-4 px-4 text-primary">
                      {device.currentReading.humidity}%
                    </td>
                    <td className="py-4 px-4 text-primary">
                      {device.currentReading.co2} ppm
                    </td>
                    <td className="py-4 px-4">
                      {device.battery ? (
                        <div className="flex items-center gap-1">
                          <Battery className={`w-4 h-4 ${getBatteryColor(device.battery)}`} />
                          <span className={`text-sm ${getBatteryColor(device.battery)}`}>
                            {device.battery}%
                          </span>
                        </div>
                      ) : (
                        <span className="text-secondary">N/A</span>
                      )}
                    </td>
                    <td className="py-4 px-4 text-secondary text-sm">
                      {formatTimeAgo(device.lastSync)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Add Device Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-secondary rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-primary">Add New Device</h2>
              <button
                onClick={() => setShowAddModal(false)}
                className="text-secondary hover:text-primary"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            {/* Mode Selection */}
            <div className="flex bg-elevated rounded-lg p-1 mb-6">
              <button
                onClick={() => setAddMode('manual')}
                className={`flex-1 py-2 px-4 rounded text-sm font-medium transition-colors ${
                  addMode === 'manual' ? 'bg-accent text-primary' : 'text-secondary hover:text-primary'
                }`}
              >
                Manual Entry
              </button>
              <button
                onClick={() => setAddMode('qr')}
                className={`flex-1 py-2 px-4 rounded text-sm font-medium transition-colors ${
                  addMode === 'qr' ? 'bg-accent text-primary' : 'text-secondary hover:text-primary'
                }`}
              >
                QR Code
              </button>
            </div>

            {addMode === 'manual' ? (
              <div className="space-y-4">
                <div className="form-group">
                  <label className="form-label">Device ID</label>
                  <input
                    type="text"
                    value={newDeviceId}
                    onChange={(e) => setNewDeviceId(e.target.value)}
                    placeholder="Enter device ID (e.g., ECO-001)"
                    className="form-input"
                  />
                </div>
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowAddModal(false)}
                    className="btn btn-secondary flex-1"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleAddDevice}
                    disabled={!newDeviceId.trim()}
                    className="btn btn-primary flex-1"
                  >
                    Add Device
                  </button>
                </div>
              </div>
            ) : (
              <div className="text-center space-y-4">
                <div className="w-32 h-32 bg-elevated rounded-lg flex items-center justify-center mx-auto">
                  <QrCode className="w-16 h-16 text-secondary" />
                </div>
                <p className="text-secondary">
                  Position the QR code within the frame to scan
                </p>
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowAddModal(false)}
                    className="btn btn-secondary flex-1"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleAddDevice}
                    className="btn btn-primary flex-1"
                  >
                    Simulate Scan
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}