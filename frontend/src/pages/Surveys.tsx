import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FileText, Plus, Calendar, Eye, Trash2, X } from 'lucide-react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { EmptyState } from '@/components/ui/empty-state';
import { surveyAPI } from '@/services/api';
import type { Survey } from '@/types';
import { toast } from 'sonner';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

export const Surveys = () => {
  const navigate = useNavigate();
  const [surveys, setSurveys] = useState<Survey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [newSurvey, setNewSurvey] = useState({
    title: '',
    questions: [''],
  });

  useEffect(() => {
    fetchSurveys();
  }, []);

  const fetchSurveys = async () => {
    try {
      const response = await surveyAPI.getSurveys();
      setSurveys(response);
    } catch (error: any) {
      toast.error('Failed to load surveys');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddQuestion = () => {
    setNewSurvey((prev) => ({
      ...prev,
      questions: [...prev.questions, ''],
    }));
  };

  const handleRemoveQuestion = (index: number) => {
    if (newSurvey.questions.length === 1) return;
    setNewSurvey((prev) => ({
      ...prev,
      questions: prev.questions.filter((_, i) => i !== index),
    }));
  };

  const handleQuestionChange = (index: number, value: string) => {
    setNewSurvey((prev) => ({
      ...prev,
      questions: prev.questions.map((q, i) => (i === index ? value : q)),
    }));
  };

  const handleCreateSurvey = async () => {
    if (!newSurvey.title.trim()) {
      toast.error('Please enter a survey title');
      return;
    }

    const validQuestions = newSurvey.questions.filter((q) => q.trim());
    if (validQuestions.length === 0) {
      toast.error('Please add at least one question');
      return;
    }

    setIsCreating(true);
    try {
      await surveyAPI.createSurvey({
        title: newSurvey.title,
        questions: validQuestions,
      });
      toast.success('Survey created successfully!');
      setIsCreateOpen(false);
      setNewSurvey({ title: '', questions: [''] });
      fetchSurveys();
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to create survey');
    } finally {
      setIsCreating(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Surveys</h1>
            <p className="text-muted-foreground mt-1">Create and manage employee surveys</p>
          </div>
          <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Create Survey
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle>Create New Survey</DialogTitle>
                <DialogDescription>
                  Create a survey to gather feedback from your employees.
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-6 py-4">
                <div className="space-y-2">
                  <Label htmlFor="title">Survey Title</Label>
                  <Input
                    id="title"
                    placeholder="e.g., Q4 Employee Wellbeing Survey"
                    value={newSurvey.title}
                    onChange={(e) =>
                      setNewSurvey((prev) => ({ ...prev, title: e.target.value }))
                    }
                  />
                </div>

                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <Label>Questions</Label>
                    <Button type="button" variant="outline" size="sm" onClick={handleAddQuestion}>
                      <Plus className="h-4 w-4 mr-1" />
                      Add Question
                    </Button>
                  </div>
                  <div className="space-y-3">
                    {newSurvey.questions.map((question, index) => (
                      <div key={index} className="flex gap-2">
                        <div className="flex-1">
                          <Textarea
                            placeholder={`Question ${index + 1}`}
                            value={question}
                            onChange={(e) => handleQuestionChange(index, e.target.value)}
                            rows={2}
                          />
                        </div>
                        {newSurvey.questions.length > 1 && (
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon"
                            onClick={() => handleRemoveQuestion(index)}
                            className="text-muted-foreground hover:text-destructive"
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        )}
                      </div>
                    ))}
                  </div>
                </div>

                <Button
                  className="w-full"
                  onClick={handleCreateSurvey}
                  disabled={isCreating}
                >
                  {isCreating ? (
                    <LoadingSpinner size="sm" />
                  ) : (
                    'Create Survey'
                  )}
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>

        {/* Surveys Grid */}
        {isLoading ? (
          <div className="flex items-center justify-center min-h-[400px]">
            <LoadingSpinner size="lg" text="Loading surveys..." />
          </div>
        ) : surveys.length > 0 ? (
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {surveys.map((survey, index) => (
              <div
                key={survey.id}
                className="bg-card rounded-xl border border-border p-6 shadow-sm hover:shadow-md transition-all duration-200 animate-slide-up"
                style={{ animationDelay: `${index * 50}ms` }}
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="w-12 h-12 rounded-xl gradient-primary flex items-center justify-center">
                    <FileText className="h-6 w-6 text-primary-foreground" />
                  </div>
                  <Badge variant={survey.is_active ? 'default' : 'secondary'}>
                    {survey.is_active ? 'Active' : 'Inactive'}
                  </Badge>
                </div>

                <h3 className="text-lg font-semibold text-foreground mb-2 line-clamp-2">
                  {survey.title}
                </h3>

                <div className="flex items-center gap-4 text-sm text-muted-foreground mb-4">
                  <span>{survey.questions.length} questions</span>
                  <span className="flex items-center gap-1">
                    <Calendar className="h-4 w-4" />
                    {formatDate(survey.created_at)}
                  </span>
                </div>

                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    className="flex-1"
                    onClick={() => navigate(`/surveys/${survey.id}/results`)}
                  >
                    <Eye className="h-4 w-4 mr-1" />
                    Results
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => toast.info('Edit coming soon!')}
                  >
                    Edit
                  </Button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <EmptyState
            icon={FileText}
            title="No surveys yet"
            description="Create your first survey to start gathering employee feedback and insights."
            action={
              <Button onClick={() => setIsCreateOpen(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Create Survey
              </Button>
            }
          />
        )}
      </div>
    </DashboardLayout>
  );
};

export default Surveys;
