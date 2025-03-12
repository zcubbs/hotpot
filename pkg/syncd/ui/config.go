package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zcubbs/hotpot/pkg/syncd"
	"strings"
)

var (
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Border(lipgloss.NormalBorder()).Padding(0, 3).Render("Submit")
	blurredButton = blurredStyle.Border(lipgloss.NormalBorder()).Padding(0, 3).Render("Submit")
)

type ConfigModel struct {
	repoURL    textinput.Model
	branch     textinput.Model
	authType   textinput.Model
	authValue  textinput.Model
	localPath  textinput.Model
	remotePath textinput.Model
	syncFreq   textinput.Model
	focused    int
	done       bool
	err        error
	config     *syncd.Config
}

// Done returns true if configuration is complete
func (m *ConfigModel) Done() bool {
	return m.done
}

func InitialConfigModel() *ConfigModel {
	// Initialize config
	config := &syncd.Config{}
	config.Repository.Branch = "main" // Set default branch

	// Common styles
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Helper function to set up input
	setupInput := func(placeholder string, width int) textinput.Model {
		input := textinput.New()
		input.Placeholder = placeholder
		input.Width = width
		input.PromptStyle = promptStyle
		input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		return input
	}

	// Repository URL
	repoURL := setupInput("https://github.com/user/repo", 60)
	repoURL.CharLimit = 150
	repoURL.Focus()

	// Branch
	branch := setupInput("main", 30)
	branch.CharLimit = 50
	branch.SetValue("main") // Set default value

	// Auth Type
	authType := setupInput("token or ssh", 10)
	authType.CharLimit = 5

	// Token/SSH Key
	authValue := setupInput("Enter token or path to SSH key", 60)
	authValue.CharLimit = 200

	// Local Path
	localPath := setupInput("/path/to/local/recipe.yaml", 60)
	localPath.CharLimit = 150

	// Remote Path
	remotePath := setupInput("path/to/recipe.yaml", 60)
	remotePath.CharLimit = 150

	// Sync Frequency
	syncFreq := setupInput("5m", 20)
	syncFreq.CharLimit = 10
	syncFreq.SetValue("5m") // Set default value

	return &ConfigModel{
		repoURL:    repoURL,
		branch:     branch,
		authType:   authType,
		authValue:  authValue,
		localPath:  localPath,
		remotePath: remotePath,
		syncFreq:   syncFreq,
		config:     config,
	}
}

func (m *ConfigModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == 7 {
				if err := m.validateAndSave(); err != nil {
					m.err = err
					return m, nil
				}
				m.done = true
				return m, tea.Quit
			}

			// Update focus
			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > 7 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = 7
			}

			// Update input focus
			m.repoURL.Blur()
			m.branch.Blur()
			m.authType.Blur()
			m.authValue.Blur()
			m.localPath.Blur()
			m.remotePath.Blur()
			m.syncFreq.Blur()

			switch m.focused {
			case 0:
				m.repoURL.Focus()
			case 1:
				m.branch.Focus()
			case 2:
				m.authType.Focus()
			case 3:
				m.authValue.Focus()
			case 4:
				m.localPath.Focus()
			case 5:
				m.remotePath.Focus()
			case 6:
				m.syncFreq.Focus()
			}

			return m, nil
		}
	}

	// Handle character input
	var cmd tea.Cmd
	switch m.focused {
	case 0:
		m.repoURL, cmd = m.repoURL.Update(msg)
	case 1:
		m.branch, cmd = m.branch.Update(msg)
	case 2:
		m.authType, cmd = m.authType.Update(msg)
	case 3:
		m.authValue, cmd = m.authValue.Update(msg)
	case 4:
		m.localPath, cmd = m.localPath.Update(msg)
	case 5:
		m.remotePath, cmd = m.remotePath.Update(msg)
	case 6:
		m.syncFreq, cmd = m.syncFreq.Update(msg)
	}

	return m, cmd
}

func (m *ConfigModel) View() string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	b.WriteString(titleStyle.Render("🔧 Configure Hotpot Sync Daemon") + "\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	b.WriteString(helpStyle.Render("Use Tab/Shift+Tab or Up/Down arrows to navigate. Press Enter to submit.") + "\n\n")

	// Debug info
	b.WriteString(fmt.Sprintf("Focused: %d\n\n", m.focused))

	// Helper function to render a field
	renderField := func(label string, input textinput.Model, isFocused bool) string {
		// Label style with fixed width for alignment
		labelStyle := lipgloss.NewStyle().Bold(true)
		if isFocused {
			labelStyle = labelStyle.Foreground(lipgloss.Color("205"))
		}

		// Input style with consistent width and padding
		inputStyle := lipgloss.NewStyle().Width(50)
		if isFocused {
			inputStyle = inputStyle.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("205")).Padding(0, 1)
		} else {
			inputStyle = inputStyle.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
		}

		// Create a container style for the input box with left margin
		containerStyle := lipgloss.NewStyle().MarginLeft(2)

		// Right-align the label by padding with spaces
		padding := strings.Repeat(" ", 16-len(label))
		return padding + labelStyle.Render(label) + ":\n" + containerStyle.Render(inputStyle.Render(input.View()))
	}

	// Repository section
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	b.WriteString(sectionStyle.Render("Repository Settings") + "\n")
	b.WriteString(renderField("Repository URL", m.repoURL, m.focused == 0) + "\n")
	b.WriteString(renderField("Branch", m.branch, m.focused == 1) + "\n")
	b.WriteString(renderField("Auth Type", m.authType, m.focused == 2) + "\n")
	b.WriteString(renderField("Token/SSH Key", m.authValue, m.focused == 3) + "\n")

	// Sync section
	b.WriteString("\n" + sectionStyle.Render("Sync Settings") + "\n")
	b.WriteString(renderField("Local Path", m.localPath, m.focused == 4) + "\n")
	b.WriteString(renderField("Remote Path", m.remotePath, m.focused == 5) + "\n")
	b.WriteString(renderField("Sync Frequency", m.syncFreq, m.focused == 6) + "\n")

	// Button
	button := &blurredButton
	if m.focused == 7 {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n", *button)

	// Error
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		b.WriteString("\n" + errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n")
	}

	return b.String()
}

func (m *ConfigModel) validateAndSave() error {
	// Set repository settings
	m.config.Repository.URL = m.repoURL.Value()
	m.config.Repository.Branch = m.branch.Value()

	// Set auth type
	authType := m.authType.Value()
	if authType != "token" && authType != "ssh" {
		return fmt.Errorf("invalid auth type: must be 'token' or 'ssh'")
	}
	m.config.Repository.AuthType = authType

	// Set token or SSH key based on auth type
	if authType == "token" {
		m.config.Repository.Token = m.authValue.Value()
	} else {
		m.config.Repository.SSHKey = m.authValue.Value()
	}

	// Set sync settings
	m.config.Sync.LocalPath = m.localPath.Value()
	m.config.Sync.RemotePath = m.remotePath.Value()
	m.config.Sync.Frequency = m.syncFreq.Value()

	// Validate required fields
	if m.config.Repository.URL == "" {
		return fmt.Errorf("repository URL is required")
	}
	if m.config.Sync.LocalPath == "" {
		return fmt.Errorf("local path is required")
	}
	if m.config.Sync.Frequency == "" {
		return fmt.Errorf("sync frequency is required")
	}

	return syncd.SaveConfig(m.config)
}
