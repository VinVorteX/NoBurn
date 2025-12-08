export interface User {
  id: number;
  email: string;
  name: string;
  role: string;
  company_id: number;
}

export interface Survey {
  id: number;
  company_id: number;
  title: string;
  questions: string[];
  is_active: boolean;
  created_at: string;
}

export interface DashboardData {
  total_employees: number;
  at_risk_employees: number;
  avg_sentiment: number;
  churn_rate: number;
  top_risk_factors: string[];
  attrition_risks: AttritionRisk[];
}

export interface AttritionRisk {
  id: number;
  user_id: number;
  risk_score: number;
  user: User;
}

export interface AuthResponse {
  token: string;
  user: User;
  message?: string;
}

export interface SurveyResponse {
  survey_id: number;
  user_token: number;
  responses: string[];
}

export interface PublicSurvey {
  id: number;
  title: string;
  description?: string;
  questions: string[];
}
