DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'gateways') THEN
        CREATE TABLE gateways (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            data_format_supported VARCHAR(50) NOT NULL,
            health_status VARCHAR(20) DEFAULT 'healthy', -- Track gateway health: 'healthy', 'unhealthy'
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
            amount DECIMAL(10, 2) NOT NULL,
            type VARCHAR(50) NOT NULL, -- deposit/withdrawal
            status VARCHAR(50) NOT NULL, -- pending, completed, failed
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
            gateway_id INT NOT NULL,  
            country_id INT NOT NULL,  
            user_id INT NOT NULL,
            FOREIGN KEY (gateway_id) REFERENCES gateways(id) ON DELETE SET NULL,
            FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        );
    END IF;
END $$;

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL UNIQUE,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            balance DECIMAL(12, 2) NOT NULL DEFAULT 0.00, -- User balance column
            country_id INT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;

-- Add indexes to optimize queries for priority, health status, and status updates
CREATE INDEX IF NOT EXISTS idx_gateway_health_status ON gateways(health_status);
CREATE INDEX IF NOT EXISTS idx_gateway_last_checked ON gateways(last_checked_at);
CREATE INDEX IF NOT EXISTS idx_transaction_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transaction_gateway_id ON transactions(gateway_id);