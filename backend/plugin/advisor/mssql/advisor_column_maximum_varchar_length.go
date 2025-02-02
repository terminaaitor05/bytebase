// Package mssql is the advisor for MSSQL database.
package mssql

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/bytebase/tsql-parser"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/advisor/db"
	bbparser "github.com/bytebase/bytebase/backend/plugin/parser/sql"
)

var (
	_ advisor.Advisor = (*ColumnMaximumVarcharLengthAdvisor)(nil)
)

func init() {
	advisor.Register(db.MSSQL, advisor.MSSQLColumnMaximumVarcharLength, &ColumnMaximumVarcharLengthAdvisor{})
}

// ColumnMaximumVarcharLengthAdvisor is the advisor checking for maximum varchar length..
type ColumnMaximumVarcharLengthAdvisor struct {
}

// Check checks for maximum varchar length..
func (*ColumnMaximumVarcharLengthAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	tree, ok := ctx.AST.(antlr.Tree)
	if !ok {
		return nil, errors.Errorf("failed to convert to Tree")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	payload, err := advisor.UnmarshalNumberTypeRulePayload(ctx.Rule.Payload)
	if err != nil {
		return nil, err
	}

	listener := &columnMaximumVarcharLengthChecker{
		level: level,
		title: string(ctx.Rule.Type),
		checkTypeString: map[string]any{
			"varchar":  nil,
			"nvarchar": nil,
			"char":     nil,
			"nchar":    nil,
		},
		maximum: payload.Number,
	}

	if listener.maximum > 0 {
		antlr.ParseTreeWalkerDefault.Walk(listener, tree)
	}

	return listener.generateAdvice()
}

// columnMaximumVarcharLengthChecker is the listener for maximum varchar length.
type columnMaximumVarcharLengthChecker struct {
	*parser.BaseTSqlParserListener

	level           advisor.Status
	title           string
	checkTypeString map[string]any
	maximum         int

	adviceList []advisor.Advice
}

// generateAdvice returns the advices generated by the listener, the advices must not be empty.
func (l *columnMaximumVarcharLengthChecker) generateAdvice() ([]advisor.Advice, error) {
	if len(l.adviceList) == 0 {
		l.adviceList = append(l.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return l.adviceList, nil
}

func (l *columnMaximumVarcharLengthChecker) EnterData_type(ctx *parser.Data_typeContext) {
	currentLength := 0
	line := ctx.GetStart().GetLine()
	if ctx.MAX() != nil && (ctx.VARCHAR() != nil || ctx.NVARCHAR() != nil) {
		// https://learn.microsoft.com/en-us/sql/t-sql/data-types/data-types-transact-sql?view=sql-server-ver16&redirectedfrom=MSDN
		currentLength = math.MaxInt32 // 2 ^ 31 - 1
		line = ctx.MAX().GetSymbol().GetLine()
	} else if ctx.GetExt_type() != nil && ctx.GetScale() != nil && ctx.GetPrec() == nil && ctx.GetInc() == nil {
		normalizedTypeString := bbparser.NormalizeTSQLIdentifier(ctx.GetExt_type())
		if _, ok := l.checkTypeString[normalizedTypeString]; !ok {
			return
		}
		length, err := strconv.Atoi(ctx.GetScale().GetText())
		if err != nil {
			slog.Error("failed to convert scale to int", log.BBError(err))
		}
		currentLength = length
		line = ctx.GetScale().GetLine()
	} else if ctx.GetUnscaled_type() != nil {
		normalizedTypeString := bbparser.NormalizeTSQLIdentifier(ctx.GetUnscaled_type())
		if _, ok := l.checkTypeString[normalizedTypeString]; !ok {
			return
		}
		line = ctx.GetUnscaled_type().GetStart().GetLine()
	}
	if currentLength > l.maximum {
		l.adviceList = append(l.adviceList, advisor.Advice{
			Status:  l.level,
			Code:    advisor.VarcharLengthExceedsLimit,
			Title:   l.title,
			Content: fmt.Sprintf("The maximum varchar length is %d.", l.maximum),
			Line:    line,
		})
	}
}
