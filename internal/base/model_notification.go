package base

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64          `gorm:"not null;index:idx_user_read" json:"user_id"`
	Title           string         `gorm:"type:varchar(200);not null" json:"title"`
	Content         string         `gorm:"type:text;not null" json:"content"`
	ContentTemplate string         `gorm:"type:varchar(100)" json:"content_template,omitempty"`
	TemplateParams  string         `gorm:"type:json" json:"template_params,omitempty"`
	Channel         int8           `gorm:"type:tinyint;not null;index:idx_channel" json:"channel"`
	Category        int8           `gorm:"type:tinyint;not null;index:idx_category" json:"category"`
	TargetType      string         `gorm:"type:varchar(30)" json:"target_type,omitempty"`
	TargetID        *int64         `gorm:"" json:"target_id,omitempty"`
	RedirectURL     string         `gorm:"type:varchar(500)" json:"redirect_url,omitempty"`
	IconURL         string         `gorm:"type:varchar(500)" json:"icon_url,omitempty"`
	IsRead          bool           `gorm:"not null;default:false" json:"is_read"`
	ReadAt          *time.Time     `gorm:"type:datetime(3)" json:"read_at,omitempty"`
	IsProcessed     bool           `gorm:"not null;default:false;index:idx_processed" json:"is_processed"`
	ProcessedAt     *time.Time     `gorm:"type:datetime(3)" json:"processed_at,omitempty"`
	ProcessResult   string         `gorm:"type:varchar(200)" json:"process_result,omitempty"`
	IsDeletedByUser bool           `gorm:"not null;default:false" json:"is_deleted_by_user"`
	Priority        int8           `gorm:"type:tinyint;not null;default:1" json:"priority"`
	CreatedBy       int64          `gorm:"not null;default:0" json:"created_by"`
	CreatedAt       time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);index:idx_created" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"type:datetime(3)" json:"-"`
}

func (Notification) TableName() string { return "base_notifications" }

type NotificationTemplate struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateCode    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"template_code"`
	Channel         int8           `gorm:"type:tinyint;not null;index:idx_code_channel" json:"channel"`
	TitleTemplate   string         `gorm:"type:varchar(200);not null" json:"title_template"`
	ContentTemplate string         `gorm:"type:text;not null" json:"content_template"`
	Category        int8           `gorm:"type:tinyint" json:"category"`
	Priority        int8           `gorm:"type:tinyint;not null;default:1" json:"priority"`
	Status          int8           `gorm:"type:tinyint;not null;default:1" json:"status"`
	CreatedAt       time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);onUpdate:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"type:datetime(3)" json:"-"`
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
