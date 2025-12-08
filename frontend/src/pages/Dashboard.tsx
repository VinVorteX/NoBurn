import { useEffect, useState } from 'react';
import { Users, AlertTriangle, TrendingDown, Smile, FileText, Bell, Plus } from 'lucide-react';
import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { StatCard } from '@/components/ui/stat-card';
import { Button } from '@/components/ui/button';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { EmptyState } from '@/components/ui/empty-state';
import { analyticsAPI } from '@/services/api';
import type { DashboardData } from '@/types';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

export const Dashboard = () => {
  const [data, setData] = useState<DashboardData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const response = await analyticsAPI.getDashboard();
        setData(response);
      } catch (error: any) {
        toast.error('Failed to load dashboard data');
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboard();
  }, []);

  const getRiskBadgeVariant = (score: number) => {
    if (score >= 0.7) return 'destructive';
    if (score >= 0.4) return 'secondary';
    return 'outline';
  };

  const getRiskLabel = (score: number) => {
    if (score >= 0.7) return 'High Risk';
    if (score >= 0.4) return 'Medium Risk';
    return 'Low Risk';
  };

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
            <p className="text-muted-foreground mt-1">Monitor your team's wellbeing at a glance</p>
          </div>
          <div className="flex gap-3">
            <Button variant="outline" onClick={() => toast.info('Alerts feature coming soon!')}>
              <Bell className="h-4 w-4 mr-2" />
              Send Alerts
            </Button>
            <Button onClick={() => navigate('/surveys')}>
              <Plus className="h-4 w-4 mr-2" />
              Create Survey
            </Button>
          </div>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center min-h-[400px]">
            <LoadingSpinner size="lg" text="Loading dashboard..." />
          </div>
        ) : data ? (
          <>
            {/* Stats Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              <StatCard
                title="Total Employees"
                value={data.total_employees}
                icon={Users}
                variant="primary"
              />
              <StatCard
                title="At Risk Employees"
                value={data.at_risk_employees}
                icon={AlertTriangle}
                variant="destructive"
              />
              <StatCard
                title="Churn Rate"
                value={`${data.churn_rate.toFixed(1)}%`}
                icon={TrendingDown}
                variant="warning"
              />
              <StatCard
                title="Avg Sentiment"
                value={`${(data.avg_sentiment * 100).toFixed(0)}%`}
                icon={Smile}
                variant="success"
              />
            </div>

            <div className="grid lg:grid-cols-3 gap-6">
              {/* Risk Factors */}
              <div className="bg-card rounded-xl border border-border p-6 shadow-sm animate-slide-up">
                <h2 className="text-lg font-semibold text-foreground mb-4">Top Risk Factors</h2>
                {data.top_risk_factors.length > 0 ? (
                  <div className="space-y-3">
                    {data.top_risk_factors.map((factor, index) => (
                      <div
                        key={index}
                        className="flex items-center gap-3 p-3 rounded-lg bg-muted/50"
                      >
                        <div className="w-8 h-8 rounded-full bg-destructive/10 flex items-center justify-center">
                          <span className="text-sm font-semibold text-destructive">{index + 1}</span>
                        </div>
                        <span className="text-sm font-medium text-foreground">{factor}</span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">No risk factors identified</p>
                )}
              </div>

              {/* High Risk Employees Table */}
              <div className="lg:col-span-2 bg-card rounded-xl border border-border p-6 shadow-sm animate-slide-up">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-lg font-semibold text-foreground">High Risk Employees</h2>
                  <Button variant="ghost" size="sm" onClick={() => navigate('/employees')}>
                    View All
                  </Button>
                </div>
                {data.attrition_risks.length > 0 ? (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Employee</TableHead>
                        <TableHead>Email</TableHead>
                        <TableHead className="text-right">Risk Score</TableHead>
                        <TableHead className="text-right">Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {data.attrition_risks.slice(0, 5).map((risk) => (
                        <TableRow key={risk.id}>
                          <TableCell className="font-medium">{risk.user.name}</TableCell>
                          <TableCell className="text-muted-foreground">{risk.user.email}</TableCell>
                          <TableCell className="text-right font-semibold">
                            {(risk.risk_score * 100).toFixed(0)}%
                          </TableCell>
                          <TableCell className="text-right">
                            <Badge variant={getRiskBadgeVariant(risk.risk_score)}>
                              {getRiskLabel(risk.risk_score)}
                            </Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                ) : (
                  <EmptyState
                    icon={Users}
                    title="No high-risk employees"
                    description="Great news! None of your employees are currently at high risk of burnout."
                  />
                )}
              </div>
            </div>

            {/* Quick Actions */}
            <div className="bg-card rounded-xl border border-border p-6 shadow-sm animate-slide-up">
              <h2 className="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
              <div className="grid sm:grid-cols-3 gap-4">
                <Button
                  variant="outline"
                  className="h-auto py-4 flex flex-col items-center gap-2"
                  onClick={() => navigate('/surveys')}
                >
                  <FileText className="h-6 w-6 text-primary" />
                  <span>Create Survey</span>
                </Button>
                <Button
                  variant="outline"
                  className="h-auto py-4 flex flex-col items-center gap-2"
                  onClick={() => toast.info('Reports feature coming soon!')}
                >
                  <TrendingDown className="h-6 w-6 text-primary" />
                  <span>View Reports</span>
                </Button>
                <Button
                  variant="outline"
                  className="h-auto py-4 flex flex-col items-center gap-2"
                  onClick={() => toast.info('Alerts feature coming soon!')}
                >
                  <Bell className="h-6 w-6 text-primary" />
                  <span>Send Alerts</span>
                </Button>
              </div>
            </div>
          </>
        ) : (
          <EmptyState
            icon={AlertTriangle}
            title="Failed to load data"
            description="We couldn't load the dashboard data. Please try refreshing the page."
            action={
              <Button onClick={() => window.location.reload()}>Refresh Page</Button>
            }
          />
        )}
      </div>
    </DashboardLayout>
  );
};

export default Dashboard;
