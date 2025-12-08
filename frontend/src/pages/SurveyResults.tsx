import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { FileText, TrendingUp, TrendingDown, Minus, Calendar, User } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import axios from 'axios';
import { toast } from 'sonner';

interface SurveyResponse {
  id: number;
  user_id: number;
  responses: string[];
  sentiment: number;
  created_at: string;
  user: {
    name: string;
    email: string;
  };
}

interface Survey {
  id: number;
  title: string;
  questions: string[];
  created_at: string;
}

export const SurveyResults = () => {
  const { surveyId } = useParams<{ surveyId: string }>();
  const [loading, setLoading] = useState(true);
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [responses, setResponses] = useState<SurveyResponse[]>([]);

  useEffect(() => {
    fetchSurveyResults();
  }, [surveyId]);

  const fetchSurveyResults = async () => {
    try {
      const token = localStorage.getItem('token');
      const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';
      const response = await axios.get(
        `${API_BASE_URL}/api/surveys/${surveyId}/responses`,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      
      setSurvey(response.data.data.survey);
      setResponses(response.data.data.responses || []);
    } catch (error: any) {
      toast.error('Failed to load survey results');
    } finally {
      setLoading(false);
    }
  };

  const getSentimentIcon = (sentiment: number) => {
    if (sentiment > 0.3) return <TrendingUp className="h-5 w-5 text-green-500" />;
    if (sentiment < -0.3) return <TrendingDown className="h-5 w-5 text-red-500" />;
    return <Minus className="h-5 w-5 text-yellow-500" />;
  };

  const getSentimentLabel = (sentiment: number) => {
    if (sentiment > 0.3) return 'Positive';
    if (sentiment < -0.3) return 'Negative';
    return 'Neutral';
  };

  const getSentimentColor = (sentiment: number) => {
    if (sentiment > 0.3) return 'text-green-600 bg-green-50';
    if (sentiment < -0.3) return 'text-red-600 bg-red-50';
    return 'text-yellow-600 bg-yellow-50';
  };

  const avgSentiment = responses.length > 0
    ? responses.reduce((sum, r) => sum + r.sentiment, 0) / responses.length
    : 0;

  if (loading) {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        </div>
      </DashboardLayout>
    );
  }

  if (!survey) {
    return (
      <DashboardLayout>
        <div className="text-center py-12">
          <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-medium mb-2">Survey not found</h3>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-foreground">{survey.title}</h1>
          <p className="text-muted-foreground mt-1">
            Survey Results • {responses.length} responses
          </p>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Total Responses
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{responses.length}</div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Average Sentiment
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-2">
                <div className="text-3xl font-bold">{avgSentiment.toFixed(2)}</div>
                {getSentimentIcon(avgSentiment)}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Overall Mood
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className={`inline-flex px-3 py-1 rounded-full text-sm font-medium ${getSentimentColor(avgSentiment)}`}>
                {getSentimentLabel(avgSentiment)}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Responses */}
        <Card>
          <CardHeader>
            <CardTitle>Individual Responses</CardTitle>
            <CardDescription>
              Detailed responses from each participant
            </CardDescription>
          </CardHeader>
          <CardContent>
            {responses.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                No responses yet
              </div>
            ) : (
              <div className="space-y-6">
                {responses.map((response, idx) => (
                  <div key={response.id} className="border rounded-lg p-6 space-y-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center">
                          <User className="h-5 w-5 text-primary" />
                        </div>
                        <div>
                          <p className="font-medium">{response.user?.name || 'Anonymous'}</p>
                          <p className="text-sm text-muted-foreground">{response.user?.email}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-4">
                        <div className={`flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium ${getSentimentColor(response.sentiment)}`}>
                          {getSentimentIcon(response.sentiment)}
                          {getSentimentLabel(response.sentiment)}
                        </div>
                        <div className="flex items-center gap-1 text-sm text-muted-foreground">
                          <Calendar className="h-4 w-4" />
                          {new Date(response.created_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>

                    <div className="space-y-3 pl-13">
                      {survey.questions.map((question, qIdx) => (
                        <div key={qIdx} className="space-y-1">
                          <p className="text-sm font-medium text-muted-foreground">
                            Q{qIdx + 1}: {question}
                          </p>
                          <p className="text-sm pl-4 border-l-2 border-primary/20 py-1">
                            {response.responses[qIdx] || 'No answer'}
                          </p>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  );
};

export default SurveyResults;