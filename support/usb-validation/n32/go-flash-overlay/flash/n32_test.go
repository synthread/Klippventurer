package flash

import "testing"

func TestN32ValidateFlashRangeAcceptsFullFlash(t *testing.T) {
	mc := &Microcontroller{config: &Config{}}
	payload := make([]byte, n32FlashSize)

	if err := mc.n32ValidateFlashRange(payload, n32FlashBase); err != nil {
		t.Fatalf("expected full flash range to be valid: %v", err)
	}
}

func TestN32ValidateFlashRangeRejectsBeforeFlash(t *testing.T) {
	mc := &Microcontroller{config: &Config{}}

	if err := mc.n32ValidateFlashRange([]byte{0x00}, n32FlashBase-1); err == nil {
		t.Fatal("expected address before flash base to be rejected")
	}
}

func TestN32ValidateFlashRangeRejectsAfterFlash(t *testing.T) {
	mc := &Microcontroller{config: &Config{}}

	if err := mc.n32ValidateFlashRange([]byte{0x00}, n32FlashBase+n32FlashSize); err == nil {
		t.Fatal("expected address at flash end to be rejected")
	}
}

func TestN32ValidateFlashRangeRejectsOverflowPastFlash(t *testing.T) {
	mc := &Microcontroller{config: &Config{}}

	if err := mc.n32ValidateFlashRange([]byte{0x00, 0x01}, n32FlashBase+n32FlashSize-1); err == nil {
		t.Fatal("expected payload extending past flash end to be rejected")
	}
}

func TestN32ValidateFlashRangeRejectsEmptyPayload(t *testing.T) {
	mc := &Microcontroller{config: &Config{}}

	if err := mc.n32ValidateFlashRange(nil, n32FlashBase); err == nil {
		t.Fatal("expected empty payload to be rejected")
	}
}

func TestN32ValidateFlashRangeUsesConfiguredBounds(t *testing.T) {
	mc := &Microcontroller{config: &Config{
		TargetFlashBase: 0x08010000,
		TargetFlashSize: 0x00010000,
	}}

	if err := mc.n32ValidateFlashRange([]byte{0x00}, 0x08010000); err != nil {
		t.Fatalf("expected configured flash base to be valid: %v", err)
	}
	if err := mc.n32ValidateFlashRange([]byte{0x00}, 0x08000000); err == nil {
		t.Fatal("expected default flash base to be rejected when configured bounds are set")
	}
}

func TestFamilyIdentityPrefixPreservesSTMCompatibility(t *testing.T) {
	mc := &Microcontroller{config: &Config{}, target: TargetFamilySTM32}

	if got := mc.familyIdentityPrefix(); got != "STM_" {
		t.Fatalf("expected STM_ identity prefix, got %q", got)
	}
}

func TestFamilyIdentityPrefixUsesN32Family(t *testing.T) {
	mc := &Microcontroller{config: &Config{}, target: TargetFamilyN32G45}

	if got := mc.familyIdentityPrefix(); got != "N32G45x_" {
		t.Fatalf("expected N32G45x_ identity prefix, got %q", got)
	}
}
