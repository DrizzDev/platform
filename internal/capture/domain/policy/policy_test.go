package policy_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/policy"
)

func TestPolicy(test *testing.T) {
	test.Parallel()

	rule, failure := policy.New(policy.Input{Limit: 1024, Retention: time.Hour, Upload: true, Redaction: false})
	if failure != nil {
		test.Fatal(failure)
	}
	if rule.Limit() != 1024 || rule.Retention() != time.Hour || !rule.Upload() || rule.Redaction() {
		test.Fatalf("policy = %d/%s/%v/%v", rule.Limit(), rule.Retention(), rule.Upload(), rule.Redaction())
	}
}

func TestPolicyRejects(test *testing.T) {
	test.Parallel()

	rejected := map[string]policy.Input{
		"limit":     {Limit: 0, Retention: time.Hour},
		"retention": {Limit: 1, Retention: 0},
	}
	for name, input := range rejected {
		if _, failure := policy.New(input); failure == nil {
			test.Fatalf("%s policy was accepted", name)
		}
	}
}
