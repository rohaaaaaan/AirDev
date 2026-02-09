-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Projects Table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    repo_url VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'Idle', -- Display status
    state VARCHAR(20) DEFAULT 'idle',   -- success, error, idle, running
    last_build_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Agents Table
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    machine_fingerprint VARCHAR(255),
    status VARCHAR(20) DEFAULT 'OFFLINE', -- ONLINE, OFFLINE, BUSY
    last_seen_at TIMESTAMP WITH TIME ZONE,
    version VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Jobs Table
CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- BUILD, TEST, DEPLOY
    status VARCHAR(50) DEFAULT 'CREATED', -- CREATED, QUEUED, RUNNING, COMPLETED, FAILED, CANCELLED
    result TEXT, -- JSON result
    input_params TEXT, -- JSON params
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
