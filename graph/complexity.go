package graph

import "mygql/internal"

func ComplexityConfig() internal.ComplexityRoot {
	var c internal.ComplexityRoot

	// Repositoryオブジェクト内のIssues数とクエリの複雑度を比例させる
	c.Repository.Issues = func(childComplexity int, after, before *string, first, last *int) int {
		var cnt int
		switch {
		case first != nil && last != nil :
			if *first < *last {
				cnt = *last
			} else {
				cnt = *first
			}
		case first != nil && last == nil:
			cnt = *first
		case first == nil && last != nil:
			cnt = *last
		default:
			cnt = 1
		}
		// childComplexity は取得されるIssuesのクエリ複雑度
		return cnt * childComplexity
	}
	return c
}
