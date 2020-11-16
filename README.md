# emosearch-api

## Commands

### Development

```bash
# Launch DynamoDB local (:8000) and DynamoDB Admin (:8001)
$ docker-compose up -d
# Create DynamoDB table
$ aws dynamodb create-table --cli-input-json file://config/dynamodb.json --endpoint-url http://localhost:8000

# Create environments file, and you have to edit some secrets.
$ cp config/sam-dev-env.example.json config/sam-dev-env.json

# Start API (:9000)
$ make
$ sam local start-api --port 9000 --env-vars config/sam-dev-env.json --docker-network emosearch-api_default

# Invoke a function manually
$ sam local invoke "ListSearchesToUpdateFunction" --env-vars config/sam-dev-env.json --docker-network emosearch-api_default

# Invoke a function with Lambda Event
$ touch event.json && echo '{"search_id":"123","user_id":"123"}' >> event.json
$ sam local invoke "ListSearchesToUpdateFunction" --event event.json --env-vars config/sam-dev-env.json --docker-network emosearch-api_default
```

### Deployment

```bash
# Deploy by AWS CloudFormation, and append tags for cost management.
sam deploy --tags "Project=EmoSearchAPI"
```

Next, open AWS Secrets Manager console, and edit secrets to "GoogleServiceAccountKey", "TwitterConsumerSecret" and "TwitterConsumerKey".

After editing, you can see the API endpoint from CloudFormation output resoures.
