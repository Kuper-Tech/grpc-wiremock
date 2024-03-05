package wearable_service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nktch1/wearable/pkg/clients/push_sender"
	"github.com/nktch1/wearable/pkg/server/wearable"
)

const expectedStatusCode = 1126

func (p *Service) BeatsPerMinute(in *wearable.BeatsPerMinuteRequest, stream wearable.WearableService_BeatsPerMinuteServer) error {
	const batchSize = 30

	for idx := 0; idx < batchSize; idx++ {
		heartRate := newRandInt()

		fmt.Printf("current heart rate: %d\n", heartRate)

		if somethingIsGoingWrong(heartRate) {
			fmt.Printf("\t something is going wrong, notify\n")

			notifyResponse, err := p.sender.Notify(context.Background(), &push_sender.NotifyRequest{
				Uuid:    "some_uuid",
				Message: "Something is going wrong!",
			})

			if err != nil {
				return fmt.Errorf("notify: %w", err)
			}

			fmt.Printf("\t notifying status is %d\n", notifyResponse.Status)

			if notifyResponse.Status != expectedStatusCode {
				return fmt.Errorf("status code %d is incorrect, expected: %d", notifyResponse.Status, expectedStatusCode)
			}
		}

		response := wearable.BeatsPerMinuteResponse{
			Value:  newRandInt(),
			Minute: uint32(idx),
		}

		if err := stream.Send(&response); err != nil {
			return fmt.Errorf("send: %w", err)
		}

		time.Sleep(time.Second)
	}

	fmt.Println("batch processing is done!")

	return nil
}

func newRandInt() uint32 {
	min := 30
	max := 160
	return uint32(rand.Intn(max-min) + min)
}

func somethingIsGoingWrong(heartRate uint32) bool {
	return heartRate < 60 || heartRate > 140
}
