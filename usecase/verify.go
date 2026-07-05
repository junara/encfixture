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
