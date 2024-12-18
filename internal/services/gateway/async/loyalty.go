package async

import (
	"time"

	"github.com/RohanPoojary/gomq"
	"github.com/silazemli/lab3-template/internal/services/gateway/clients"
)

type Retry struct {
	Username string
	Time     time.Time
}

func LoyaltyDecrementRetry(broker gomq.Broker, loyaltyClient *clients.LoyaltyClient) {
	poller := broker.Subscribe(gomq.ExactMatcher("decrement loyalty counter"))
	go func() {
		for {
			value, ok := poller.Poll()
			if !ok {
				return
			}

			retry, ok := value.(Retry)
			if !ok {
				continue
			}

			for time.Since(retry.Time) <= 10*time.Second {
				continue
			}

			err := loyaltyClient.DecrementCounter(retry.Username)
			if err != nil {
				broker.Publish("decrement loyalty counter", Retry{Username: retry.Username})
			}

		}
	}()
}
