package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lazyprj/internal/config"
)

// SetupModel はスキャンディレクトリ設定画面
type SetupModel struct {
	width, height int
	input         textinput.Model
	cfg           *config.Config
	canSkip       bool // 既存の値があればEscでスキップ可
	done          bool
	status        string
	statusIsErr   bool
}

type setupSavedMsg struct{ err error }

func NewSetup(cfg *config.Config) *SetupModel {
	ti := textinput.New()
	ti.Placeholder = "例: ~/work  または  /home/user/projects"
	ti.CharLimit = 256
	ti.Width = 50

	canSkip := cfg.Settings.ScanDirectory != ""
	if canSkip {
		ti.SetValue(cfg.Settings.ScanDirectory)
	}
	ti.Focus()

	return &SetupModel{
		input:   ti,
		cfg:     cfg,
		canSkip: canSkip,
	}
}

func (m *SetupModel) Init() tea.Cmd { return textinput.Blink }

func (m *SetupModel) IsDone() bool { return m.done }

func (m *SetupModel) Resize(w, h int) {
	m.width = w
	m.height = h
	if w > 20 {
		m.input.Width = w - 20
	}
}

func (m *SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setupSavedMsg:
		if msg.err != nil {
			m.status = "保存エラー: " + msg.err.Error()
			m.statusIsErr = true
		} else {
			m.done = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := strings.TrimSpace(m.input.Value())
			if val == "" {
				m.status = "ディレクトリを入力してください"
				m.statusIsErr = true
				return m, nil
			}
			m.cfg.Settings.ScanDirectory = val
			cfg := m.cfg
			return m, func() tea.Msg {
				return setupSavedMsg{err: config.Save(cfg)}
			}

		case "esc":
			if m.canSkip {
				m.done = true
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *SetupModel) View() string {
	title := lipgloss.NewStyle().
		Background(colorBg2).
		Foreground(colorCyan).
		Bold(true).
		Padding(0, 1).
		Width(m.width).
		Render("スキャンディレクトリの設定")

	label := styleNormal.Render("プロジェクトをスキャンするディレクトリを入力してください。")
	note := styleDim.Render("チルダ（~/）が使えます。設定後にスキャン画面（n キー）で追加できます。")

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorderFocus).
		Padding(0, 1).
		Width(m.input.Width + 4).
		Render(m.input.View())

	var keys string
	if m.canSkip {
		keys = styleKeyDesc.Render("保存") + styleKeySep.Render(": ") + styleKeyName.Render("Enter") +
			styleKeySep.Render("  |  ") + styleKeyDesc.Render("キャンセル") + styleKeySep.Render(": ") + styleKeyName.Render("Esc")
	} else {
		keys = styleKeyDesc.Render("保存") + styleKeySep.Render(": ") + styleKeyName.Render("Enter")
	}

	var statusLine string
	if m.status != "" {
		if m.statusIsErr {
			statusLine = styleStatusErr.Width(m.width).Render(m.status)
		} else {
			statusLine = styleStatusOk.Width(m.width).Render(m.status)
		}
	} else {
		statusLine = styleStatusOk.Width(m.width).Render("")
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		label,
		"",
		note,
		"",
		inputBox,
	)

	panel := stylePanelFocused.
		Width(m.width - 2).
		Height(m.height - 4).
		Render(body)

	keysLine := lipgloss.NewStyle().Background(colorBg3).Width(m.width).Render(keys)

	return lipgloss.JoinVertical(lipgloss.Left, title, panel, statusLine, keysLine)
}
