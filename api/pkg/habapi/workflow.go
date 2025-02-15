package habapi

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/domain"
	"github.com/habiliai/habiliai/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/pkg/errors"
)

func (s *server) GetWorkflowStatus(ctx context.Context, req *ThreadId) (*WorkflowStatus, error) {
	_, thread, err := s.getThread(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var works []domain.ActionWork
	if err := helpers.GetTx(ctx).
		Preload("Action").
		Preload("Action.Agent").
		Find(&works, "thread_id = ?", thread.ID).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find action works")
	}

	return &WorkflowStatus{
		Works: gog.Map(works, func(w domain.ActionWork) *ActionWork {
			res := &ActionWork{
				Action: newActionPbFromDb(&w.Action),
				Done:   w.Done,
			}

			return res
		}),
	}, nil
}
