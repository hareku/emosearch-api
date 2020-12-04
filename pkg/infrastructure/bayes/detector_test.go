package bayes

import (
	"context"
	"log"
	"testing"

	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

func Test_bayesDetector_BatchDetect(t *testing.T) {
	d := NewBayesDetector()
	text1 := "これは最高に面白い映画です"
	text2 := "これは最悪で最低の映画です"

	output, err := d.BatchDetect(context.Background(), []*string{&text1, &text2})
	if err != nil {
		t.Fatalf("BatchDetect returned error: %s", err)
	}
	log.Printf("%v", output)

	if output[0].Label != sentiment.LabelPositive {
		t.Errorf("%q is labeled as %s", text1, output[0].Label)
	}

	if output[1].Label != sentiment.LabelNegative {
		t.Errorf("%q is labeled as %s", text2, output[1].Label)
	}
}
