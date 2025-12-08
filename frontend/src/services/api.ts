import axios from 'axios';
import type { User, Survey, DashboardData, AttritionRisk, AuthResponse, SurveyResponse, PublicSurvey } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', { email, password });
    return response.data.data;
  },
  
  register: async (data: {
    company_name: string;
    name: string;
    email: string;
    password: string;
  }): Promise<AuthResponse> => {
    const response = await api.post('/auth/register', data);
    const registerData = response.data.data;
    // Register returns only token and message, need to login to get user
    if (registerData.token && !registerData.user) {
      const loginResponse = await api.post('/auth/login', { 
        email: data.email, 
        password: data.password 
      });
      return loginResponse.data.data;
    }
    return registerData;
  },
};

export const surveyAPI = {
  getSurveys: async (): Promise<Survey[]> => {
    const response = await api.get('/api/surveys');
    return response.data.data || [];
  },
  
  createSurvey: async (data: { title: string; questions: string[] }): Promise<Survey> => {
    const response = await api.post('/api/surveys', data);
    return response.data.data;
  },
  
  getPublicSurvey: async (surveyId: number): Promise<PublicSurvey> => {
    const response = await api.get(`/api/surveys/${surveyId}/public`);
    return response.data.data;
  },
  
  submitPublicResponse: async (data: SurveyResponse): Promise<{ message: string }> => {
    const response = await api.post('/api/surveys/responses/public', data);
    return response.data.data;
  },
};

export const analyticsAPI = {
  getDashboard: async (): Promise<DashboardData> => {
    const response = await api.get('/api/dashboard');
    return response.data.data;
  },
  
  getAttritionRisks: async (): Promise<AttritionRisk[]> => {
    const response = await api.get('/api/attrition-risks');
    return response.data.data || [];
  },
};

export const employeeAPI = {
  getEmployees: async (): Promise<User[]> => {
    const response = await api.get('/api/employees');
    return response.data.data || [];
  },
  
  addEmployee: async (data: { email: string; name: string }): Promise<User> => {
    const response = await api.post('/api/employees', data);
    return response.data.data;
  },
  
  deleteEmployee: async (employeeId: number): Promise<{ message: string }> => {
    const response = await api.delete(`/api/employees/${employeeId}`);
    return response.data.data;
  },
  
  bulkUpload: async (file: File): Promise<{ message: string }> => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await api.post('/api/employees/bulk', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data.data;
  },
};

export const settingsAPI = {
  getSMTP: async (): Promise<any> => {
    const response = await api.get('/api/settings/smtp');
    return response.data.data;
  },
  
  updateSMTP: async (data: any): Promise<{ message: string }> => {
    const response = await api.put('/api/settings/smtp', data);
    return response.data.data;
  },
};

export default api;
