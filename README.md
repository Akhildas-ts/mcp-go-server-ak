# MCP Vector Search Server

A Go-based server for vector search functionality using Pinecone and OpenAI.

## Setup

### 1. Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=8081
JWT_SECRET=mcp-secret-key

# Pinecone Configuration (Required)
PINECONE_API_KEY=your-pinecone-api-key
PINECONE_INDEX_NAME=your-pinecone-index-name
PINECONE_HOST=your-pinecone-host

# OpenAI Configuration (Required)
OPENAI_API_KEY=your-openai-api-key-here

# GitHub OAuth Configuration (Optional for development)
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_OAUTH_REDIRECT_URL=http://localhost:8081/auth/github/callback
```

### 2. Required API Keys

You need to obtain the following API keys:

1. **Pinecone API Key**: 
   - Sign up at [Pinecone](https://www.pinecone.io/)
   - Create an index and get your API key
   - Set `PINECONE_API_KEY`, `PINECONE_INDEX_NAME`, and `PINECONE_HOST`

2. **OpenAI API Key**:
   - Sign up at [OpenAI](https://platform.openai.com/)
   - Get your API key from the dashboard
   - Set `OPENAI_API_KEY`

### 3. Running the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go
```

The server will start on port 8081 (or the port specified in your .env file).

## API Endpoints

- Health check: `GET /health`
- Search: `POST /search`
- Index: `POST /index`
- Authentication endpoints: `/auth/*`

## Development

For development, you can set dummy values for optional fields, but `PINECONE_API_KEY` and `OPENAI_API_KEY` are required for the application to start. # mcp-go-server-ak
