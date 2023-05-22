package v2

import (
	"errors"
	"fmt"
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
	MinValue       float64     `json:"minValue"`
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
	MinValue       float64     `json:"minValue"`
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
	Unit                 FormatUnit `json:"unit"`
	InputFormat          *string    `json:"inputFormat"`
	DisplayFormat        *string    `json:"displayFormat"`
	Decimals             *int       `json:"decimals"`
	YAxis                *string    `json:"yAxis"`
	MinInterval          *string    `json:"minInterval"`
	NullValueDisplayMode *string    `json:"nullValueDisplayMode"`
}

func NewFormat(
	unit FormatUnit,
	inputFormat string,
	displayFormat string,
	decimals int,
	yAxis string,
	minInterval string,
	nullValueDisplayMode string) *Format {
	return &Format{
		Unit:                 unit,
		InputFormat:          &inputFormat,
		DisplayFormat:        &displayFormat,
		Decimals:             &decimals,
		YAxis:                &yAxis,
		MinInterval:          &minInterval,
		NullValueDisplayMode: &nullValueDisplayMode,
	}
}

func newPercentFormat() Format {
	return *NewFormat(
		FormatUnitPercentage,
		"0-100",
		"auto",
		0,
		"auto",
		"",
		"nullGap",
	)
}

func newDataFormat() Format {
	return *NewFormat(
		FormatUnitData,
		"B",
		"auto",
		0,
		"auto",
		"",
		"nullGap",
	)
}

func newDataRateFormat() Format {
	return *NewFormat(
		FormatUnitDataRate,
		"B/s",
		"auto",
		0,
		"auto",
		"",
		"nullGap",
	)
}

func newNumberFormat() Format {
	return *NewFormat(
		FormatUnitNumber,
		"1",
		"auto",
		0,
		"auto",
		"",
		"nullGap",
	)
}

func newNumberRateFormat() Format {
	return *NewFormat(
		FormatUnitNumberRate,
		"/s",
		"auto",
		0,
		"auto",
		"",
		"nullGap",
	)
}

func newTimeFormat() Format {
	return *NewFormat(
		FormatUnitTime,
		"ns",
		"auto",
		0,
		"auto",
		"",
		"nullGap",
	)
}

type AdvancedQueries struct {
	Enabled     bool        `json:"enabled"`
	DisplayInfo DisplayInfo `json:"displayInfo"`
	Format      Format      `json:"format"`
	Query       string      `json:"query"`
	ID          int         `json:"id"`
	ParentPanel *Panels     `json:"-"`
}

func NewPromqlQuery(query string, parentPanel *Panels, displayInfo DisplayInfo) *AdvancedQueries {
	newQuery := &AdvancedQueries{
		Enabled: true,
		DisplayInfo: DisplayInfo{
			DisplayName:                   displayInfo.DisplayName,
			TimeSeriesDisplayNameTemplate: displayInfo.TimeSeriesDisplayNameTemplate,
			Type:                          displayInfo.Type,
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

func (q *AdvancedQueries) updateFormat(f *Format) {
	if f == nil {
		return
	}

	if f.Unit != "" {
		q.Format.Unit = f.Unit
	}

	if f.DisplayFormat != nil {
		q.Format.DisplayFormat = f.DisplayFormat
	}

	if f.InputFormat != nil {
		q.Format.InputFormat = f.InputFormat
	}

	if f.Decimals != nil {
		q.Format.Decimals = f.Decimals
	}

	if f.YAxis != nil {
		q.Format.YAxis = f.YAxis
	}

	if f.NullValueDisplayMode != nil {
		q.Format.NullValueDisplayMode = f.NullValueDisplayMode
	}

	if f.MinInterval != nil {
		q.Format.MinInterval = f.MinInterval
	}
}

func (q *AdvancedQueries) WithPercentFormat(f *Format) *AdvancedQueries {
	q.Format = newPercentFormat()
	q.updateFormat(f)
	return q
}

func (q *AdvancedQueries) WithDataFormat(f *Format) *AdvancedQueries {
	q.Format = newDataFormat()
	q.updateFormat(f)
	return q
}

func (q *AdvancedQueries) WithDataRateFormat(f *Format) *AdvancedQueries {
	q.Format = newDataRateFormat()
	q.updateFormat(f)
	return q
}

func (q *AdvancedQueries) WithNumberFormat(f *Format) *AdvancedQueries {
	q.Format = newNumberFormat()
	q.updateFormat(f)
	return q
}

func (q *AdvancedQueries) WithNumberRateFormat(f *Format) *AdvancedQueries {
	q.Format = newNumberRateFormat()
	q.updateFormat(f)
	return q
}

func (q *AdvancedQueries) WithTimeFormat(f *Format) *AdvancedQueries {
	q.Format = newTimeFormat()
	q.updateFormat(f)
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
	MarkdownSource         *string              `json:"markdownSource,omitempty"`
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
	Base   NumberThresholdBase `json:"base"`
	Values []interface{}       `json:"values"`
}

type NumberThresholdBase struct {
	DisplayText string `json:"displayText"`
	Severity    string `json:"severity"`
}

type TeamSharingOptions struct {
	Type          string        `json:"type"`
	UserTeamsRole string        `json:"userTeamsRole"`
	SelectedTeams []interface{} `json:"selectedTeams"`
}

type SharingOptions struct {
	Member SharingMember `json:"member"`
	Role   string        `json:"role"`
}

type SharingMember struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type ScopeExpressionList struct {
	Operand     string      `json:"operand"`
	Operator    string      `json:"operator"`
	DisplayName string      `json:"displayName"`
	Value       []string    `json:"value"`
	Descriptor  interface{} `json:"descriptor"`
	IsVariable  bool        `json:"isVariable"`
}

type Dashboard struct {
	Version                 int                    `json:"version,omitempty"`
	CustomerID              interface{}            `json:"customerId"`
	TeamID                  int                    `json:"teamId"`
	Schema                  int                    `json:"schema"`
	AutoCreated             bool                   `json:"autoCreated"`
	PublicToken             string                 `json:"publicToken"`
	ScopeExpressionList     []*ScopeExpressionList `json:"scopeExpressionList"`
	Layout                  []*Layout              `json:"layout"`
	TeamScope               interface{}            `json:"teamScope"`
	EventDisplaySettings    EventDisplaySettings   `json:"eventDisplaySettings"`
	ID                      int                    `json:"id,omitempty"`
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	Username                string                 `json:"username"`
	Shared                  bool                   `json:"shared"`
	SharingSettings         []*SharingOptions      `json:"sharingSettings"`
	Public                  bool                   `json:"public"`
	Favorite                bool                   `json:"favorite"`
	CreatedOn               int64                  `json:"createdOn"`
	ModifiedOn              int64                  `json:"modifiedOn"`
	Panels                  []*Panels              `json:"panels"`
	TeamScopeExpressionList []interface{}          `json:"teamScopeExpressionList"`
	CreatedOnDate           string                 `json:"createdOnDate"`
	ModifiedOnDate          string                 `json:"modifiedOnDate"`
	TeamSharingOptions      TeamSharingOptions     `json:"teamSharingOptions"`
}

type dashboardWrapper struct {
	Dashboard *Dashboard `json:"dashboard"`
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
