// package main
//
// import (
// 	"context"
// 	"fmt"
// 	"os"
//
// 	"github.com/ItsNotGoodName/ipcmanview/internal/db"
// 	"github.com/ItsNotGoodName/ipcmanview/migrations"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/rs/zerolog/log"
// )
//
// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }
//
// type model struct {
// 	textInput textinput.Model
// 	err       error
// }
//
// func initialModel() model {
// 	ti := textinput.New()
// 	ti.Placeholder = "Pikachu"
// 	ti.Focus()
// 	ti.CharLimit = 156
// 	ti.Width = 20
//
// 	return model{
// 		textInput: ti,
// 		err:       nil,
// 	}
// }
//
// func (m model) View() string {
// 	return fmt.Sprintf(
// 		"What’s your favorite Pokémon?\n\n%s\n\n%s",
// 		m.textInput.View(),
// 		"(esc to quit)",
// 	) + "\n"
// }
//
// func (m model) Init() tea.Cmd {
// 	return textinput.Blink
// }
//
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			return m, tea.Quit
// 		}
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
// 	}
//
// 	var cmd tea.Cmd
// 	m.textInput, cmd = m.textInput.Update(msg)
// 	return m, cmd
// }
//
// func initStack() tea.Msg {
// 	ctx, shutdown := context.WithCancel(interrupt.Context())
//
// 	// Database
// 	pool, err := db.New(ctx, os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		shutdown()
// 		return errMsg{fmt.Errorf("Failed to create database connection pool: %w", err)}
// 	}
//
// 	// Database migrate
// 	if err := migrations.Migrate(ctx, pool); err != nil {
// 		shutdown()
// 		pool.Close()
// 		return errMsg{fmt.Errorf("Failed to migrate database: %w", err)}
// 	}
//
// 	return app{
// 		ready: true,
// 		db:    pool,
// 		shutdown: func() {
// 			shutdown()
// 			pool.Close()
// 		},
// 	}
// }
//
// type app struct {
// 	ready    bool
// 	db       qes.Querier
// 	shutdown context.CancelFunc
// }
//
// type errMsg struct {
// 	err error
// }
//
// func (err errMsg) Error() string {
// 	return err.Error()
// }
//
// func init() {
// 	log.Logger = log.Output(nil)
// }
