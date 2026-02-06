package audio

import "math"

const sampleRate = 44100

// GenerateDiceRoll creates a rattling dice sound.
func GenerateDiceRoll() []byte {
	duration := 0.6
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	lfsr := uint16(0xACE1)
	for i := 0; i < samples; i++ {
		progress := float64(i) / float64(samples)
		t := float64(i) / float64(sampleRate)

		// Noise via LFSR
		bit := ((lfsr >> 0) ^ (lfsr >> 2) ^ (lfsr >> 3) ^ (lfsr >> 5)) & 1
		lfsr = (lfsr >> 1) | (bit << 15)
		noise := float64(int16(lfsr)) / 32768.0

		// Rattling clicks at decreasing frequency
		clickFreq := 20.0 - 15.0*progress
		click := math.Sin(2*math.Pi*clickFreq*t) * 0.3

		val := noise*0.4 + click*0.6
		env := 1.0 - progress*0.7
		sample := int16(val * env * 4000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GeneratePurchase creates a pleasant cash register ching.
func GeneratePurchase() []byte {
	duration := 0.3
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		val := math.Sin(2*math.Pi*2000*t)*0.4 +
			math.Sin(2*math.Pi*3000*t)*0.3 +
			math.Sin(2*math.Pi*4000*t)*0.2

		env := 1.0 - progress
		env *= env
		sample := int16(val * env * 5000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateRent creates a coin dropping sound.
func GenerateRent() []byte {
	duration := 0.25
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		freq := 1500.0 + 500.0*math.Sin(2*math.Pi*20*t)
		val := math.Sin(2 * math.Pi * freq * t)

		env := 1.0 - progress
		sample := int16(val * env * 4000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateCardDraw creates a card flip sound.
func GenerateCardDraw() []byte {
	duration := 0.15
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	lfsr := uint16(0xBEEF)
	for i := 0; i < samples; i++ {
		progress := float64(i) / float64(samples)

		bit := ((lfsr >> 0) ^ (lfsr >> 2) ^ (lfsr >> 3) ^ (lfsr >> 5)) & 1
		lfsr = (lfsr >> 1) | (bit << 15)
		noise := float64(int16(lfsr)) / 32768.0

		env := 1.0 - progress
		env *= env
		sample := int16(noise * env * 3000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateJail creates an ominous low tone.
func GenerateJail() []byte {
	duration := 0.5
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		val := math.Sin(2*math.Pi*150*t)*0.5 +
			math.Sin(2*math.Pi*200*t)*0.3 +
			math.Sin(2*math.Pi*100*t)*0.2

		env := 1.0 - progress*0.5
		sample := int16(val * env * 6000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GeneratePassGo creates a cheerful ascending tone.
func GeneratePassGo() []byte {
	duration := 0.4
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	notes := []float64{523.25, 659.25, 783.99} // C5, E5, G5
	noteLen := samples / len(notes)

	for i := 0; i < samples; i++ {
		noteIdx := i / noteLen
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		val := math.Sin(2*math.Pi*freq*t)*0.7 +
			math.Sin(2*math.Pi*freq*2*t)*0.2

		env := 1.0 - progress*0.4
		sample := int16(val * env * 5000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateBuild creates a construction hammering sound.
func GenerateBuild() []byte {
	duration := 0.3
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		// Two taps
		tap1 := 0.0
		tap2 := 0.0
		if progress < 0.15 {
			tap1 = math.Sin(2*math.Pi*300*t) * (1.0 - progress/0.15)
		}
		if progress > 0.4 && progress < 0.55 {
			p2 := (progress - 0.4) / 0.15
			tap2 = math.Sin(2*math.Pi*350*t) * (1.0 - p2)
		}

		val := tap1 + tap2
		sample := int16(val * 6000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateBankruptcy creates a sad descending tone.
func GenerateBankruptcy() []byte {
	duration := 1.0
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		freq := 400.0 - 250.0*progress
		val := math.Sin(2*math.Pi*freq*t)*0.6 +
			math.Sin(2*math.Pi*freq*0.5*t)*0.3

		env := 1.0 - progress
		sample := int16(val * env * 7000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateWin creates a victory fanfare.
func GenerateWin() []byte {
	duration := 1.5
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	notes := []float64{392.0, 523.25, 659.25, 783.99, 1046.50} // G4, C5, E5, G5, C6
	noteLen := samples / len(notes)

	for i := 0; i < samples; i++ {
		noteIdx := i / noteLen
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		t := float64(i) / float64(sampleRate)
		localT := float64(i%noteLen) / float64(noteLen)

		val := math.Sin(2*math.Pi*freq*t)*0.6 +
			math.Sin(2*math.Pi*freq*2*t)*0.2 +
			math.Sin(2*math.Pi*freq*3*t)*0.1

		env := 1.0
		if localT > 0.8 {
			env = (1.0 - localT) * 5.0
		}

		sample := int16(val * env * 6000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}

// GenerateMenuSelect creates a short blip.
func GenerateMenuSelect() []byte {
	duration := 0.05
	samples := int(float64(sampleRate) * duration)
	buf := make([]byte, samples*2)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(samples)

		val := math.Sin(2 * math.Pi * 1000 * t)
		env := 1.0 - progress
		sample := int16(val * env * 5000)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}
	return buf
}
