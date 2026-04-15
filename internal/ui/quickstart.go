package ui

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/masaway/lazyprj/internal/config"
)

// QuickstartModel は任意ディレクトリからセッションを作成する画面
type QuickstartModel struct {
	width, height int
	cfg           *config.Config

	focus     int // 0=dirInput, 1=nameInput
	dirInput  textinput.Model
	nameInput textinput.Model

	status      string
	statusIsErr bool
	done        bool
	result      *config.Project // 確定後にAppが読む
}

type quickstartSavedMsg struct {
	project config.Project
	err     error
}

func NewQuickstart(cfg *config.Config) *QuickstartModel {
	dir := textinput.New()
	dir.Placeholder = "例: ~/work/myproject  または  /home/user/repos/myapp"
	dir.CharLimit = 256
	dir.Width = 52
	dir.Focus()

	name := textinput.New()
	name.Placeholder = "例: myproject"
	name.CharLimit = 64
	name.Width = 52

	return &QuickstartModel{
		dirInput:  dir,
		nameInput: name,
		cfg:       cfg,
	}
}

func (m *QuickstartModel) Init() tea.Cmd { return textinput.Blink }

func (m *QuickstartModel) IsDone() bool { return m.done }

func (m *QuickstartModel) Result() *config.Project { return m.result }

func (m *QuickstartModel) Resize(w, h int) {
	m.width = w
	m.height = h
	inputW := w - 20
	if inputW < 20 {
		inputW = 20
	}
	m.dirInput.Width = inputW
	m.nameInput.Width = inputW
}

func (m *QuickstartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case quickstartSavedMsg:
		if msg.err != nil {
			m.status = "保存エラー: " + msg.err.Error()
			m.statusIsErr = true
		} else {
			m.result = &msg.project
			m.done = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.result = nil
			m.done = true
			return m, nil

		case "tab":
			m = m.cycleFocus(true)
			// ディレクトリ入力からnameへ移動したとき、nameが空なら自動補完
			if m.focus == 1 && strings.TrimSpace(m.nameInput.Value()) == "" {
				dir := strings.TrimSpace(m.dirInput.Value())
				if dir != "" {
					expanded := config.ExpandPath(dir)
					base := filepath.Base(expanded)
					if base != "" && base != "." && base != "/" {
						m.nameInput.SetValue(base)
						m.nameInput.CursorEnd()
					}
				}
			}
			return m, nil

		case "shift+tab":
			m = m.cycleFocus(false)
			return m, nil

		case "enter":
			dir := strings.TrimSpace(m.dirInput.Value())
			if dir == "" {
				m.status = "ディレクトリを入力してください"
				m.statusIsErr = true
				return m, nil
			}
			// nameが空なら自動補完
			sessionName := strings.TrimSpace(m.nameInput.Value())
			if sessionName == "" {
				expanded := config.ExpandPath(dir)
				sessionName = filepath.Base(expanded)
			}
			if sessionName == "" || sessionName == "." || sessionName == "/" {
				m.status = "セッション名を入力してください"
				m.statusIsErr = true
				return m, nil
			}
			// 重複チェック
			for _, p := range m.cfg.Projects {
				if p.Name == sessionName {
					m.status = "'" + sessionName + "' は既に登録済みです"
					m.statusIsErr = true
					return m, nil
				}
			}
			proj := config.Project{Name: sessionName, Path: dir}
			cfg := m.cfg
			return m, func() tea.Msg {
				cfg.Projects = append(cfg.Projects, proj)
				if err := config.Save(cfg); err != nil {
					return quickstartSavedMsg{err: err}
				}
				return quickstartSavedMsg{project: proj}
			}
		}
	}

	var cmd tea.Cmd
	if m.focus == 0 {
		m.dirInput, cmd = m.dirInput.Update(msg)
	} else {
		m.nameInput, cmd = m.nameInput.Update(msg)
	}
	return m, cmd
}

func (m *QuickstartModel) cycleFocus(fwd bool) *QuickstartModel {
	if fwd {
		m.focus = (m.focus + 1) % 2
	} else {
		m.focus = (m.focus - 1 + 2) % 2
	}
	if m.focus == 0 {
		m.dirInput.Focus()
		m.nameInput.Blur()
	} else {
		m.nameInput.Focus()
		m.dirInput.Blur()
	}
	return m
}

func (m *QuickstartModel) View() string {
	dialogW := m.width - 12
	if dialogW > 80 {
		dialogW = 80
	}
	if dialogW < 54 {
		dialogW = 54
	}
	dialogH := 18

	innerW := dialogW - 4
	inputW := innerW - 4
	if inputW < 20 {
		inputW = 20
	}
	m.dirInput.Width = inputW
	m.nameInput.Width = inputW

	borderColorDir := colorBorder
	borderColorName := colorBorder
	if m.focus == 0 {
		borderColorDir = colorBorderFocus
	} else {
		borderColorName = colorBorderFocus
	}

	dirLabel := styleNormal.Render("ディレクトリ")
	nameLabel := styleNormal.Render("セッション名")
	note := styleDim.Render("チルダ（~/）が使えます。セッション名は省略するとディレクトリ名になります。")

	dirBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColorDir).
		Padding(0, 1).
		Width(m.dirInput.Width + 4).
		Render(m.dirInput.View())

	nameBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColorName).
		Padding(0, 1).
		Width(m.nameInput.Width + 4).
		Render(m.nameInput.View())

	var statusLine string
	if m.status != "" {
		if m.statusIsErr {
			statusLine = styleStatusErr.Width(innerW).Render(m.status)
		} else {
			statusLine = styleStatusOk.Width(innerW).Render(m.status)
		}
	} else {
		statusLine = styleStatusOk.Width(innerW).Render("")
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		note,
		"",
		dirLabel,
		dirBox,
		"",
		nameLabel,
		nameBox,
		"",
		statusLine,
		m.renderKeys(),
	)

	return panelBorderColored(body, dialogW, dialogH, 0, "新規セッション作成", colorYellow, colorYellow)
}

func (m *QuickstartModel) renderKeys() string {
	type hint struct{ key, desc string }
	hints := []hint{
		{"Enter", "作成"}, {"Esc", "キャンセル"},
	}
	var parts []string
	for _, h := range hints {
		parts = append(parts, styleKeyDesc.Render(h.desc)+styleKeySep.Render(": ")+styleKeyName.Render(h.key))
	}
	sep := styleKeySep.Render("  |  ")
	return lipgloss.NewStyle().Padding(0, 1).Render(strings.Join(parts, sep))
}
