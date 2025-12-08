import { useEffect, useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { FileText, CheckCircle, AlertCircle, Flame } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { surveyAPI } from '@/services/api';
import type { PublicSurvey as PublicSurveyType } from '@/types';
import { toast } from 'sonner';

export const PublicSurvey = () => {
  const { surveyId } = useParams();
  const [searchParams] = useSearchParams();
  const userToken = searchParams.get('token');

  const [survey, setSurvey] = useState<PublicSurveyType | null>(null);
  const [responses, setResponses] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSurvey = async () => {
      if (!surveyId) {
        setError('Invalid survey link');
        setIsLoading(false);
        return;
      }

      try {
        const response = await surveyAPI.getPublicSurvey(parseInt(surveyId));
        setSurvey(response);
        setResponses(new Array(response.questions.length).fill(''));
      } catch (error: any) {
        setError(error.response?.data?.error || 'Survey not found');
      } finally {
        setIsLoading(false);
      }
    };

    fetchSurvey();
  }, [surveyId]);

  const handleResponseChange = (index: number, value: string) => {
    setResponses((prev) => prev.map((r, i) => (i === index ? value : r)));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!surveyId || !userToken) {
      toast.error('Invalid survey link');
      return;
    }

    const emptyResponses = responses.filter((r) => !r.trim());
    if (emptyResponses.length > 0) {
      toast.error('Please answer all questions');
      return;
    }

    setIsSubmitting(true);
    try {
      await surveyAPI.submitPublicResponse({
        survey_id: parseInt(surveyId),
        user_token: parseInt(userToken),
        responses,
      });
      setIsSubmitted(true);
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to submit survey');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <LoadingSpinner size="lg" text="Loading survey..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center p-4">
        <div className="text-center animate-fade-in">
          <div className="w-16 h-16 rounded-full bg-destructive/10 flex items-center justify-center mx-auto mb-4">
            <AlertCircle className="h-8 w-8 text-destructive" />
          </div>
          <h1 className="text-2xl font-bold text-foreground mb-2">Survey Not Found</h1>
          <p className="text-muted-foreground">{error}</p>
        </div>
      </div>
    );
  }

  if (isSubmitted) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center p-4">
        <div className="text-center animate-fade-in max-w-md">
          <div className="w-20 h-20 rounded-full bg-success/10 flex items-center justify-center mx-auto mb-6">
            <CheckCircle className="h-10 w-10 text-success" />
          </div>
          <h1 className="text-3xl font-bold text-foreground mb-3">Thank You!</h1>
          <p className="text-lg text-muted-foreground mb-6">
            Your responses have been submitted successfully. We appreciate your feedback!
          </p>
          <div className="p-4 rounded-lg bg-muted/50">
            <p className="text-sm text-muted-foreground">
              Your honest feedback helps us create a better work environment for everyone.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="bg-card border-b border-border">
        <div className="max-w-2xl mx-auto px-4 py-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl gradient-primary flex items-center justify-center">
              <Flame className="h-5 w-5 text-primary-foreground" />
            </div>
            <span className="text-lg font-semibold text-foreground">NoBurn HR</span>
          </div>
        </div>
      </header>

      {/* Survey Content */}
      <main className="max-w-2xl mx-auto px-4 py-8">
        <div className="animate-slide-up">
          {/* Survey Header */}
          <div className="bg-card rounded-xl border border-border p-6 mb-6 shadow-sm">
            <div className="w-14 h-14 rounded-xl gradient-primary flex items-center justify-center mb-4">
              <FileText className="h-7 w-7 text-primary-foreground" />
            </div>
            <h1 className="text-2xl font-bold text-foreground mb-2">{survey?.title}</h1>
            <p className="text-muted-foreground">
              Please take a few minutes to complete this survey. Your responses are anonymous and
              will help us improve the workplace.
            </p>
          </div>

          {/* Questions */}
          <form onSubmit={handleSubmit} className="space-y-6">
            {survey?.questions.map((question, index) => (
              <div
                key={index}
                className="bg-card rounded-xl border border-border p-6 shadow-sm animate-fade-in"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <Label className="text-base font-medium text-foreground mb-3 block">
                  <span className="text-primary mr-2">{index + 1}.</span>
                  {question}
                </Label>
                <Textarea
                  placeholder="Type your answer here..."
                  value={responses[index] || ''}
                  onChange={(e) => handleResponseChange(index, e.target.value)}
                  rows={4}
                  className="mt-3"
                  required
                />
              </div>
            ))}

            <Button type="submit" size="lg" className="w-full h-14 text-base" disabled={isSubmitting}>
              {isSubmitting ? <LoadingSpinner size="sm" /> : 'Submit Survey'}
            </Button>
          </form>
        </div>
      </main>
    </div>
  );
};

export default PublicSurvey;
