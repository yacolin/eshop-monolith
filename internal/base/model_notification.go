package base

import (
	"time"

	"gorm.io/gorm"

	"eshop-monolith/pkg/utils"
)

type Notification struct {
	ID              int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64           `gorm:"not null;index:idx_user_id" json:"user_id"`
	MerchantID      int64           `gorm:"not null;default:0" json:"merchant_id"`
	Title           string          `gorm:"type:varchar(200);not null" json:"title"`
	Content         string          `gorm:"type:text;not null" json:"content"`
	ContentTemplate string          `gorm:"type:varchar(100)" json:"content_template,omitempty"`
	TemplateParams  *string         `gorm:"type:json" json:"template_params,omitempty"`
	Channel         int8            `gorm:"type:tinyint;not null;index:idx_channel" json:"channel"`
	Category        int8            `gorm:"type:tinyint;not null;index:idx_category" json:"category"`
	TargetType      string          `gorm:"type:varchar(30)" json:"target_type,omitempty"`
	TargetID        *int64          `json:"target_id,omitempty"`
	RedirectURL     string          `gorm:"type:varchar(500)" json:"redirect_url,omitempty"`
	IconURL         string          `gorm:"type:varchar(500)" json:"icon_url,omitempty"`
	Priority        int8            `gorm:"type:tinyint;not null;default:1" json:"priority"`
	CreatedBy       int64           `gorm:"not null;default:0" json:"created_by"`
	CreatedAt       utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);index:idx_created" json:"created_at"`
	UpdatedAt       utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"type:datetime(3)" json:"-"`

	// computed: LEFT JOIN base_notification_reads, 只读不写入
	IsRead bool `gorm:"->;default:0" json:"is_read"`
}

func (Notification) TableName() string { return "base_notifications" }

// NotificationRead 已读记录
type NotificationRead struct {
	ID             int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	NotificationID int64           `gorm:"not null;uniqueIndex:uk_notification_user" json:"notification_id"`
	UserID         int64           `gorm:"not null;uniqueIndex:uk_notification_user;index:idx_user" json:"user_id"`
	ReadAt         time.Time       `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"read_at"`
	CreatedAt      utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
}

func (NotificationRead) TableName() string { return "base_notification_reads" }

type NotificationTemplate struct {
	ID              int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateCode    string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"template_code"`
	Channel         int8            `gorm:"type:tinyint;not null;index:idx_code_channel" json:"channel"`
	TitleTemplate   string          `gorm:"type:varchar(200);not null" json:"title_template"`
	ContentTemplate string          `gorm:"type:text;not null" json:"content_template"`
	Category        int8            `gorm:"type:tinyint" json:"category"`
	Priority        int8            `gorm:"type:tinyint;not null;default:1" json:"priority"`
	Status          int8            `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt       utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt       utils.Timestamp `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
}

func (NotificationTemplate) TableName() string { return "base_notification_templates" }

const (
	ChannelInApp  int8 = 1
	ChannelPush   int8 = 2
	ChannelSMS    int8 = 3
	ChannelEmail  int8 = 4
	ChannelWechat int8 = 5
)

const (
	CategorySystem    int8 = 1
	CategoryOrder     int8 = 2
	CategoryMarketing int8 = 3
	CategoryInteract  int8 = 4
	CategorySecurity  int8 = 5
)
