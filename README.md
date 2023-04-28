# Everest Backend

## Running

Run services
```
docker compose up -d
```

Run backend
```
go run .
```

## API
### Token
Login at http://localhost:3000/login  
Copy token from the response.

### Calling API
```sh
TOKEN=your_token_here
curl localhost:3000/api/v1/project/create \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "my project"}'
```

### Endpoints

- POST `/api/v1/project/create`
  ```json
  {
    "name": "string"
  }
  ```
- POST `/api/v1/user/add-role`
  ```json
  {
    "project_id": 0,
    "role": "string",
    "user_id": 0
  }
  ```
- POST `/api/v1/user/remove-role`
  ```json
  {
    "project_id": 0,
    "role": "string",
    "user_id": 0
  }
  ```
- POST `/api/v1/cluster/create/:projectID`
- POST `/api/v1/cluster/delete/:projectID`

### Roles

- `cluster.create`
- `cluster.delete`
