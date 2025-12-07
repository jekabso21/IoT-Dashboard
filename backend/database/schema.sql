-- IoT Dashboard Database Schema
-- PostgreSQL 16+

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    company TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Pairing codes table (for user-initiated device pairing)
CREATE TABLE pairing_codes (
    code TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Devices table
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    device_name TEXT DEFAULT 'New Device',
    device_type TEXT DEFAULT 'esp32_scd4x',
    last_seen TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    paired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Device tokens table (JWT access + refresh tokens)
CREATE TABLE device_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(device_id)
);

-- Sensor data table
CREATE TABLE sensor_data (
    id BIGSERIAL PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    co2 INTEGER
);

-- Create indexes for better query performance
CREATE INDEX idx_sensor_data_device_time ON sensor_data (device_id, time DESC);
CREATE INDEX idx_sensor_data_user_time ON sensor_data (user_id, time DESC);
CREATE INDEX idx_sensor_data_time ON sensor_data (time DESC);
CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_devices_last_seen ON devices(last_seen DESC) WHERE is_active = TRUE;
CREATE INDEX idx_pairing_codes_expires ON pairing_codes(expires_at) WHERE used = FALSE;
CREATE INDEX idx_device_tokens_device ON device_tokens(device_id);
CREATE INDEX idx_device_tokens_expires ON device_tokens(expires_at);

-- Enable Row Level Security
ALTER TABLE devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE sensor_data ENABLE ROW LEVEL SECURITY;
ALTER TABLE pairing_codes ENABLE ROW LEVEL SECURITY;

-- RLS Policies for devices
CREATE POLICY devices_isolation ON devices
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true)::UUID);

-- RLS Policies for sensor_data
CREATE POLICY sensor_data_isolation ON sensor_data
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true)::UUID);

-- RLS Policies for pairing_codes
CREATE POLICY pairing_codes_isolation ON pairing_codes
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true)::UUID);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for devices table
CREATE TRIGGER update_devices_updated_at BEFORE UPDATE ON devices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to cleanup expired pairing codes
CREATE OR REPLACE FUNCTION cleanup_expired_pairing_codes()
RETURNS void AS $$
BEGIN
    DELETE FROM pairing_codes
    WHERE expires_at < NOW() AND used = FALSE;
END;
$$ LANGUAGE plpgsql;

-- Function to cleanup expired device tokens
CREATE OR REPLACE FUNCTION cleanup_expired_device_tokens()
RETURNS void AS $$
BEGIN
    DELETE FROM device_tokens
    WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Helper view: Online devices (last_seen within 30 seconds)
CREATE VIEW devices_online AS
SELECT
    d.*,
    CASE
        WHEN d.last_seen > NOW() - INTERVAL '30 seconds' THEN TRUE
        ELSE FALSE
    END as online
FROM devices d
WHERE d.is_active = TRUE;

-- Helper function: Generate pairing code
CREATE OR REPLACE FUNCTION generate_pairing_code()
RETURNS TEXT AS $$
DECLARE
    new_code TEXT;
    code_exists BOOLEAN;
BEGIN
    LOOP
        -- Generate random 6-digit code
        new_code := LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');

        -- Check if code already exists and is not expired
        SELECT EXISTS(
            SELECT 1 FROM pairing_codes
            WHERE code = new_code
            AND expires_at > NOW()
        ) INTO code_exists;

        -- Exit loop if code is unique
        EXIT WHEN NOT code_exists;
    END LOOP;

    RETURN new_code;
END;
$$ LANGUAGE plpgsql;

-- Helper function: Check if device should be considered online
CREATE OR REPLACE FUNCTION is_device_online(device_uuid UUID)
RETURNS BOOLEAN AS $$
DECLARE
    device_last_seen TIMESTAMPTZ;
BEGIN
    SELECT last_seen INTO device_last_seen
    FROM devices
    WHERE id = device_uuid;

    IF device_last_seen IS NULL THEN
        RETURN FALSE;
    END IF;

    RETURN device_last_seen > NOW() - INTERVAL '30 seconds';
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE users IS 'User accounts for the IoT dashboard';
COMMENT ON TABLE pairing_codes IS 'Temporary codes for pairing new devices (user-initiated)';
COMMENT ON TABLE devices IS 'ESP32 devices paired to user accounts';
COMMENT ON TABLE device_tokens IS 'JWT access and refresh tokens for device authentication';
COMMENT ON TABLE sensor_data IS 'Time-series sensor data from devices';
COMMENT ON VIEW devices_online IS 'View of active devices with online status (last_seen < 30s)';

COMMENT ON COLUMN devices.last_seen IS 'Updated automatically on every data transmission';
COMMENT ON COLUMN pairing_codes.code IS '6-digit numeric code generated by user on dashboard';
COMMENT ON COLUMN pairing_codes.expires_at IS 'Pairing codes expire after 10 minutes';
COMMENT ON COLUMN device_tokens.expires_at IS 'Access token expiry (refresh token separately managed)';
