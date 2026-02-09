package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

// UnitPrefs holds user-selected unit preferences for the dashboard.
type UnitPrefs struct {
	TempF    bool
	WindMph  bool
	PrecipIn bool
	DistMi   bool
}

// DashboardModel is the Bubbletea model for the live weather dashboard.
type DashboardModel struct {
	// Station metadata
	stationName string
	timezone    string
	deviceID    int
	apiToken    string

	// Live data
	currentObs  *ObsData
	rapidWind   *RapidWindData
	windHistory []RapidWindData
	events      []EventData

	// Connection state
	connected         bool
	lastUpdate        time.Time
	errMsg            string
	reconnecting      bool
	reconnectAttempts int

	// Display
	width  int
	height int

	// Preferences
	unitPrefs UnitPrefs

	// WS internals
	wsConn *websocket.Conn
	wsDone chan struct{}
	msgCh  chan tea.Msg
}

// --- Tea message types ---

type wsObsMsg struct{ obs ObsData }
type wsRapidWindMsg struct{ wind RapidWindData }
type wsEventMsg struct{ event EventData }
type wsConnectedMsg struct{ conn *websocket.Conn }
type wsErrorMsg struct{ err error }
type wsReconnectMsg struct{}
type tickMsg time.Time

// Init starts the WS connection and background tickers.
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.connectWSCmd(),
		m.waitForWSMsg(),
		tickEvery(),
	)
}

// Update handles all incoming messages.
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.wsConn != nil {
				m.wsConn.Close()
			}
			if m.wsDone != nil {
				select {
				case <-m.wsDone:
				default:
					close(m.wsDone)
				}
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case wsConnectedMsg:
		m.wsConn = msg.conn
		m.connected = true
		m.errMsg = ""
		m.reconnecting = false
		m.reconnectAttempts = 0
		return m, nil

	case wsObsMsg:
		m.currentObs = &msg.obs
		m.lastUpdate = time.Now()
		m.errMsg = ""
		return m, m.waitForWSMsg()

	case wsRapidWindMsg:
		m.rapidWind = &msg.wind
		m.windHistory = append(m.windHistory, msg.wind)
		if len(m.windHistory) > 20 {
			m.windHistory = m.windHistory[len(m.windHistory)-20:]
		}
		m.lastUpdate = time.Now()
		return m, m.waitForWSMsg()

	case wsEventMsg:
		m.events = append([]EventData{msg.event}, m.events...)
		if len(m.events) > 10 {
			m.events = m.events[:10]
		}
		m.lastUpdate = time.Now()
		return m, m.waitForWSMsg()

	case wsErrorMsg:
		m.connected = false
		m.errMsg = msg.err.Error()
		if m.wsConn != nil {
			m.wsConn.Close()
			m.wsConn = nil
		}
		if m.reconnectAttempts < 5 {
			m.reconnecting = true
			m.reconnectAttempts++
			return m, tea.Batch(
				m.waitForWSMsg(),
				reconnectAfter(5*time.Second),
			)
		}
		m.reconnecting = false
		return m, m.waitForWSMsg()

	case wsReconnectMsg:
		if !m.connected && m.reconnecting {
			return m, tea.Batch(
				m.connectWSCmd(),
				m.waitForWSMsg(),
			)
		}
		return m, nil

	case tickMsg:
		return m, tickEvery()
	}

	return m, nil
}

// View delegates to the dashboard renderer.
func (m DashboardModel) View() string {
	return renderDashboard(m)
}

// --- WS command functions ---

func (m DashboardModel) connectWSCmd() tea.Cmd {
	return func() tea.Msg {
		wsURL := fmt.Sprintf("wss://ws.weatherflow.com/swd/data?token=%s", m.apiToken)
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			return wsErrorMsg{err: fmt.Errorf("ws dial: %w", err)}
		}

		// Send listen_start
		listenStart := WSListenStart{
			Type:     "listen_start",
			DeviceID: m.deviceID,
			ID:       "tempest-cli-obs",
		}
		if err := conn.WriteJSON(listenStart); err != nil {
			conn.Close()
			return wsErrorMsg{err: fmt.Errorf("ws listen_start: %w", err)}
		}

		// Send listen_rapid_start
		listenRapid := WSListenRapidStart{
			Type:     "listen_rapid_start",
			DeviceID: m.deviceID,
			ID:       "tempest-cli-rapid",
		}
		if err := conn.WriteJSON(listenRapid); err != nil {
			conn.Close()
			return wsErrorMsg{err: fmt.Errorf("ws listen_rapid_start: %w", err)}
		}

		// Start the read loop goroutine
		go wsReadLoop(conn, m.msgCh, m.wsDone)

		return wsConnectedMsg{conn: conn}
	}
}

func wsReadLoop(conn *websocket.Conn, msgCh chan<- tea.Msg, done chan struct{}) {
	defer func() {
		select {
		case <-done:
		default:
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
		}

		// Set read deadline to detect stale connections (11 min > obs_st ~60s interval)
		conn.SetReadDeadline(time.Now().Add(11 * time.Minute))

		_, raw, err := conn.ReadMessage()
		if err != nil {
			select {
			case <-done:
				return
			default:
				msgCh <- wsErrorMsg{err: fmt.Errorf("ws read: %w", err)}
				return
			}
		}

		var envelope WSMessage
		if err := json.Unmarshal(raw, &envelope); err != nil {
			continue
		}

		switch envelope.Type {
		case "obs_st":
			var obs WSObservation
			if err := json.Unmarshal(raw, &obs); err != nil {
				continue
			}
			if len(obs.Obs) > 0 {
				parsed := ParseObsArray(obs.Obs[0])
				msgCh <- wsObsMsg{obs: parsed}
			}

		case "rapid_wind":
			var rw WSRapidWind
			if err := json.Unmarshal(raw, &rw); err != nil {
				continue
			}
			if len(rw.Ob) >= 3 {
				parsed := ParseRapidWind(rw.Ob)
				msgCh <- wsRapidWindMsg{wind: parsed}
			}

		case "evt_precip":
			var ep WSEventPrecip
			if err := json.Unmarshal(raw, &ep); err != nil {
				continue
			}
			ts := time.Now()
			if len(ep.Evt) > 0 {
				ts = time.Unix(int64(ep.Evt[0]), 0)
			}
			msgCh <- wsEventMsg{event: EventData{
				Timestamp: ts,
				Type:      "rain",
				Detail:    "Rain started",
			}}

		case "evt_strike":
			var es WSEventStrike
			if err := json.Unmarshal(raw, &es); err != nil {
				continue
			}
			ts := time.Now()
			detail := "Lightning detected"
			if len(es.Evt) >= 3 {
				ts = time.Unix(int64(es.Evt[0]), 0)
				dist := es.Evt[1]
				detail = fmt.Sprintf("Lightning %.0fkm away", dist)
			}
			msgCh <- wsEventMsg{event: EventData{
				Timestamp: ts,
				Type:      "lightning",
				Detail:    detail,
			}}

		case "ack":
			// Acknowledged, nothing to do
		}
	}
}

func (m DashboardModel) waitForWSMsg() tea.Cmd {
	return func() tea.Msg {
		return <-m.msgCh
	}
}

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func reconnectAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return wsReconnectMsg{}
	})
}
