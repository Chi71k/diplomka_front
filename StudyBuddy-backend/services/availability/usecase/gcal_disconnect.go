package usecase

import "fmt"

type GCalDisconnectInput struct {
	UserID              string
	DeleteImportedSlots bool
}

type GCalDisconnect interface {
	Disconnect(in GCalDisconnectInput) error
}

type gcalDisconnect struct {
	gcalRepo GCalRepository
	slotRepo SlotRepository
}

func NewGCalDisconnect(gcalRepo GCalRepository, slotRepo SlotRepository) GCalDisconnect {
	return &gcalDisconnect{
		gcalRepo: gcalRepo,
		slotRepo: slotRepo,
	}
}

func (gc *gcalDisconnect) Disconnect(in GCalDisconnectInput) error {
	conn, err := gc.gcalRepo.GetConnection(in.UserID)
	if err != nil {
		return fmt.Errorf("get gcal connection: %w", err)
	}
	if conn == nil {
		return nil
	}

	if in.DeleteImportedSlots {
		if err := gc.slotRepo.DeleteAllForUser(in.UserID); err != nil {
			return fmt.Errorf("delete imported slots: %w", err)
		}
	}

	if err := gc.gcalRepo.DeleteConnection(in.UserID); err != nil {
		return fmt.Errorf("delete gcal connection: %w", err)
	}
	return nil
}
