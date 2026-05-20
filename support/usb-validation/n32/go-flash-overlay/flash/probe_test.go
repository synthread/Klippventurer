package flash

import "testing"

func TestSTMApplyGetResponseMapsVersionAndCommands(t *testing.T) {
	mc := &Microcontroller{}
	mc.stmApplyGetResponse([]byte{0x31, 0x00, 0x01, 0x02, 0x11, 0x21})

	if mc.stmBootloaderVersion != 0x31 {
		t.Fatalf("expected bootloader version 0x31, got 0x%02x", mc.stmBootloaderVersion)
	}
	if got := mc.stmCmdCodes[CommandCodeGet]; got != 0x00 {
		t.Fatalf("expected GET command byte 0x00, got 0x%02x", got)
	}
	if got := mc.stmCmdCodes[CommandCodeGetID]; got != 0x02 {
		t.Fatalf("expected GET ID mapped byte 0x02 from response order, got 0x%02x", got)
	}
	if got := mc.stmCmdCodes[CommandCodeGo]; got != 0x21 {
		t.Fatalf("expected GO mapped byte 0x21, got 0x%02x", got)
	}
}

func TestCommandBytesFromGetSkipsBootloaderVersion(t *testing.T) {
	got := commandBytesFromGet([]byte{0x31, 0x00, 0x02})
	want := []byte{0x00, 0x02}
	if len(got) != len(want) {
		t.Fatalf("expected %d command bytes, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("byte %d: expected 0x%02x, got 0x%02x", i, want[i], got[i])
		}
	}
}

func TestCommandBytesFromGetHandlesMissingCommands(t *testing.T) {
	if got := commandBytesFromGet([]byte{0x31}); got != nil {
		t.Fatalf("expected nil command bytes, got %#v", got)
	}
}
