DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL UNIQUE,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            country_id INT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gateways') THEN
        CREATE TYPE health_status_enum AS ENUM ('healthy', 'unhealthy');

        CREATE TABLE gateways (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            data_format_supported VARCHAR(50) NOT NULL,
            health_status health_status_enum DEFAULT 'healthy', -- Track gateway health: ENUM type
            last_checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of the last health check
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'countries') THEN
        CREATE TABLE countries (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            code CHAR(2) NOT NULL UNIQUE,
            currency CHAR(3) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gateway_countries') THEN
        CREATE TABLE gateway_countries (
            gateway_id INT NOT NULL,
            country_id INT NOT NULL,
            priority INT DEFAULT 1, -- Priority for gateways in specific countries
            PRIMARY KEY (gateway_id, country_id),
            FOREIGN KEY (gateway_id) REFERENCES gateways(id) ON DELETE CASCADE,
            FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
        );
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'transactions') THEN
        CREATE TABLE transactions (
            id SERIAL PRIMARY KEY,
            reference_id UUID NOT NULL UNIQUE, -- unique identifier for the transaction
            amount DECIMAL(10, 2) NOT NULL,
            currency CHAR(3) NOT NULL,
            type VARCHAR(50) NOT NULL, -- deposit/withdrawal
            status VARCHAR(50) NOT NULL, -- pending, retry, completed, failed
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
            gateway_id INT,
            country_id INT NOT NULL,
            user_id INT NOT NULL,
            FOREIGN KEY (gateway_id) REFERENCES gateways(id) ON DELETE SET NULL,
            FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        );
    END IF;
END $$;

-- Add indexes to optimize queries for priority, health status, and status updates
CREATE INDEX IF NOT EXISTS idx_gateway_health_status ON gateways(health_status);
CREATE INDEX IF NOT EXISTS idx_gateway_last_checked ON gateways(last_checked_at);
CREATE INDEX IF NOT EXISTS idx_transaction_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transaction_gateway_id ON transactions(gateway_id);
CREATE INDEX IF NOT EXISTS idx_countries_code ON countries(code);
CREATE INDEX IF NOT EXISTS idx_gateway_countries_country_id ON gateway_countries(country_id);
CREATE INDEX IF NOT EXISTS idx_gateway_countries_composite ON gateway_countries(country_id, priority, gateway_id);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_id ON transactions(reference_id);

-- Populate countries, gateways, and a user
INSERT INTO countries (name, code, currency, created_at, updated_at)
VALUES 
    ('United States', 'US', 'USD', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Indonesia', 'ID', 'IDR', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

INSERT INTO gateways (name, data_format_supported, health_status, created_at, updated_at)
VALUES 
    ('A', 'json', 'healthy', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('B', 'soap', 'healthy', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('C', 'json', 'healthy', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- Populate countries, gateways, and a user
INSERT INTO gateway_countries (gateway_id, country_id, priority)
VALUES
    (1, 1, 1),
    (2, 1, 2),
    (3, 2, 1)
ON CONFLICT DO NOTHING;

INSERT INTO users (username, email, password, country_id, created_at, updated_at)
VALUES 
    ('test_user', 'test_user@example.com', 'hashed_password', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
