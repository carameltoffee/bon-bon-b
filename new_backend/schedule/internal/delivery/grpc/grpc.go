package delivery

import (
	"context"
	schedule "schedule/internal/delivery/gen"
	"schedule/internal/models"
	"schedule/internal/usecase"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	logger *zap.Logger
	uc     *usecase.ScheduleUsecase
	schedule.UnimplementedScheduleServiceServer
}

func NewScheduleServerGrpc(gserver *grpc.Server, l *zap.Logger, uc *usecase.ScheduleUsecase) {
	server := &server{
		logger: l,
		uc:     uc,
	}
	schedule.RegisterScheduleServiceServer(gserver, server)
	reflection.Register(gserver)
}

func (s *server) CreateSlot(ctx context.Context, req *schedule.CreateSlotRequest) (*schedule.SlotResponse, error) {
	slot, err := s.uc.CreateSlot(ctx, req.UserId, req.Time)
	if err != nil {
		return nil, err
	}

	return &schedule.SlotResponse{Slot: convertToPB(slot)}, nil
}

func (s *server) GetSlot(ctx context.Context, req *schedule.SlotIdRequest) (*schedule.SlotResponse, error) {
	slot, err := s.uc.GetSlot(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &schedule.SlotResponse{Slot: convertToPB(slot)}, nil
}

func (s *server) UpdateSlot(ctx context.Context, req *schedule.UpdateSlotRequest) (*schedule.SlotResponse, error) {
	slot, err := s.uc.UpdateSlot(ctx, req.Id, req.IsAvailable)
	if err != nil {
		return nil, err
	}

	return &schedule.SlotResponse{Slot: convertToPB(slot)}, nil
}

func (s *server) DeleteSlot(ctx context.Context, req *schedule.SlotIdRequest) (*schedule.Empty, error) {
	if err := s.uc.DeleteSlot(ctx, req.Id); err != nil {
		return nil, err
	}

	return &schedule.Empty{}, nil
}

func (s *server) ListBusySlotsForUser(ctx context.Context, req *schedule.UserIdRequest) (*schedule.SlotListResponse, error) {
	slots, err := s.uc.ListBusySlotsForUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &schedule.SlotListResponse{Slots: convertToPBList(slots)}, nil
}

func (s *server) ListAvailableSlotsForUser(ctx context.Context, req *schedule.UserIdRequest) (*schedule.SlotListResponse, error) {
	slots, err := s.uc.ListAvailableSlotsForUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &schedule.SlotListResponse{Slots: convertToPBList(slots)}, nil
}

func (s *server) SaveSchedulePatternForUser(
	ctx context.Context,
	req *schedule.SaveSchedulePatternRequest,
) (*schedule.SchedulePatternResponse, error) {
	domainPattern := convertToDomainPattern(req.Pattern)
	saved, err := s.uc.SaveSchedulePatternForUser(ctx, domainPattern)
	if err != nil {
		return nil, err
	}

	return &schedule.SchedulePatternResponse{
		Pattern: convertToPBPattern(saved),
	}, nil
}

func (s *server) GetSchedulePatternsForUser(
	ctx context.Context,
	req *schedule.UserIdRequest,
) (*schedule.SchedulePatternResponse, error) {
	pattern, err := s.uc.GetSchedulePatternForUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &schedule.SchedulePatternResponse{
		Pattern: convertToPBPattern(pattern),
	}, nil
}

func convertToPBPattern(p *models.SchedulePattern) *schedule.SchedulePattern {
	slots := make([]*schedule.SlotPattern, len(p.Slots))
	for i, s := range p.Slots {
		slots[i] = &schedule.SlotPattern{
			Weekday: s.Weekday,
			Time:    s.Time,
		}
	}

	return &schedule.SchedulePattern{
		UserId:       p.UserId,
		ScheduleType: p.Type,
		TimeSlots:    slots,
		DaysAhead:    int32(p.DaysAhead),
	}
}

func convertToDomainPattern(p *schedule.SchedulePattern) *models.SchedulePattern {
	slots := make([]models.SlotPattern, len(p.TimeSlots))
	for i, s := range p.TimeSlots {
		slots[i] = models.SlotPattern{
			Weekday: s.Weekday,
			Time:    s.Time,
		}
	}

	return &models.SchedulePattern{
		UserId:    p.UserId,
		Type:      p.ScheduleType,
		Slots:     slots,
		DaysAhead: int(p.DaysAhead),
	}
}

func convertToPB(s *models.Slot) *schedule.Slot {
	return &schedule.Slot{
		Id:          s.ID,
		UserId:      s.UserID,
		Time:        s.Time.Format(time.RFC3339),
		IsAvailable: s.IsAvailable,
	}
}

func convertToPBList(slots []models.Slot) []*schedule.Slot {
	res := make([]*schedule.Slot, len(slots))
	for i := range slots {
		res[i] = convertToPB(&slots[i])
	}
	return res
}
