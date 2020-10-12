package telegram

import (
	"context"
	"sync"

	"github.com/cordalace/wireguard-for-homies/internal/manager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type Telegram struct {
	db                   DB
	token                string
	bot                  *tgbotapi.BotAPI
	wgManager            *manager.Manager
	fsm                  fsm
	messageTypeByCommand map[string]messageType
	ctx                  context.Context
	cancelCtx            context.CancelFunc
	wg                   *sync.WaitGroup
	logger               *zap.Logger
}

func NewTelegram(db DB, token string, wgManager *manager.Manager, logger *zap.Logger) *Telegram {
	ctx, cancel := context.WithCancel(context.Background())
	return &Telegram{
		db:        db,
		token:     token,
		wgManager: wgManager,
		ctx:       ctx,
		cancelCtx: cancel,
		wg:        &sync.WaitGroup{},
		logger:    logger,
	}
}

func (t *Telegram) Init() error {
	var err error
	t.bot, err = tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return err
	}

	t.ctx = context.Background()

	t.fsm = fsm{
		fsmKey{ChatStateInitial, messageTypeCmdUnknown}:            t.commandUnrecognized,
		fsmKey{ChatStateInitial, messageTypeText}:                  t.showHelp,
		fsmKey{ChatStateInitial, messageTypeCmdCancel}:             t.withReadWriteTx(t.commandCancel),
		fsmKey{ChatStateInitial, messageTypeCmdHelp}:               t.showHelp,
		fsmKey{ChatStateInitial, messageTypeCmdNewSubnet}:          t.withReadWriteTx(t.commandNewSubnet),
		fsmKey{ChatStateSubnetExpectCIDR, messageTypeCmdUnknown}:   t.commandUnrecognized,
		fsmKey{ChatStateSubnetExpectCIDR, messageTypeText}:         t.withReadWriteTx(t.provideSubnetCIDR),
		fsmKey{ChatStateSubnetExpectCIDR, messageTypeCmdCancel}:    t.withReadWriteTx(t.commandCancel),
		fsmKey{ChatStateSubnetExpectCIDR, messageTypeCmdHelp}:      t.showHelp,
		fsmKey{ChatStateSubnetExpectCIDR, messageTypeCmdNewSubnet}: t.withReadWriteTx(t.commandNewSubnet),
	}

	t.messageTypeByCommand = map[string]messageType{
		"cancel":    messageTypeCmdCancel,
		"help":      messageTypeCmdHelp,
		"newsubnet": messageTypeCmdNewSubnet,
	}

	return t.setMyCommands()
}

func (t *Telegram) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for {
		select {
		case update := <-updates:
			t.logger.Info("received", zap.String("text", update.Message.Text))
			t.wg.Add(1)
			go t.updateGoroutine(&update)
		case <-t.ctx.Done():
			t.logger.Info("done")
			t.bot.StopReceivingUpdates()
			return nil
		}
		t.logger.Info("loop")
	}
}

func (t *Telegram) updateGoroutine(update *tgbotapi.Update) {
	defer t.wg.Done()
	if err := t.handleUpdate(update); err != nil {
		t.logger.Error("error handling update", zap.Error(err))
	}
}

// Close stops accepting new messages and waits handlers to finish.
func (t *Telegram) Close() {
	t.cancelCtx()
	t.wg.Wait()
}
