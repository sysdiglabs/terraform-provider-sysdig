package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Layout struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	W       int `json:"w"`
	H       int `json:"h"`
	PanelID int `json:"panelId"`
}

type QueryParams struct {
	Severities    []interface{} `json:"severities"`
	AlertStatuses []interface{} `json:"alertStatuses"`
	Categories    []interface{} `json:"categories"`
	Filter        string        `json:"filter"`
	TeamScope     bool          `json:"teamScope"`
}

type EventDisplaySettings struct {
	Enabled     bool        `json:"enabled"`
	QueryParams QueryParams `json:"queryParams"`
}

type Bottom struct {
	Enabled bool `json:"enabled"`
}

type Left struct {
	Enabled        bool        `json:"enabled"`
	DisplayName    interface{} `json:"displayName"`
	Unit           string      `json:"unit"`
	DisplayFormat  string      `json:"displayFormat"`
	Decimals       interface{} `json:"decimals"`
	MinValue       int         `json:"minValue"`
	MaxValue       interface{} `json:"maxValue"`
	MinInputFormat string      `json:"minInputFormat"`
	MaxInputFormat string      `json:"maxInputFormat"`
	Scale          string      `json:"scale"`
}

type Right struct {
	Enabled        bool        `json:"enabled"`
	DisplayName    interface{} `json:"displayName"`
	Unit           string      `json:"unit"`
	DisplayFormat  string      `json:"displayFormat"`
	Decimals       interface{} `json:"decimals"`
	MinValue       int         `json:"minValue"`
	MaxValue       interface{} `json:"maxValue"`
	MinInputFormat string      `json:"minInputFormat"`
	MaxInputFormat string      `json:"maxInputFormat"`
	Scale          string      `json:"scale"`
}

type AxesConfiguration struct {
	Bottom Bottom `json:"bottom"`
	Left   Left   `json:"left"`
	Right  Right  `json:"right"`
}

type LegendConfiguration struct {
	Enabled     bool        `json:"enabled"`
	Position    string      `json:"position"`
	Layout      string      `json:"layout"`
	ShowCurrent bool        `json:"showCurrent"`
	Width       interface{} `json:"width"`
	Height      interface{} `json:"height"`
}

type DisplayInfo struct {
	DisplayName                   string `json:"displayName"`
	TimeSeriesDisplayNameTemplate string `json:"timeSeriesDisplayNameTemplate"`
	Type                          string `json:"type"`
	Color                         string `json:"color,omitempty"`
	LineWidth                     int    `json:"lineWidth,omitempty"`
}

type FormatUnit string

const (
	FormatUnitPercentage FormatUnit = "%"
	FormatUnitData       FormatUnit = "byte"
	FormatUnitDataRate   FormatUnit = "byteRate"
	FormatUnitNumber     FormatUnit = "number"
	FormatUnitNumberRate FormatUnit = "numberRate"
	FormatUnitTime       FormatUnit = "relativeTime"
)

type Format struct {
	Unit          FormatUnit  `json:"unit"`
	InputFormat   string      `json:"inputFormat"`
	DisplayFormat string      `json:"displayFormat"`
	Decimals      interface{} `json:"decimals"`
	YAxis         string      `json:"yAxis"`
}

func newPercentFormat() Format {
	return Format{
		Unit:          FormatUnitPercentage,
		InputFormat:   "0-100",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func newDataFormat() Format {
	return Format{
		Unit:          FormatUnitData,
		InputFormat:   "B",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func newDataRateFormat() Format {
	return Format{
		Unit:          FormatUnitDataRate,
		InputFormat:   "B/s",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func newNumberFormat() Format {
	return Format{
		Unit:          FormatUnitNumber,
		InputFormat:   "1",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func newNumberRateFormat() Format {
	return Format{
		Unit:          FormatUnitNumberRate,
		InputFormat:   "/s",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

func newTimeFormat() Format {
	return Format{
		Unit:          FormatUnitTime,
		InputFormat:   "ns",
		DisplayFormat: "auto",
		Decimals:      nil,
		YAxis:         "auto",
	}
}

type AdvancedQueries struct {
	Enabled     bool        `json:"enabled"`
	DisplayInfo DisplayInfo `json:"displayInfo"`
	Format      Format      `json:"format"`
	Query       string      `json:"query"`
	ID          int         `json:"id"`
	ParentPanel *Panels     `json:"-"`
}

func NewPromqlQuery(query string, parentPanel *Panels) *AdvancedQueries {
	newQuery := &AdvancedQueries{
		Enabled: true,
		DisplayInfo: DisplayInfo{
			DisplayName:                   "",
			TimeSeriesDisplayNameTemplate: "",
			Type:                          "lines",
		},
		Format:      newPercentFormat(),
		Query:       query,
		ID:          0,
		ParentPanel: parentPanel,
	}

	if parentPanel.Type == "advancedNumber" {
		newQuery.DisplayInfo.Color = "mixed"
		newQuery.DisplayInfo.LineWidth = 2
	}

	return newQuery
}

func (q *AdvancedQueries) Enable(val bool) *AdvancedQueries {
	q.Enabled = val
	return q
}

func (q *AdvancedQueries) WithPercentFormat() *AdvancedQueries {
	q.Format = newPercentFormat()
	return q
}

func (q *AdvancedQueries) WithDataFormat() *AdvancedQueries {
	q.Format = newDataFormat()
	return q
}

func (q *AdvancedQueries) WithDataRateFormat() *AdvancedQueries {
	q.Format = newDataRateFormat()
	return q
}

func (q *AdvancedQueries) WithNumberFormat() *AdvancedQueries {
	q.Format = newNumberFormat()
	return q
}

func (q *AdvancedQueries) WithNumberRateFormat() *AdvancedQueries {
	q.Format = newNumberRateFormat()
	return q
}

func (q *AdvancedQueries) WithTimeFormat() *AdvancedQueries {
	q.Format = newTimeFormat()
	return q
}

type Panels struct {
	ID                     int                  `json:"id"`
	Name                   string               `json:"name"`
	Description            string               `json:"description"`
	AxesConfiguration      *AxesConfiguration   `json:"axesConfiguration,omitempty"`
	LegendConfiguration    *LegendConfiguration `json:"legendConfiguration,omitempty"`
	ApplyScopeToAll        bool                 `json:"applyScopeToAll,omitempty"`
	ApplySegmentationToAll bool                 `json:"applySegmentationToAll,omitempty"`
	AdvancedQueries        []*AdvancedQueries   `json:"advancedQueries,omitempty"`
	NumberThresholds       *NumberThresholds    `json:"numberThresholds,omitempty"`
	MarkdownSource         *string               `json:"markdownSource,omitempty"`
	PanelTitleVisible      bool                 `json:"panelTitleVisible"`
	TextAutosized          bool                 `json:"textAutosized"`
	TransparentBackground  bool                 `json:"transparentBackground"`
	Type                   PanelType            `json:"type"`
	// Just a helper to the client, the actual field is in Dashboard
	Layout *Layout `json:"-"`
}



type PanelType string

const (
	PanelTypeTimechart PanelType = "advancedTimechart"
	PanelTypeNumber    PanelType = "advancedNumber"
	PanelTypeText      PanelType = "text"
)

func (p *Panels) AddQueries(queries ...*AdvancedQueries) (*Panels, error) {
	if p.Type == PanelTypeNumber && len(p.AdvancedQueries)+len(queries) > 1 {
		return nil, fmt.Errorf("a panel of type 'number' can only contain one query")
	}

	maxIndex := 0
	for _, query := range p.AdvancedQueries {
		if maxIndex < query.ID {
			maxIndex = query.ID
		}
	}

	for _, query := range queries {
		maxIndex++
		query.ID = maxIndex
		p.AdvancedQueries = append(p.AdvancedQueries, query)
	}

	return p, nil
}

func (p *Panels) WithLayout(xPos, yPos, width, height int) (*Panels, error) {
	if xPos < 0 {
		return p, errors.New("x position must be at least 0")
	}

	if yPos < 0 {
		return p, errors.New("y position must be at least 0")
	}

	if xPos+width > 24 {
		return p, errors.New("the sum of the x position and the width must be lower or equal to 24")
	}

	// no limit in the height

	p.Layout = &Layout{
		X:       xPos,
		Y:       yPos,
		W:       width,
		H:       height,
		PanelID: p.ID,
	}

	return p, nil
}

type NumberThresholds struct {
	Base   Base          `json:"base"`
	Values []interface{} `json:"values"`
}

type Base struct {
	DisplayText string `json:"displayText"`
	Severity    string `json:"severity"`
}

type TeamSharingOptions struct {
	Type          string        `json:"type"`
	UserTeamsRole string        `json:"userTeamsRole"`
	SelectedTeams []interface{} `json:"selectedTeams"`
}
type Dashboard struct {
	Version                 int                  `json:"version,omitempty"`
	CustomerID              interface{}          `json:"customerId"`
	TeamID                  int                  `json:"teamId"`
	Schema                  int                  `json:"schema"`
	AutoCreated             bool                 `json:"autoCreated"`
	PublicToken             string               `json:"publicToken"`
	ScopeExpressionList     interface{}          `json:"scopeExpressionList"`
	Layout                  []*Layout            `json:"layout"`
	TeamScope               interface{}          `json:"teamScope"`
	EventDisplaySettings    EventDisplaySettings `json:"eventDisplaySettings"`
	ID                      int                  `json:"id,omitempty"`
	Name                    string               `json:"name"`
	Description             string               `json:"description"`
	Username                string               `json:"username"`
	Shared                  bool                 `json:"shared"`
	SharingSettings         []interface{}        `json:"sharingSettings"`
	Public                  bool                 `json:"public"`
	Favorite                bool                 `json:"favorite"`
	CreatedOn               int64                `json:"createdOn"`
	ModifiedOn              int64                `json:"modifiedOn"`
	Panels                  []*Panels            `json:"panels"`
	TeamScopeExpressionList []interface{}        `json:"teamScopeExpressionList"`
	CreatedOnDate           string               `json:"createdOnDate"`
	ModifiedOnDate          string               `json:"modifiedOnDate"`
	TeamSharingOptions      TeamSharingOptions   `json:"teamSharingOptions"`
}

type dashboardWrapper struct {
	Dashboard *Dashboard `json:"dashboard"`
}

func (db *Dashboard) ToJSON() io.Reader {
	payload, _ := json.Marshal(dashboardWrapper{db})
	return bytes.NewBuffer(payload)
}

func (db *Dashboard) AddPanels(panels ...*Panels) {
	maxPanelID := 0
	for _, existingPanel := range db.Panels {
		if maxPanelID < existingPanel.ID {
			maxPanelID = existingPanel.ID
		}
	}

	for _, panelToAdd := range panels {
		maxPanelID++
		panelToAdd.ID = maxPanelID
		panelToAdd.Layout.PanelID = maxPanelID

		db.Panels = append(db.Panels, panelToAdd)
		db.Layout = append(db.Layout, panelToAdd.Layout)
	}
}

func NewDashboard(name, description string) *Dashboard {
	return &Dashboard{
		Name:        name,
		Description: description,
		Schema:      3,
	}
}

func (db *Dashboard) AsPublic(value bool) *Dashboard {
	db.Public = value
	return db
}

func DashboardFromJSON(body []byte) *Dashboard {
	var result dashboardWrapper
	json.Unmarshal(body, &result)

	return result.Dashboard
}
