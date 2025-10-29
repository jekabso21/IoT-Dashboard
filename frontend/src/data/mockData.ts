export interface Device {
  id: string;
  name: string;
  location: string;
  type: 'temperature' | 'humidity' | 'co2' | 'multi';
  status: 'online' | 'offline' | 'warning' | 'calibration';
  lastSync: string;
  battery?: number;
  currentReading: {
    temperature?: number;
    humidity?: number;
    co2?: number;
  };
  historicalData: Array<{
    timestamp: string;
    temperature?: number;
    humidity?: number;
    co2?: number;
  }>;
}

export interface Alert {
  id: string;
  deviceId: string;
  deviceName: string;
  type: 'temperature' | 'humidity' | 'co2' | 'offline' | 'battery';
  severity: 'low' | 'medium' | 'high' | 'critical';
  message: string;
  timestamp: string;
  acknowledged: boolean;
}

// Generate realistic historical data
const generateHistoricalData = (days: number = 7): Array<{
  timestamp: string;
  temperature?: number;
  humidity?: number;
  co2?: number;
}> => {
  const data = [];
  const now = new Date();
  
  for (let i = days * 24; i >= 0; i -= 2) { // Generate data every 2 hours instead of every hour
    const timestamp = new Date(now.getTime() - i * 60 * 60 * 1000);
    const hour = timestamp.getHours();
    
    // Simulate daily patterns
    const tempBase = 22 + Math.sin((hour - 6) * Math.PI / 12) * 4;
    const humidityBase = 45 + Math.sin((hour - 12) * Math.PI / 12) * 15;
    const co2Base = 400 + (hour > 8 && hour < 18 ? 200 : 0) + Math.random() * 100;
    
    data.push({
      timestamp: timestamp.toISOString(),
      temperature: Math.round((tempBase + (Math.random() - 0.5) * 2) * 10) / 10,
      humidity: Math.round((humidityBase + (Math.random() - 0.5) * 10) * 10) / 10,
      co2: Math.round(co2Base + (Math.random() - 0.5) * 50)
    });
  }
  
  return data;
};

export const mockDevices: Device[] = [
  {
    id: '1',
    name: 'Living Room Sensor',
    location: 'Living Room',
    type: 'multi',
    status: 'online',
    lastSync: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    battery: 85,
    currentReading: {
      temperature: 23.5,
      humidity: 45.2,
      co2: 420
    },
    historicalData: generateHistoricalData()
  },
  {
    id: '2',
    name: 'Server Room Monitor',
    location: 'Server Room',
    type: 'multi',
    status: 'warning',
    lastSync: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
    battery: 92,
    currentReading: {
      temperature: 28.1,
      humidity: 35.8,
      co2: 380
    },
    historicalData: generateHistoricalData()
  },
  {
    id: '3',
    name: 'Greenhouse Station',
    location: 'Greenhouse',
    type: 'multi',
    status: 'online',
    lastSync: new Date(Date.now() - 1 * 60 * 1000).toISOString(),
    battery: 78,
    currentReading: {
      temperature: 26.3,
      humidity: 65.4,
      co2: 1150
    },
    historicalData: generateHistoricalData()
  },
  {
    id: '4',
    name: 'Office Workspace',
    location: 'Office',
    type: 'multi',
    status: 'online',
    lastSync: new Date(Date.now() - 3 * 60 * 1000).toISOString(),
    battery: 67,
    currentReading: {
      temperature: 24.7,
      humidity: 42.1,
      co2: 650
    },
    historicalData: generateHistoricalData()
  },
  {
    id: '5',
    name: 'Bedroom Monitor',
    location: 'Master Bedroom',
    type: 'multi',
    status: 'offline',
    lastSync: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    battery: 15,
    currentReading: {
      temperature: 21.2,
      humidity: 48.9,
      co2: 520
    },
    historicalData: generateHistoricalData()
  },
  {
    id: '6',
    name: 'Kitchen Sensor',
    location: 'Kitchen',
    type: 'multi',
    status: 'calibration',
    lastSync: new Date(Date.now() - 15 * 60 * 1000).toISOString(),
    battery: 88,
    currentReading: {
      temperature: 25.8,
      humidity: 55.3,
      co2: 480
    },
    historicalData: generateHistoricalData()
  }
];

export const mockAlerts: Alert[] = [
  {
    id: '1',
    deviceId: '2',
    deviceName: 'Server Room Monitor',
    type: 'temperature',
    severity: 'high',
    message: 'Temperature exceeds safe operating range (28.1°C)',
    timestamp: new Date(Date.now() - 10 * 60 * 1000).toISOString(),
    acknowledged: false
  },
  {
    id: '2',
    deviceId: '5',
    deviceName: 'Bedroom Monitor',
    type: 'offline',
    severity: 'critical',
    message: 'Device has been offline for 2 hours',
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    acknowledged: false
  },
  {
    id: '3',
    deviceId: '5',
    deviceName: 'Bedroom Monitor',
    type: 'battery',
    severity: 'medium',
    message: 'Battery level is low (15%)',
    timestamp: new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString(),
    acknowledged: true
  },
  {
    id: '4',
    deviceId: '3',
    deviceName: 'Greenhouse Station',
    type: 'co2',
    severity: 'medium',
    message: 'CO2 levels are elevated (1150 ppm)',
    timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
    acknowledged: false
  },
  {
    id: '5',
    deviceId: '6',
    deviceName: 'Kitchen Sensor',
    type: 'temperature',
    severity: 'low',
    message: 'Sensor requires calibration',
    timestamp: new Date(Date.now() - 45 * 60 * 1000).toISOString(),
    acknowledged: false
  }
];

export const systemStats = {
  totalDevices: mockDevices.length,
  onlineDevices: mockDevices.filter(d => d.status === 'online').length,
  activeAlerts: mockAlerts.filter(a => !a.acknowledged).length,
  systemHealth: 85
};