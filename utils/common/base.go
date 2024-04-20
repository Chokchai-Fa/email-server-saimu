package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DateCommon struct {
	CreateBy   string `json:"create_by" dynamodbav:"create_by" db:"create_by" differ:"exclude"`
	CreateDate string `json:"create_date" dynamodbav:"create_date" db:"create_date" differ:"exclude"`
	UpdateBy   string `json:"update_by" dynamodbav:"update_by,omitempty" db:"update_by,omitempty" differ:"exclude"`
	UpdateDate string `json:"update_date" dynamodbav:"update_date,omitempty" db:"update_date,omitempty" differ:"exclude"`
}

type DateCommonPostgres struct {
	CreateBy   NullString `json:"create_by" db:"create_by"`
	CreateDate NullTime   `json:"create_date" db:"create_date"`
	UpdateBy   NullString `json:"update_by"  db:"update_by,omitempty"`
	UpdateDate NullTime   `json:"update_date"  db:"update_date,omitempty"`
}

type VersioningCommon struct {
	Active        bool   `json:"active,omitempty"  dynamodbav:"active" differ:"exclude"`
	Version       uint16 `json:"version,omitempty" dynamodbav:"version,omitempty" differ:"exclude"`
	LatestVersion bool   `json:"latest_version,omitempty" dynamodbav:"latest_version,omitempty"`
	CreateDate    string `json:"create_date" dynamodbav:"create_date" differ:"exclude"`
	UpdateDate    string `json:"update_date" dynamodbav:"update_date,omitempty"`
	DeleteDate    string `json:"delete_date,omitempty" dynamodbav:"delete_date,omitempty"` // if item is not deleted this field is empty
}

type Key struct {
	ObjectID string `json:"object_id" copier:"objectid" dynamodbav:"object_id"` //PK
	ID       string `json:"id" dynamodbav:"id"`                                 //FK
}

type KeyReq struct {
	ObjectID string `json:"object_id" copier:"objectid"  validate:"required"` //PK
	ID       string `json:"id"  validate:"required"`
}

type MultiLang struct {
	ThTH string `json:"th_TH" dynamodbav:"th_TH,omitempty" db:"th_TH"`
	EnUS string `json:"en_US" dynamodbav:"en_US,omitempty" db:"en_US"`
}

// Make the MultiLang struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (m MultiLang) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Make the MultiLang struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (m *MultiLang) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		fmt.Println("type assertion to []byte failed")
		return nil
		// return errors.New("type assertion to []byte failed") // comment because jsonb field can null value
	}

	return json.Unmarshal(b, &m)
}

type MultiLangReq struct {
	ThTH string `json:"th_TH" dynamodbav:"th_TH,omitempty" validate:"required,notblank"`
	EnUS string `json:"en_US" dynamodbav:"en_US,omitempty" validate:"required,notblank"`
}

// Make the MultiLangReq struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (m MultiLangReq) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Make the MultiLangReq struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (m *MultiLangReq) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		fmt.Println("type assertion to []byte failed")
		return nil
		// return errors.New("type assertion to []byte failed") // comment because jsonb field can null value
	}

	return json.Unmarshal(b, &m)
}

type SaleStatus struct {
	Status    bool   `json:"status" dynamodbav:"status,omitempty"`
	StartDate string `json:"start_date,omitempty" dynamodbav:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty" dynamodbav:"end_date,omitempty"`
}

type MultiLanguages map[string]string

func (m MultiLanguages) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MultiLanguages) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		fmt.Println("type assertion to []byte failed")
		return nil
		// return errors.New("type assertion to []byte failed") // comment because jsonb field can null value
	}

	return json.Unmarshal(b, &m)
}

type TranslationTable struct {
	ID           int
	Key          string
	Collection   string
	LanguageCode string
	Name         string
	Description  string
}

func (t *TranslationTable) MapToTranslationsTable(name, description map[string]string, tableName, planType string) []TranslationTable {
	translations := []TranslationTable{}
	for key, value := range name {
		lang := key
		name := value
		description := description[key]

		translation := TranslationTable{
			Key:          planType,
			Collection:   tableName,
			LanguageCode: lang,
			Name:         name,
			Description:  description,
		}
		translations = append(translations, translation)
	}
	return translations
}
