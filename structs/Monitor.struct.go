package structs

import "time"

type MonitorType string

const (
	MonitorTypeHTTP MonitorType = "http"
	MonitorTypeTCP  MonitorType = "tcp"
	MonitorTypeICMP MonitorType = "icmp"
	MonitorTypeDNS  MonitorType = "dns"
	MonitorTypeTLS  MonitorType = "tls"
)

type Monitor struct {
	ID           int         `db:"id" json:"id"`
	Name         string      `db:"name" json:"name"`
	Type         MonitorType `db:"type" json:"type"`
	Endpoint     string      `db:"endpoint" json:"endpoint"` // URL, host:port, or domain
	IntervalS    int         `db:"interval_s" json:"interval_s"`
	TimeoutMS    int         `db:"timeout_ms" json:"timeout_ms"`
	ExpectedCode *int        `db:"expected_code" json:"expected_code,omitempty"`
	Enabled      bool        `db:"enabled" json:"enabled"`
	CreatedBy    int         `db:"created_by" json:"created_by"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
	InsertedAt   time.Time   `db:"inserted_at" json:"inserted_at"`
}

type MonitorResult struct {
	ID         int       `db:"id" json:"id"`
	MonitorID  int       `db:"monitor_id" json:"monitor_id"`
	StartedAt  time.Time `db:"started_at" json:"started_at"`
	DurationMS int       `db:"duration_ms" json:"duration_ms"`
	Status     string    `db:"status" json:"status"`
	Code       *int      `db:"code" json:"code,omitempty"`
	Message    *string   `db:"message" json:"message,omitempty"`
}
