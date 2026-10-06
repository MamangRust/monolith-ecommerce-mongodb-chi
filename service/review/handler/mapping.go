package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/repository"
	pbreviews "github.com/MamangRust/monolith-ecommerce-pb/reviews"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// ReviewResult aliases repository.ReviewResult so the proto mapping type switch
// sees the concrete type returned by the service layer.
type ReviewResult = repository.ReviewResult

func formatDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func formatDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05.000")
}

func (h *reviewHandleGrpc) mapResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *ReviewResult:
		return h.mapReviewResult(v)
	case *salemongo.Review:
		return h.mapReview(v)
	case []*ReviewResult:
		res := make([]*pbreviews.ReviewResponse, len(v))
		for i, r := range v {
			res[i] = h.mapReviewResultResponse(r)
		}
		return res
	case *salemongo.ReviewRow:
		return h.mapGetReviewByIDRow(v)
	default:
		return nil
	}
}

func (h *reviewHandleGrpc) mapReviewResult(v *ReviewResult) *pbreviews.ReviewResponseDeleteAt {
	return &pbreviews.ReviewResponseDeleteAt{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: strRev(v.CreatedAt),
		UpdatedAt: strRev(v.UpdatedAt),
		DeletedAt: &wrapperspb.StringValue{Value: strRev(v.DeletedAt)},
	}
}

func (h *reviewHandleGrpc) mapReviewResultResponse(v *ReviewResult) *pbreviews.ReviewResponse {
	return &pbreviews.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: strRev(v.CreatedAt),
		UpdatedAt: strRev(v.UpdatedAt),
	}
}

// mapReviewResultsDetail maps review rows for the product/merchant listing endpoints,
// embedding any attached review detail payloads.
func (h *reviewHandleGrpc) mapReviewResultsDetail(v []*ReviewResult) []*pbreviews.ReviewsDetailResponse {
	res := make([]*pbreviews.ReviewsDetailResponse, len(v))
	for i, r := range v {
		detail := &pbreviews.ReviewsDetailResponse{
			Id:        r.ReviewID,
			UserId:    r.UserID,
			ProductId: r.ProductID,
			Name:      r.Name,
			Comment:   r.Comment,
			Rating:    r.Rating,
			CreatedAt: strRev(r.CreatedAt),
			UpdatedAt: strRev(r.UpdatedAt),
		}
		if details, ok := r.ReviewDetails.(string); ok && details != "" {
			h.mapReviewDetails(details, detail)
		}
		res[i] = detail
	}
	return res
}

// mapReviewResultsDeleteAt maps review rows carrying soft-delete state for the
// active/trashed listing endpoints.
func (h *reviewHandleGrpc) mapReviewResultsDeleteAt(v []*ReviewResult) []*pbreviews.ReviewResponseDeleteAt {
	res := make([]*pbreviews.ReviewResponseDeleteAt, len(v))
	for i, r := range v {
		res[i] = &pbreviews.ReviewResponseDeleteAt{
			Id:        r.ReviewID,
			UserId:    r.UserID,
			ProductId: r.ProductID,
			Name:      r.Name,
			Comment:   r.Comment,
			Rating:    r.Rating,
			CreatedAt: strRev(r.CreatedAt),
			UpdatedAt: strRev(r.UpdatedAt),
			DeletedAt: &wrapperspb.StringValue{Value: strRev(r.DeletedAt)},
		}
	}
	return res
}

func (h *reviewHandleGrpc) mapReview(v *salemongo.Review) *pbreviews.ReviewResponse {
	return &pbreviews.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: formatDatePtr(v.CreatedAt),
		UpdatedAt: formatDatePtr(v.UpdatedAt),
	}
}

// mapReviewDeleteAt maps a review that carries soft-delete state (trash/restore).
func (h *reviewHandleGrpc) mapReviewDeleteAt(v *salemongo.Review) *pbreviews.ReviewResponseDeleteAt {
	return &pbreviews.ReviewResponseDeleteAt{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: formatDatePtr(v.CreatedAt),
		UpdatedAt: formatDatePtr(v.UpdatedAt),
		DeletedAt: &wrapperspb.StringValue{Value: formatDatePtr(v.DeletedAt)},
	}
}

func (h *reviewHandleGrpc) mapGetReviewByIDRow(v *salemongo.ReviewRow) *pbreviews.ReviewResponse {
	return &pbreviews.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: formatDate(&v.CreatedAt),
		UpdatedAt: formatDate(&v.UpdatedAt),
	}
}

func (h *reviewHandleGrpc) mapReviewDetails(reviewDetails string, res *pbreviews.ReviewsDetailResponse) {
	var details []struct {
		DetailID  int    `json:"detail_id"`
		Type      string `json:"type"`
		URL       string `json:"url"`
		Caption   string `json:"caption"`
		CreatedAt string `json:"created_at"`
	}

	err := json.Unmarshal([]byte(reviewDetails), &details)
	if err != nil {
		log.Printf("Error unmarshaling review details: %v", err)
		return
	}

	if len(details) > 0 {
		firstDetail := details[0]
		res.ReviewDetail = &pbreviews.ReviewDetailResponse{
			Id:        int32(firstDetail.DetailID),
			Type:      firstDetail.Type,
			Url:       firstDetail.URL,
			Caption:   firstDetail.Caption,
			CreatedAt: firstDetail.CreatedAt,
		}
	}
}

func strRev(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
