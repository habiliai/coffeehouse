package aliceapi

import (
	"context"
	"github.com/Masterminds/sprig/v3"
	"github.com/habiliai/alice/api/domain"
	"github.com/habiliai/alice/api/internal/db"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"strings"
	"text/template"
)

func (s *server) GetMissions(ctx context.Context, _ *emptypb.Empty) (*GetMissionsResponse, error) {
	ctx, tx := db.OpenSession(ctx, s.db)

	stmt := tx.Model(&domain.Mission{})
	var numTotal int64
	if err := stmt.Count(&numTotal).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count missions")
	}

	var missions []domain.Mission
	if err := stmt.Preload(clause.Associations).
		Preload("Steps.Actions").
		Order("id ASC").
		Find(&missions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find missions")
	}

	resp := &GetMissionsResponse{
		Missions: make([]*Mission, len(missions)),
		NumTotal: int32(numTotal),
	}

	var err error
	for i := range missions {
		resp.Missions[i], err = s.newMissionPbFromDb(ctx, &missions[i])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert mission %d", missions[i].ID)
		}
	}

	return resp, nil
}

func (s *server) GetMission(ctx context.Context, req *MissionId) (*Mission, error) {
	ctx, tx := db.OpenSession(ctx, s.db)

	var mission domain.Mission
	err := tx.
		Preload(clause.Associations).
		Preload("Steps.Actions").
		First(&mission, req.Id).Error
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find mission by id %d", req.Id)
	}

	return s.newMissionPbFromDb(ctx, &mission)
}

func (s *server) GetMissionStepStatus(ctx context.Context, req *GetMissionStepStatusRequest) (*MissionStepStatus, error) {
	ctx, tx := db.OpenSession(ctx, s.db)
	thread, err := s.getThread(ctx, req.ThreadId)
	if err != nil {
		return nil, err
	}

	var works []domain.ActionWork
	if err := tx.
		InnerJoins("Action").
		InnerJoins("Action.Step", tx.Where("\"Action__Step\".seq_no = ? AND \"Action__Step\".mission_id = ?", req.StepSeqNo, thread.MissionID)).
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

	res := MissionStepStatus{}
	for _, w := range works {
		a, err := s.newActionPb(ctx, &w.Action, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert action")
		}

		actionWork := &ActionWork{
			Action: a,
			Done:   w.Done,
		}
		if w.Error != "" {
			actionWork.Error = &w.Error
		}

		res.ActionWorks = append(res.ActionWorks, actionWork)
	}
	res.Step, err = s.newStepPbFromDb(ctx, &step)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert step")
	}

	return &res, nil
}

func (s *server) generateResult(ctx context.Context, thd *domain.Thread) (string, error) {
	tpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(thd.Mission.ResultTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse result template")
	}

	s.logger.Debug("execute result template", "memory", "")

	var resultBuilder strings.Builder
	if err := tpl.Execute(&resultBuilder, "test"); err != nil {
		return "", errors.Wrapf(err, "failed to execute result template")
	}

	return resultBuilder.String(), nil
}
