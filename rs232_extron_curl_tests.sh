#!/bin/bash

# RS232 Extron Microservice test script
# Replace these variables with your actual values
MICROSERVICE_URL="your-microservice-url"
DEVICE_FQDN="your-device-fqdn"

echo "Testing RS232 Extron Microservice..."
echo "Microservice URL: $MICROSERVICE_URL"
echo "Device FQDN: $DEVICE_FQDN"
echo "----------------------------------------"

# GET videomute
echo "Testing GET videomute..."
curl -X GET "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/videomute/1"
sleep 1

# SET videomute
echo "Testing SET videomute..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/videomute/1" \
     -H "Content-Type: application/json" \
     -d '"true"'
sleep 1

# GET videosyncmute
echo "Testing GET videosyncmute..."
curl -X GET "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/videosyncmute/1"
sleep 1

# SET videosyncmute
echo "Testing SET videosyncmute..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/videosyncmute/1" \
     -H "Content-Type: application/json" \
     -d '"true"'
sleep 1

# SET rawcommand
echo "Testing SET rawcommand..."
curl -X PUT "http://${MICROSERVICE_URL}/${DEVICE_FQDN}/rawcommand" \
     -H "Content-Type: application/json" \
     -d '"<command>"'
sleep 1

echo "----------------------------------------"
echo "All tests completed."
