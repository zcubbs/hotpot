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

	// Initialize inputs with default values
	repoURL := setupInput("https://github.com/user/repo", 60)
	repoURL.CharLimit = 150
	repoURL.Focus()

	branch := setupInput("main", 30)
	branch.CharLimit = 50
	branch.SetValue("main") // Default value

	authType := setupInput("token or ssh", 10)
	authType.CharLimit = 5

	authValue := setupInput("Enter token or path to SSH key", 60)
	authValue.CharLimit = 200

	localPath := setupInput("/path/to/local/recipe.yaml", 60)
	localPath.CharLimit = 150

	remotePath := setupInput("path/to/recipe.yaml", 60)
	remotePath.CharLimit = 150

	syncFreq := setupInput("5m", 20)
	syncFreq.CharLimit = 10
	syncFreq.SetValue("5m") // Default value

	// Create model with initialized config
	config := &syncd.Config{}

	model := &ConfigModel{
		repoURL:    repoURL,
		branch:     branch,
		authType:   authType,
		authValue:  authValue,
		localPath:  localPath,
		remotePath: remotePath,
		syncFreq:   syncFreq,
		config:     config,
	}

	// Try to load existing config
	existingConfig, err := syncd.LoadConfig()
	if err != nil && !strings.Contains(err.Error(), "file not found") {
		// Only set error if it's not a "file not found" error
		model.err = fmt.Errorf("failed to load existing config: %w", err)
	} else if existingConfig != nil {
		// Update the model's config with the existing one
		model.config = existingConfig

		// Prefill form fields with existing values
		if existingConfig.Repository.URL != "" {
			model.repoURL.SetValue(existingConfig.Repository.URL)
		}
		if existingConfig.Repository.Branch != "" {
			model.branch.SetValue(existingConfig.Repository.Branch)
		} else {
			// Set default branch if not set
			model.config.Repository.Branch = "main"
			model.branch.SetValue("main")
		}

		// Set auth type and value
		if existingConfig.Repository.Token != "" {
			model.authType.SetValue("token")
			model.authValue.SetValue(existingConfig.Repository.Token)
		} else if existingConfig.Repository.SSHKey != "" {
			model.authType.SetValue("ssh")
			model.authValue.SetValue(existingConfig.Repository.SSHKey)
		}

		if existingConfig.Sync.LocalPath != "" {
			model.localPath.SetValue(existingConfig.Sync.LocalPath)
		}
		if existingConfig.Sync.RemotePath != "" {
			model.remotePath.SetValue(existingConfig.Sync.RemotePath)
		}
		if existingConfig.Sync.Frequency != "" {
			model.syncFreq.SetValue(existingConfig.Sync.Frequency)
		} else {
			// Set default frequency if not set
			model.config.Sync.Frequency = "5m"
			model.syncFreq.SetValue("5m")
		}
	} else {
		// Set default values for new config
		model.config.Repository.Branch = "main"
		model.config.Sync.Frequency = "5m"
	}

	return model
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

	// Show error if any
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		b.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)) + "\n\n")
	}

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
	// Validate and set repository settings
	if m.repoURL.Value() == "" {
		return fmt.Errorf("repository URL is required")
	}
	m.config.Repository.URL = m.repoURL.Value()

	if m.branch.Value() == "" {
		return fmt.Errorf("branch is required")
	}
	m.config.Repository.Branch = m.branch.Value()

	// Validate and set auth type
	authType := m.authType.Value()
	if authType != "token" && authType != "ssh" {
		return fmt.Errorf("invalid auth type: must be 'token' or 'ssh'")
	}
	m.config.Repository.AuthType = authType

	// Validate and set auth value
	authValue := m.authValue.Value()
	if authValue == "" {
		return fmt.Errorf("auth value is required")
	}

	// Set token or SSH key based on auth type and clear the other
	if authType == "token" {
		m.config.Repository.Token = authValue
		m.config.Repository.SSHKey = "" // Clear SSH key when using token
	} else {
		m.config.Repository.SSHKey = authValue
		m.config.Repository.Token = "" // Clear token when using SSH
	}

	// Validate and set sync settings
	if m.localPath.Value() == "" {
		return fmt.Errorf("local path is required")
	}
	m.config.Sync.LocalPath = m.localPath.Value()

	if m.remotePath.Value() == "" {
		return fmt.Errorf("remote path is required")
	}
	m.config.Sync.RemotePath = m.remotePath.Value()

	if m.syncFreq.Value() == "" {
		return fmt.Errorf("sync frequency is required")
	}
	m.config.Sync.Frequency = m.syncFreq.Value()

	// Save the config
	return syncd.SaveConfig(m.config)
}
