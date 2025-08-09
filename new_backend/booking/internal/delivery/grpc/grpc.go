package delivery

import (
	booking "booking/internal/delivery/gen"
	"booking/internal/models"
	"booking/internal/usecase"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const timeLayout = time.RFC3339

type server struct {
	uc *usecase.BookingUsecase
	l  *zap.Logger
	booking.UnimplementedBookingServiceServer
}

func NewBookingServerGrpc(gserver *grpc.Server, l *zap.Logger, uc *usecase.BookingUsecase) {
	server := &server{
		l:  l,
		uc: uc,
	}
	booking.RegisterBookingServiceServer(gserver, server)
	reflection.Register(gserver)
}

func (s *server) CreateBooking(ctx context.Context, req *booking.CreateBookingRequest) (*booking.BookingResponse, error) {
	book := &models.Booking{
		UserID:       req.GetUserId(),
		SlotID:       req.GetSlotId(),
		Type:         req.GetType(),
		ContactPhone: req.GetContactPhone(),
		Comment:      req.GetComment(),
	}

	created, err := s.uc.CreateBooking(ctx, book)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create booking failed: %v", err)
	}

	return &booking.BookingResponse{Booking: bookingModelToProto(created)}, nil
}

func (s *server) GetBooking(ctx context.Context, req *booking.BookingIdRequest) (*booking.BookingResponse, error) {
	book, err := s.uc.GetBooking(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, usecase.ErrBookingNotFound) {
			return nil, status.Errorf(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "get booking failed: %v", err)
	}

	return &booking.BookingResponse{Booking: bookingModelToProto(book)}, nil
}

func (s *server) UpdateBooking(ctx context.Context, req *booking.UpdateBookingRequest) (*booking.BookingResponse, error) {
	updated, err := s.uc.UpdateBooking(ctx, req.GetId(), req.GetStatus(), req.GetComment())
	if err != nil {
		if errors.Is(err, usecase.ErrBookingNotFound) {
			return nil, status.Errorf(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "update booking failed: %v", err)
	}

	return &booking.BookingResponse{Booking: bookingModelToProto(updated)}, nil
}

func (s *server) DeleteBooking(ctx context.Context, req *booking.BookingIdRequest) (*booking.Empty, error) {
	err := s.uc.DeleteBooking(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, usecase.ErrBookingNotFound) {
			return nil, status.Errorf(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "delete booking failed: %v", err)
	}

	return &booking.Empty{}, nil
}
func (s *server) ListBookings(ctx context.Context, req *booking.UserIdRequest) (*booking.BookingListResponse, error) {
	bookings, err := s.uc.ListBookingsByUser(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list bookings failed: %v", err)
	}

	resp := &booking.BookingListResponse{}
	for _, b := range bookings {
		resp.Bookings = append(resp.Bookings, bookingModelToProto(&b))
	}

	return resp, nil
}

func bookingModelToProto(m *models.Booking) *booking.Booking {
	var createdAt, updatedAt string
	if !m.CreatedAt.IsZero() {
		createdAt = m.CreatedAt.UTC().Format(timeLayout)
	}
	if !m.UpdatedAt.IsZero() {
		updatedAt = m.UpdatedAt.UTC().Format(timeLayout)
	}

	return &booking.Booking{
		Id:           m.ID,
		UserId:       m.UserID,
		SlotId:       m.SlotID,
		Type:         m.Type,
		ContactPhone: m.ContactPhone,
		Comment:      m.Comment,
		Status:       m.Status,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
