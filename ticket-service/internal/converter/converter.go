package converter

import (
	ticketservicev1 "github.com/MaxFando/lms/ticket-service/api/grpc/gen/go/ticket-service/v1"
	"github.com/MaxFando/lms/ticket-service/internal/entity"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func ToTicketServiceFromEntity(t *entity.Ticket) *ticketservicev1.Ticket {
	var userID *wrapperspb.Int32Value
	if t.UserID != nil {
		userID = &wrapperspb.Int32Value{Value: *t.UserID}
	}
	return &ticketservicev1.Ticket{
		TicketId:  t.ID,
		UserId:    userID,
		DrawId:    t.DrawID,
		Numbers:   t.Numbers,
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}
}

func ToTicketWithDrawServiceFromEntity(t *entity.TicketWithDraw) *ticketservicev1.TicketWithDraw {
	return &ticketservicev1.TicketWithDraw{
		TicketId:  t.ID,
		UserId:    *t.UserID,
		DrawId:    t.DrawID,
		Numbers:   t.Numbers,
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		Draw:      ToDrawServiceFromEntity(&t.Draw),
	}
}

func ToDrawServiceFromEntity(d *entity.Draw) *ticketservicev1.Draw {
	return &ticketservicev1.Draw{
		DrawId:      d.ID,
		LotteryType: d.LotteryType,
		Status:      d.Status,
		StartTime:   d.StartTime.Format(time.RFC3339),
		EndTime:     d.EndTime.Format(time.RFC3339),
	}
}
