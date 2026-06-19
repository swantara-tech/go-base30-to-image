package decoder

import (
	"testing"
)

// ─── UncompressLeg ────────────────────────────────────────────────────────────

func TestUncompressLeg_SimplePositive(t *testing.T) {
	// "jz" in jSignature base30 charset
	// j = index 19, z = index 35 (tail of 'j' which is 9, mapped: q[j]=J, m[J]=j)
	// Let's test a known simple case: single digit "1" → value 1
	result := UncompressLeg("1")
	if len(result) != 1 || result[0] != 1 {
		t.Errorf("expected [1], got %v", result)
	}
}

func TestUncompressLeg_MultiDigit(t *testing.T) {
	// In jSignature: each lower char starts a new token
	// "10" = two tokens: 1, then 0+1=1 (cumulative delta: [1, 1])
	// To encode value 30 we'd need a "tail" (uppercase) char
	// e.g. 'A' is tail-remap of 'a' (index=10 in base-30)
	// "1A" = token '1' then tail 'A'→'a', so buf=[1,a]=[1,10]=1*30+10=40... not 30
	// Actually let's just test the delta-encoding behavior
	result := UncompressLeg("10")
	if len(result) != 2 {
		t.Fatalf("expected 2 tokens from '10', got %d: %v", len(result), result)
	}
	// first token: '1' → val=1*1+0=1, last=1
	// second token: '0' → val=0+last=0+1=1, last=1
	if result[0] != 1 || result[1] != 1 {
		t.Errorf("expected [1,1], got %v", result)
	}
}

func TestUncompressLeg_ZThenY(t *testing.T) {
	// "5Z3Y2" → [5, 5+(-3)=2, 2+2=4]
	result := UncompressLeg("5Z3Y2")
	if len(result) != 3 {
		t.Fatalf("expected 3 values, got %d: %v", len(result), result)
	}
	if result[0] != 5 {
		t.Errorf("expected result[0]=5, got %d", result[0])
	}
	if result[1] != 2 {
		t.Errorf("expected result[1]=2, got %d", result[1])
	}
	if result[2] != 4 {
		t.Errorf("expected result[2]=4, got %d", result[2])
	}
}

func TestUncompressLeg_TailChar(t *testing.T) {
	// charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWX"
	// base = 30, qMap[runes[r]] = runes[r+base]
	// '6' is index 6 → qMap['6'] = runes[6+30] = runes[36] = 'A'
	// So '6A': flush '6' → val=6, start new token? No: 'A' is tail.
	// '6' flushes (lower=new token start), buf=['6'], then 'A'→mMap['A']='6'
	// buf=['6','6'] = 6*30+6 = 186... let's verify: mMap maps runes[r+base]→runes[r]
	// runes[36]='A', r=6, mMap['A']='6' (charVal=6). buf=['6','6'] = 6*30+6=186. Nope.
	// Let's just test known real behavior: value from '1' flush, tail 'A'→'6'
	// '1A': buf after '1' starts = ['1'], then 'A'→mMap['A']='6', buf=['1','6']
	// val = 1*30+6 = 36, sign=1, last=0 → 36
	result := UncompressLeg("1A")
	if len(result) != 1 || result[0] != 36 {
		t.Errorf("expected [36], got %v", result)
	}
}

func TestUncompressLeg_Empty(t *testing.T) {
	result := UncompressLeg("")
	if len(result) != 0 {
		t.Errorf("expected [], got %v", result)
	}
}

// ─── UncompressStrokes ────────────────────────────────────────────────────────

func TestUncompressStrokes_WithPrefix(t *testing.T) {
	// minimal 1-stroke: xLeg="1"_yLeg="2"
	strokes := UncompressStrokes("image/jsignature;base30,1_2")
	if len(strokes) != 1 {
		t.Fatalf("expected 1 stroke, got %d", len(strokes))
	}
	if strokes[0].X[0] != 1 {
		t.Errorf("expected X[0]=1, got %d", strokes[0].X[0])
	}
	if strokes[0].Y[0] != 2 {
		t.Errorf("expected Y[0]=2, got %d", strokes[0].Y[0])
	}
}

func TestUncompressStrokes_TwoStrokes(t *testing.T) {
	// 2 strokes: "1_2_3_4"
	strokes := UncompressStrokes("1_2_3_4")
	if len(strokes) != 2 {
		t.Fatalf("expected 2 strokes, got %d", len(strokes))
	}
}

func TestUncompressStrokes_Empty(t *testing.T) {
	strokes := UncompressStrokes("")
	if len(strokes) != 0 {
		t.Errorf("expected empty, got %v", strokes)
	}
}

// ─── ParseSignature ───────────────────────────────────────────────────────────

func TestParseSignature_Empty(t *testing.T) {
	_, err := ParseSignature("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParseSignature_WithPrefix(t *testing.T) {
	sig, err := ParseSignature("image/jsignature;base30,5a_2b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig.Strokes) == 0 {
		t.Error("expected at least 1 stroke")
	}
}

func TestParseSignature_WithoutPrefix(t *testing.T) {
	sig, err := ParseSignature("5a_2b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig.Strokes) == 0 {
		t.Error("expected at least 1 stroke")
	}
}

func TestParseSignature_MultipleStrokes(t *testing.T) {
	sig, err := ParseSignature("5a_2b_3c_1d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig.Strokes) < 2 {
		t.Errorf("expected 2 strokes, got %d", len(sig.Strokes))
	}
}

func TestParseSignature_RealSample(t *testing.T) {
	// Real jSignature sample from user
	data := "image/jsignature;base30,jzZ448g74b0000Y166568477ebeba6_4u4lkqkf1Nc96653222020Z3432312_lAZ27645400Y122677a98l5Z1_5X5ba88987d353200000Z35_ly46698a693_5M121010112_kF8699a986456_7u00000000Z100_pK00Z2025330000002_5y7588afcba786654_pC55a78ba7a66341322210Z5354869a9da96Y488bia554_5vZ13000000Y211344333367521110000000Z2Y5466b4421"
	sig, err := ParseSignature(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sig.Strokes) == 0 {
		t.Error("expected strokes from real sample")
	}
	if sig.Width <= 0 || sig.Height <= 0 {
		t.Errorf("expected positive dimensions, got %dx%d", sig.Width, sig.Height)
	}
	t.Logf("Real sample: strokes=%d, width=%d, height=%d", len(sig.Strokes), sig.Width, sig.Height)
}

// ─── ValidateSignature ────────────────────────────────────────────────────────

func TestValidateSignature_Empty(t *testing.T) {
	if err := ValidateSignature(""); err == nil {
		t.Error("expected error for empty input")
	}
}

func TestValidateSignature_Valid(t *testing.T) {
	if err := ValidateSignature("image/jsignature;base30,5a_2b"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
