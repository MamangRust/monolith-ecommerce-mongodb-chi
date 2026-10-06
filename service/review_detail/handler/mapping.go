package handler

import (
	"encoding/json"
	"log"

	pbreviews_detail "github.com/MamangRust/monolith-ecommerce-pb/reviews_detail"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (h *Handler) mapToReviewDetailResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *salemongo.ReviewDetailRow:
		return &pbreviews_detail.ReviewDetailsResponse{
			Id:        v.ReviewDetailID,
			ReviewId:  v.ReviewID,
			Type:      v.Type,
			Url:       v.Url,
			Caption:   derefStr(v.Caption),
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02"),
		}
	case *salemongo.ReviewDetailListRow:
		return &pbreviews_detail.ReviewDetailsResponse{
			Id:        v.ReviewDetailID,
			ReviewId:  v.ReviewID,
			Type:      v.Type,
			Url:       v.Url,
			Caption:   derefStr(v.Caption),
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02"),
		}
	case *salemongo.ReviewDetailActiveRow:
		var deletedAt string
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.Format("2006-01-02")
		}
		return &pbreviews_detail.ReviewDetailsResponseDeleteAt{
			Id:        v.ReviewDetailID,
			ReviewId:  v.ReviewID,
			Type:      v.Type,
			Url:       v.Url,
			Caption:   derefStr(v.Caption),
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	case *salemongo.ReviewDetail:
		var deletedAt string
		if v.DeletedAt != nil {
			deletedAt = v.DeletedAt.Format("2006-01-02")
		}
		return &pbreviews_detail.ReviewDetailsResponseDeleteAt{
			Id:        v.ReviewDetailID,
			ReviewId:  v.ReviewID,
			Type:      v.Type,
			Url:       v.Url,
			Caption:   derefStr(v.Caption),
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	default:
		log.Printf("Unknown type for mapping: %T", v)
		return nil
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
