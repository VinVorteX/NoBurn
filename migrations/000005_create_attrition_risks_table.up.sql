CREATE TABLE IF NOT EXISTS attrition_risks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    risk_score DECIMAL(3,2) DEFAULT 0.0,
    factors JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attrition_risks_user_id ON attrition_risks(user_id);
CREATE INDEX idx_attrition_risks_risk_score ON attrition_risks(risk_score);
CREATE INDEX idx_attrition_risks_created_at ON attrition_risks(created_at);