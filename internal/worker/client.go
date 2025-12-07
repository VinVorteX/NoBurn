package worker

import (
	"time"

	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
}

func NewClient(redisAddr string) *Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &Client{client: client}
}

func (c *Client) EnqueueSurveyProcessing(responseID, userID uint, language string) error {
	task, err := NewProcessSurveyTask(responseID, userID, language)
	if err != nil {
		return err
	}

	_, err = c.client.Enqueue(task, asynq.Queue("default"))
	return err
}

func (c *Client) EnqueueChurnCalculation(userID, companyID uint) error {
	task, err := NewCalculateChurnTask(userID, companyID)
	if err != nil {
		return err
	}

	_, err = c.client.Enqueue(task, asynq.Queue("default"), asynq.ProcessIn(5*time.Minute))
	return err
}

func (c *Client) EnqueueNotification(userID uint, notificationType, message string) error {
	task, err := NewSendNotificationTask(userID, notificationType, message)
	if err != nil {
		return err
	}

	_, err = c.client.Enqueue(task, asynq.Queue("critical"))
	return err
}

func (c *Client) EnqueueSurveyInvitation(surveyID, userID uint, email string) error {
	task, err := NewSurveyInvitationTask(surveyID, userID, email)
	if err != nil {
		return err
	}

	_, err = c.client.Enqueue(task, asynq.Queue("default"))
	return err
}

func (c *Client) ScheduleDailyChurnAnalysis(companyID uint) error {
	// Schedule daily churn analysis for all users in company
	task, err := NewCalculateChurnTask(0, companyID) // 0 means all users
	if err != nil {
		return err
	}

	_, err = c.client.Enqueue(task, 
		asynq.Queue("low"),
		asynq.ProcessAt(time.Now().Add(24*time.Hour)),
	)
	return err
}

func (c *Client) Close() error {
	return c.client.Close()
}