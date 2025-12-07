CREATE TABLE IF NOT EXISTS surveys (
    id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES companies(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    questions JSONB NOT NULL DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_surveys_company_id ON surveys(company_id);
CREATE INDEX idx_surveys_is_active ON surveys(is_active);
CREATE INDEX idx_surveys_deleted_at ON surveys(deleted_at);