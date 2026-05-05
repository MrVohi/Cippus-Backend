package models

type UserRole string
type ModerationStatus string
type ReportStatus string
type OAuthProviderName string
type NotificationType string
type ContentType string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"

	ModerationPending  ModerationStatus = "pending"
	ModerationApproved ModerationStatus = "approved"
	ModerationFlagged  ModerationStatus = "flagged"
	ModerationBlocked  ModerationStatus = "blocked"

	ReportPending   ReportStatus = "pending"
	ReportResolved  ReportStatus = "resolved"
	ReportDismissed ReportStatus = "dismissed"

	ProviderGithub OAuthProviderName = "github"
	ProviderGoogle OAuthProviderName = "google"

	NotificationReply    NotificationType = "reply"
	NotificationLike     NotificationType = "like"
	NotificationMention  NotificationType = "mention"
	NotificationDM       NotificationType = "dm"
	NotificationFlagged  NotificationType = "post_flagged"
	NotificationApproved NotificationType = "post_approved"

	ContentPost    ContentType = "post"
	ContentComment ContentType = "comment"
)
