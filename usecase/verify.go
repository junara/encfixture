package usecase

import (
	"fmt"

	"github.com/junara/encfixture/domain"
)

// VerifyUseCase inspects existing media files via a prober.
type VerifyUseCase struct {
	prober Prober
}

// NewVerifyUseCase creates a new VerifyUseCase with the given prober.
func NewVerifyUseCase(prober Prober) *VerifyUseCase {
	return &VerifyUseCase{prober: prober}
}

// VerifyWithExpectations probes the media file at path and evaluates the
// given expectations against it. The probed info is returned alongside the
// per-expectation results; expectation failures are reported in the results,
// not as an error.
func (uc *VerifyUseCase) VerifyWithExpectations(
	path string,
	exps []domain.Expectation,
) (domain.MediaInfo, []domain.CheckResult, error) {
	info, err := uc.Verify(path)
	if err != nil {
		return info, nil, err
	}

	return info, domain.EvaluateExpectations(info, exps), nil
}

// Verify returns the probed properties of the media file at path.
func (uc *VerifyUseCase) Verify(path string) (domain.MediaInfo, error) {
	var empty domain.MediaInfo

	err := uc.prober.CheckAvailable()
	if err != nil {
		return empty, fmt.Errorf("ffprobe availability check failed: %w", err)
	}

	info, probeErr := uc.prober.Probe(path)
	if probeErr != nil {
		return empty, fmt.Errorf("probe failed: %w", probeErr)
	}

	return info, nil
}
