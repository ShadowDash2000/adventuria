package reviews

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type repository interface {
	Create(ctx context.Context, review *model.Review) (*model.Review, error)
	Update(ctx context.Context, review *model.Review) (*model.Review, error)
	GetByActionAndPlayerID(ctx context.Context, actionId, playerId string) (*model.Review, error)
}

type reviewsStorage interface {
	IsValidFileURL(ctx context.Context, rawURL string) (bool, error)
	GetSize(ctx context.Context, id string) (int64, error)
	SaveFromBytes(ctx context.Context, bytes []byte, name string) (string, string, error)
	Delete(ctx context.Context, id string) error
}

type Reviews struct {
	domain         string
	repository     repository
	reviewsStorage reviewsStorage
}

const MaxTotalFilesSize = 5 * 1024 * 1024 // 5 mb

func NewReviews(domain string, repository repository, reviewsStorage reviewsStorage) *Reviews {
	return &Reviews{
		domain:         domain,
		repository:     repository,
		reviewsStorage: reviewsStorage,
	}
}

type CreateInput struct {
	Comment string
	Score   float64
}

func (r *Reviews) Create(ctx context.Context, input CreateInput) (*model.Review, error) {
	rawComment, fileIds, err := r.processComment(ctx, input.Comment, nil)
	if err != nil {
		return nil, err
	}

	review, err := model.NewReview(rawComment, input.Score)
	if err != nil {
		return nil, err
	}

	review.SetFiles(fileIds)

	return r.repository.Create(ctx, review)
}

type UpdateInput struct {
	Comment *string
	Score   *float64
}

func (r *Reviews) UpdateByActionAndPlayerID(ctx context.Context, actionId, playerId string, input UpdateInput) (*model.Review, error) {
	review, err := r.repository.GetByActionAndPlayerID(ctx, actionId, playerId)
	if err != nil {
		return nil, err
	}

	nothingToUpdate := true

	if input.Comment != nil {
		rawComment, fileIds, err := r.processComment(ctx, *input.Comment, review.Files())
		if err != nil {
			return nil, err
		}
		review.SetFiles(fileIds)
		comment, err := model.NewReviewComment(rawComment)
		if err != nil {
			return nil, err
		}
		review.SetComment(comment)
		nothingToUpdate = false
	}

	if input.Score != nil {
		score, err := model.NewReviewScore(*input.Score)
		if err != nil {
			return nil, err
		}
		review.SetScore(score)
		nothingToUpdate = false
	}

	if nothingToUpdate {
		return nil, errs.ErrNothingToUpdate
	}

	review, err = r.repository.Update(ctx, review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Reviews) processComment(ctx context.Context, comment string, oldIds []string) (string, []string, error) {
	nodes, err := html.ParseFragment(strings.NewReader(comment), &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	})
	if err != nil {
		return "", nil, err
	}

	var (
		totalBytes int64
		crawl      func(node *html.Node) error
		imgNodes   []struct {
			node      *html.Node
			attrIndex int
			attrValue string
		}
		existingIds []string
	)
	crawl = func(node *html.Node) error {
		if node.Type == html.ElementNode && node.Data == "img" {
			for i, attr := range node.Attr {
				if attr.Key != "src" {
					continue
				}

				if strings.HasPrefix(attr.Val, "data:image/") {
					totalBytes += int64(base64Len(attr.Val))
					imgNodes = append(imgNodes, struct {
						node      *html.Node
						attrIndex int
						attrValue string
					}{
						node:      node,
						attrIndex: i,
						attrValue: attr.Val,
					})
					break
				}

				isValidURL, err := r.reviewsStorage.IsValidFileURL(ctx, attr.Val)
				if err != nil {
					return err
				}

				if isValidURL {
					fileId, err := r.extractFileIDFromURL(attr.Val)
					if err != nil {
						return err
					}

					size, err := r.reviewsStorage.GetSize(ctx, fileId)
					if err != nil {
						return err
					}

					totalBytes += size
					existingIds = append(existingIds, fileId)
					break
				}
			}
		}

		for node := node.FirstChild; node != nil; node = node.NextSibling {
			err := crawl(node)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, node := range nodes {
		err = crawl(node)
		if err != nil {
			return "", nil, err
		}
	}

	if totalBytes > MaxTotalFilesSize {
		return "", nil, errs.ErrReviewCommentTooLong
	}

	resIds := helper.SlicesIntersection(existingIds, oldIds)
	for _, img := range imgNodes {
		rawBase64 := img.attrValue[strings.Index(img.attrValue, ",")+1:]
		b, err := base64.StdEncoding.DecodeString(rawBase64)
		if err != nil {
			return "", nil, err
		}

		id, path, err := r.reviewsStorage.SaveFromBytes(ctx, b, "some_name")
		if err != nil {
			return "", nil, err
		}

		resIds = append(resIds, id)
		img.node.Attr[img.attrIndex].Val = path
	}

	for _, id := range helper.SlicesDifference(oldIds, existingIds) {
		err = r.reviewsStorage.Delete(ctx, id)
		if err != nil {
			return "", nil, err
		}
	}

	var resBuf bytes.Buffer
	for _, node := range nodes {
		err = html.Render(&resBuf, node)
		if err != nil {
			return "", nil, err
		}
	}

	return resBuf.String(), resIds, nil
}

func (r *Reviews) extractFileIDFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	path := strings.FieldsFunc(u.Path, func(r rune) bool {
		return r == '/'
	})
	if len(path) != 5 {
		return "", errors.New("invalid file URL")
	}

	return path[3], nil
}

func base64Len(data string) int {
	rawBase64 := data[strings.Index(data, ",")+1:]

	padding := 0
	if strings.HasSuffix(rawBase64, "==") {
		padding = 2
	} else if strings.HasSuffix(rawBase64, "=") {
		padding = 1
	}

	return base64.StdEncoding.DecodedLen(len(rawBase64)) - padding
}
