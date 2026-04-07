package evaluator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestoreInterviewQuestionsFromTruncatedOutput(t *testing.T) {
	parsed := map[string]interface{}{
		"录用建议": map[string]interface{}{
			"面试题目": []interface{}{},
		},
	}

	rawOutput := `{"录用建议":{"薪资建议":"30K-35K/月","面试题目":[` +
		`{"题目":"Q1","考察点":"F1","参考答案要点":["A1","A2"]},` +
		`{"题目":"Q2","考察点":"F2","参考答案要点":["B1"]},` +
		`{"题目":"Q3","考察点":"F3","参考答案要点":[`

	recovered := restoreInterviewQuestions(parsed, rawOutput)
	require.Equal(t, 2, recovered)

	rec := parsed["录用建议"].(map[string]interface{})
	questions := rec["面试题目"].([]interface{})
	require.Len(t, questions, 2)

	q1 := questions[0].(map[string]interface{})
	require.Equal(t, "Q1", q1["题目"])
	require.Equal(t, "F1", q1["考察点"])
}

func TestRestoreInterviewQuestionsKeepsExistingQuestions(t *testing.T) {
	parsed := map[string]interface{}{
		"录用建议": map[string]interface{}{
			"面试题目": []interface{}{
				map[string]interface{}{"题目": "已有题目"},
			},
		},
	}

	rawOutput := `{"录用建议":{"面试题目":[{"题目":"Q1"},{"题目":"Q2"}]}}`

	recovered := restoreInterviewQuestions(parsed, rawOutput)
	require.Equal(t, 0, recovered)

	rec := parsed["录用建议"].(map[string]interface{})
	questions := rec["面试题目"].([]interface{})
	require.Len(t, questions, 1)
	q := questions[0].(map[string]interface{})
	require.Equal(t, "已有题目", q["题目"])
}
