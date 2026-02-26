package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    GoogleID  string    `gorm:"uniqueIndex;not null" json:"google_id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Phone     string    `json:"phone"`
    AvatarURL string    `json:"avatar_url"`
    IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    Submissions []Submission `json:"submissions,omitempty"`
}

type Test struct {
    ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    Title       string          `gorm:"not null" json:"title"`
    Type        string          `gorm:"not null;check:type IN ('Reading','Listening','Writing')" json:"type"`
    Content     JSONMap         `gorm:"type:jsonb;not null" json:"content"`
    Answers     JSONMap         `gorm:"type:jsonb;not null" json:"-"` // Hidden from JSON
    IsPublished bool            `gorm:"default:false" json:"is_published"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    
    Submissions []Submission    `json:"submissions,omitempty"`
}

type Submission struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    TestID      uuid.UUID `gorm:"type:uuid;not null;index" json:"test_id"`
    UserAnswers JSONMap   `gorm:"type:jsonb;not null" json:"user_answers"`
    Score       *int      `json:"score"`
    Status      string    `gorm:"default:'completed'" json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Test        Test      `gorm:"foreignKey:TestID" json:"test,omitempty"`
}

// Custom type for JSON maps with proper GORM implementations
type JSONMap map[string]interface{}

// FIXED: Implement Valuer interface for JSONMap
func (jm JSONMap) Value() (driver.Value, error) {
    if jm == nil {
        return nil, nil
    }
    return json.Marshal(jm)
}

// FIXED: Implement Scanner interface for JSONMap
func (jm *JSONMap) Scan(value interface{}) error {
    if value == nil {
        *jm = nil
        return nil
    }
    
    var bytes []byte
    switch v := value.(type) {
    case []byte:
        bytes = v
    case string:
        bytes = []byte(v)
    default:
        return errors.New("failed to unmarshal JSONB value: invalid type")
    }
    
    result := make(map[string]interface{})
    err := json.Unmarshal(bytes, &result)
    if err != nil {
        return err
    }
    
    *jm = result
    return nil
}

// BeforeCreate hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

func (t *Test) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
}

func (s *Submission) BeforeCreate(tx *gorm.DB) error {
    if s.ID == uuid.Nil {
        s.ID = uuid.New()
    }
    return nil
}
