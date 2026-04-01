package ocr

import (
	"strings"
	"testing"
)

func TestCleanExtractedTextRemovesWatermarkNoise(t *testing.T) {
	raw := `Bh
f7a249668c2366cc1Xx92tu_EVRSxY6_VvyWWOGgmPbSMBhg
朱冠州
⼿机号：13232212086
SM
Bh
SM
⾼级后台开发⼯程师
`

	cleaned := cleanExtractedText(raw)

	if cleaned == "" {
		t.Fatal("expected cleaned text to be non-empty")
	}
	if strings.Contains(cleaned, "f7a249668c2366cc1Xx92tu_EVRSxY6_VvyWWOGgmPbSMBhg") {
		t.Fatalf("expected watermark token to be removed, got %q", cleaned)
	}
	if strings.Contains(cleaned, "\nBh\n") {
		t.Fatalf("expected repeated short watermark fragment to be removed, got %q", cleaned)
	}
	if !strings.Contains(cleaned, "手机号:13232212086") {
		t.Fatalf("expected compatibility ideographs to be normalized, got %q", cleaned)
	}
	if !strings.Contains(cleaned, "高级后台开发工程师") {
		t.Fatalf("expected NFKC normalized Chinese text, got %q", cleaned)
	}
}

func TestTextQualityScorePrefersReadableResumeText(t *testing.T) {
	readable := cleanExtractedText(`朱冠州
手机号：13232212086
工作经历
高级后台开发工程师
项目经历
`)
	noisy := cleanExtractedText(`Bh
SM
Pb
gm
f7a249668c2366cc1Xx92tu_EVRSxY6_VvyWWOGgmPbSMBhg
`)

	readableScore := textQualityScore(readable)
	noisyScore := textQualityScore(noisy)

	if readableScore <= noisyScore {
		t.Fatalf("expected readable text score > noisy text score, got readable=%f noisy=%f", readableScore, noisyScore)
	}
}

func TestPickBestTextCandidatePrefersHigherQualityText(t *testing.T) {
	noisy := buildTextCandidate("go-pdf", `Bh
SM
Pb
f7a249668c2366cc1Xx92tu_EVRSxY6_VvyWWOGgmPbSMBhg`, 0.90)
	readable := buildTextCandidate("pdftotext", `朱冠州
手机号：13232212086
工作经历
项目经历`, 0.94)

	best := pickBestTextCandidate([]textCandidate{noisy, readable})
	if best.Method != "pdftotext" {
		t.Fatalf("expected pdftotext candidate to win, got %s", best.Method)
	}
}
