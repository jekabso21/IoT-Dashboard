package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/database"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
)

func InsertSensorData(ctx context.Context, deviceID, userID uuid.UUID, data *models.SensorDataRequest) (*models.SensorData, error) {
	sensorData := &models.SensorData{
		DeviceID:    deviceID,
		UserID:      userID,
		Temperature: data.Temperature,
		Humidity:    data.Humidity,
		CO2:         data.CO2,
	}

	err := database.DB.QueryRowContext(ctx,
		`INSERT INTO sensor_data (device_id, user_id, temperature, humidity, co2)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, time`,
		deviceID, userID, data.Temperature, data.Humidity, data.CO2,
	).Scan(&sensorData.ID, &sensorData.Time)

	if err != nil {
		return nil, fmt.Errorf("failed to insert sensor data: %w", err)
	}

	return sensorData, nil
}

func GetSensorData(ctx context.Context, query *models.SensorDataQuery) ([]models.SensorData, error) {
	sqlQuery := `SELECT id, time, device_id, user_id, temperature, humidity, co2
	             FROM sensor_data WHERE 1=1`
	args := []interface{}{}
	argPos := 1

	if query.DeviceID != nil {
		sqlQuery += fmt.Sprintf(" AND device_id = $%d", argPos)
		args = append(args, *query.DeviceID)
		argPos++
	}

	if query.StartTime != nil {
		sqlQuery += fmt.Sprintf(" AND time >= $%d", argPos)
		args = append(args, *query.StartTime)
		argPos++
	}

	if query.EndTime != nil {
		sqlQuery += fmt.Sprintf(" AND time <= $%d", argPos)
		args = append(args, *query.EndTime)
		argPos++
	}

	sqlQuery += " ORDER BY time DESC"

	if query.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, query.Limit)
		argPos++
	}

	if query.Offset > 0 {
		sqlQuery += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, query.Offset)
		argPos++
	}

	rows, err := database.DB.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor data: %w", err)
	}
	defer rows.Close()

	var sensorData []models.SensorData
	for rows.Next() {
		var data models.SensorData
		err := rows.Scan(
			&data.ID,
			&data.Time,
			&data.DeviceID,
			&data.UserID,
			&data.Temperature,
			&data.Humidity,
			&data.CO2,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %w", err)
		}
		sensorData = append(sensorData, data)
	}

	return sensorData, nil
}
