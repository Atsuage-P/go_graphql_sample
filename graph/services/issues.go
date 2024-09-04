package services

import (
	"context"
	"fmt"
	"log"
	"mygql/graph/db"
	"mygql/graph/model"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type issueService struct {
	exec boil.ContextExecutor
}

func convertIssue(issue *db.Issue) *model.Issue {
	issueURL, err := model.UnmarshalURI(issue.URL)
	if err != nil {
		log.Println("invalid URI", issue.URL)
	}

	return &model.Issue{
		ID:         issue.ID,
		URL:        issueURL,
		Title:      issue.Title,
		Closed:     (issue.Closed == 1),
		Number:     issue.Number,
		Author:     &model.User{ID: issue.Author},
		Repository: &model.Repository{},
	}
}

func convertIssueConnection(issues db.IssueSlice, hasPrevPage, hasNextPage bool) *model.IssueConnection {
	var result model.IssueConnection

	for _, dbi := range issues {
		issue := convertIssue(dbi)

		result.Edges = append(result.Edges, &model.IssueEdge{Cursor: issue.ID, Node: issue})
		result.Nodes = append(result.Nodes, issue)
	}
	result.TotalCount = len(issues)

	result.PageInfo = &model.PageInfo{}
	if result.TotalCount != 0 {
		result.PageInfo.StartCursor = &result.Nodes[0].ID
		result.PageInfo.EndCursor = &result.Nodes[result.TotalCount-1].ID
	}
	result.PageInfo.HasPreviousPage = hasPrevPage
	result.PageInfo.HasNextPage = hasNextPage

	return &result
}

func (i *issueService) GetIssueByID(ctx context.Context, id string) (*model.Issue, error) {
	issue, err := db.FindIssue(ctx, i.exec, id,
		db.IssueColumns.ID,
		db.IssueColumns.URL,
		db.IssueColumns.Title,
		db.IssueColumns.Closed,
		db.IssueColumns.Number,
		db.IssueColumns.Repository,
	)
	if err != nil {
		return nil, err
	}

	return convertIssue(issue), nil
}

func (i *issueService) GetIssueByRepoAndNumber(ctx context.Context, repoID string, number int) (*model.Issue, error) {
	issue, err := db.Issues(
		qm.Select(
			db.IssueColumns.ID,
			db.IssueColumns.URL,
			db.IssueColumns.Title,
			db.IssueColumns.Closed,
			db.IssueColumns.Number,
			db.IssueColumns.Author,
			db.IssueColumns.Repository,
		),
		db.IssueWhere.Repository.EQ(repoID),
		db.IssueWhere.Number.EQ(number),
	).One(ctx, i.exec)
	if err != nil {
		return nil, err
	}

	return convertIssue(issue), nil
}

func (i *issueService) ListIssueInRepository(ctx context.Context, repoID string, after *string, before *string, first *int, last *int) (*model.IssueConnection, error) {
	cond := []qm.QueryMod{
		qm.Select(
			db.IssueColumns.ID,
			db.IssueColumns.URL,
			db.IssueColumns.Title,
			db.IssueColumns.Closed,
			db.IssueColumns.Number,
			db.IssueColumns.Author,
			db.IssueColumns.Repository,
		),
		db.IssueWhere.Repository.EQ(repoID),
	}
	var scanDesc bool

	// ページネーション
	switch {
	case (after != nil) && (before != nil):
		// before < ID < after の範囲
		cond = append(cond,
			db.IssueWhere.ID.GT(*after),
			db.IssueWhere.ID.LT(*before),
		)
	case after != nil:
		// ID > after の範囲で first 件数分を昇順で取得
		cond = append(cond,
			db.IssueWhere.ID.GT(*after),
			qm.OrderBy(fmt.Sprintf("%s asc", db.IssueColumns.ID)),
		)
		if first != nil {
			cond = append(cond, qm.Limit(*first))
		}
	case before != nil:
		// ID < before の範囲で last 件数分を降順で取得
		scanDesc = true
		cond = append(cond,
			db.IssueWhere.ID.LT(*before),
			qm.OrderBy(fmt.Sprintf("%s desc", db.IssueColumns.ID)),
		)
		if last != nil {
			cond = append(cond, qm.Limit(*last))
		}
	default:
		// 全件取得
		switch {
		case last != nil:
			// last が指定されていれば降順で取得
			scanDesc = true
			cond = append(cond,
				qm.OrderBy(fmt.Sprintf("%s desc", db.IssueColumns.ID)),
				qm.Limit(*last),
			)
		case first != nil:
			// first が指定されていれば昇順で取得
			cond = append(cond,
				qm.OrderBy(fmt.Sprintf("%s asc", db.IssueColumns.ID)),
				qm.Limit(*first),
			)
		default:
			// ページネーションのパラメータなしの時は昇順で取得
			cond = append(cond,
				qm.OrderBy(fmt.Sprintf("%s asc", db.IssueColumns.ID)),
			)
		}
	}

	issues, err := db.Issues(cond...).All(ctx, i.exec)
	if err != nil {
		return nil, err
	}

	var hasNextPage, hasPrevPage bool
	if len(issues) != 0 {
		// 降順の場合は昇順に並び替え
		if scanDesc {
			for i, j := 0, len(issues)-1; i < j; i, j = i+1, j-1 {
				issues[i], issues[j] = issues[j], issues[i]
			}
		}
		// 最初の要素のissueIDと最後の要素のissueID
		startCursor, endCursor := issues[0].ID, issues[len(issues)-1].ID

		var err error
		// 最初の要素のissueIDより小さなIDがあるか(前ページが存在するか)
		hasPrevPage, err = db.Issues(
			db.IssueWhere.Repository.EQ(repoID),
			db.IssueWhere.ID.LT(startCursor),
		).Exists(ctx, i.exec)
		if err != nil {
			return nil, err
		}
		// 最後の要素のissueIDより大きなIDがあるか(次ページが存在するか)
		hasNextPage, err = db.Issues(
			db.IssueWhere.Repository.EQ(repoID),
			db.IssueWhere.ID.GT(endCursor),
		).Exists(ctx, i.exec)
		if err != nil {
			return nil, err
		}
	}

	return convertIssueConnection(issues, hasPrevPage, hasNextPage), nil
}
