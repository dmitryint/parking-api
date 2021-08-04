# parking-api

This API is used to determine the Moscow parking space's number of the Moscow Department of Transport by geographic coordinates.

### Example

```bash
curl -X POST \
  http://localhost:5000/v1/location/ \
  -d '{"lat":"55.762375", "lon":"37.594615", "hours":1, "dist":"0.015"}' \
  -H "Content-Type: application/json"
```
