package accural

import (
	"gophermart/internal/accural/client"
	"gophermart/internal/zlog"
	"time"
)

type AccuralController struct {
	client *client.AccuralClient
	done   chan struct{}
}

func StartNewController(addr string) *AccuralController {
	controller := &AccuralController{
		client: client.New(addr),
		done:   make(chan struct{}),
	}

	go func() {
		checkAccuralTicker := time.NewTicker(time.Second)

		select {
		case <-checkAccuralTicker.C:
			if err := controller.checkAccural(); err != nil {
				zlog.Logger.Errorf("check accural err=%s")
			}
		case <-controller.done:
			return
		}
	}()

	return controller
}

func (c *AccuralController) checkAccural() error {
	// get orders ids

	// check orders status

	return nil
}
