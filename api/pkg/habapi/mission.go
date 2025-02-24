package habapi

import (
	"context"
	"github.com/Masterminds/sprig/v3"
	"github.com/habiliai/alice/api/pkg/domain"
	"github.com/habiliai/alice/api/pkg/helpers"
	"github.com/mokiat/gog"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"strings"
	"text/template"
)

func (s *server) GetMissions(ctx context.Context, _ *emptypb.Empty) (*GetMissionsResponse, error) {
	tx := helpers.GetTx(ctx)

	stmt := tx.Model(&domain.Mission{})
	var numTotal int64
	if err := stmt.Count(&numTotal).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count missions")
	}

	var missions []domain.Mission
	if err := stmt.Preload(clause.Associations).
		Preload("Steps.Actions").
		Preload("Steps.Actions.Agent").
		Order("id ASC").
		Find(&missions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find missions")
	}

	resp := &GetMissionsResponse{
		Missions: make([]*Mission, len(missions)),
		NumTotal: int32(numTotal),
	}

	for i := range missions {
		resp.Missions[i] = newMissionPbFromDb(&missions[i])
	}

	return resp, nil
}

func (s *server) GetMission(ctx context.Context, req *MissionId) (*Mission, error) {
	tx := helpers.GetTx(ctx)

	var mission domain.Mission
	err := tx.
		Preload(clause.Associations).
		Preload("Steps.Actions").
		Preload("Steps.Actions.Agent").
		First(&mission, req.Id).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find mission by id %d", req.Id)
	}

	return newMissionPbFromDb(&mission), nil
}

func (s *server) GetMissionStepStatus(ctx context.Context, req *GetMissionStepStatusRequest) (*MissionStepStatus, error) {
	tx := helpers.GetTx(ctx)
	_, thread, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	var works []domain.ActionWork
	if err := helpers.GetTx(ctx).
		InnerJoins("Action").
		InnerJoins("Action.Step", tx.Where("\"Action__Step\".seq_no = ? AND \"Action__Step\".mission_id = ?", req.StepSeqNo, thread.MissionID)).
		Preload("Action.Agent").
		Find(
			&works,
			"thread_id = ?",
			thread.ID,
		).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find action works")
	}

	var step domain.Step
	if err := tx.
		First(&step, "mission_id = ? AND seq_no = ?", thread.MissionID, req.StepSeqNo).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find step")
	}

	return &MissionStepStatus{
		Step: newStepPbFromDb(&step),
		ActionWorks: gog.Map(works, func(w domain.ActionWork) *ActionWork {
			res := &ActionWork{
				Action: newActionPbFromDb(&w.Action),
				Done:   w.Done,
			}
			if w.Error != "" {
				res.Error = &w.Error
			}

			return res
		}),
	}, nil
}

func (s *server) generateResult(thread *domain.Thread) (string, error) {
	tpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(thread.Mission.ResultTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse result template")
	}

	memory := thread.Data.Data()

	logger.Debug("execute result template", "memory", memory)

	var resultBuilder strings.Builder
	if err := tpl.Execute(&resultBuilder, memory); err != nil {
		return "", errors.Wrapf(err, "failed to execute result template")
	}

	return resultBuilder.String(), nil
}
