package BoxiBus

import "errors"

type DisplayStatusCode byte
type LightingModeId byte

const (
	Booting              DisplayStatusCode = 0x00
	HostAwake            DisplayStatusCode = 0x01
	HostNoActivity       DisplayStatusCode = 0x02
	DisplayServerFailed  DisplayStatusCode = 0x03
	HostConnectionFailed DisplayStatusCode = 0x04
	Active               DisplayStatusCode = 0x05
)

const (
	Off                    LightingModeId = 0x00
	SetColor               LightingModeId = 0x01
	FadeToColor            LightingModeId = 0x02
	PaletteFade            LightingModeId = 0x03
	PaletteSwitch          LightingModeId = 0x04
	PaletteBrightnessFlash LightingModeId = 0x05
	PaletteHueFlash        LightingModeId = 0x06
	Strobe                 LightingModeId = 0x07
)

type Color struct {
	Red         byte
	Green       byte
	Blue        byte
	White       byte
	Amber       byte
	UltraViolet byte
}

func CreateDisplayStatusUpdate(statusCode DisplayStatusCode) MessageBlock {
	message := BusMessage{StatusCode, []byte{byte(statusCode)}}
	return []BusMessage{message}
}

func CreateLightingOff(applyOnBeat bool) MessageBlock {
	modeMessage := BusMessage{LightingMode, []byte{byte(Off)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return []BusMessage{modeMessage, applyMessage}
}

func CreateLightingSetColor(boxi1 Color, boxi2 Color, applyOnBeat bool) MessageBlock {
	colorAMessage := BusMessage{LightingPaletteA + 0, convertColor(boxi1)}
	colorBMessage := BusMessage{LightingPaletteA + 1, convertColor(boxi2)}
	modeMessage := BusMessage{LightingMode, []byte{byte(SetColor)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return []BusMessage{colorAMessage, colorBMessage, modeMessage, applyMessage}
}

func CreateLightingFadeToColor(boxi1 Color, boxi2 Color, speed uint16, applyOnBeat bool) MessageBlock {
	colorAMessage := BusMessage{LightingPaletteA + 0, convertColor(boxi1)}
	colorBMessage := BusMessage{LightingPaletteA + 1, convertColor(boxi2)}
	speedMessage := BusMessage{LightingSpeed, convertShort(speed)}
	modeMessage := BusMessage{LightingMode, []byte{byte(FadeToColor)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return []BusMessage{colorAMessage, colorBMessage, speedMessage, modeMessage, applyMessage}
}

func CreateLightingPaletteFade(palette []Color, speed uint16, paletteShift byte, applyOnBeat bool) (MessageBlock, error) {
	paletteMessages, err := convertPalette(palette)
	if err != nil {
		return nil, err
	}

	speedMessage := BusMessage{LightingSpeed, convertShort(speed)}
	shiftMessage := BusMessage{LightingColorShift, []byte{paletteShift}}
	modeMessage := BusMessage{LightingMode, []byte{byte(PaletteFade)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return append(paletteMessages, speedMessage, shiftMessage, modeMessage, applyMessage), nil
}

func CreateLightingPaletteSwitch(palette []Color, paletteShift byte, applyOnBeat bool) (MessageBlock, error) {
	paletteMessages, err := convertPalette(palette)
	if err != nil {
		return nil, err
	}

	shiftMessage := BusMessage{LightingColorShift, []byte{paletteShift}}
	modeMessage := BusMessage{LightingMode, []byte{byte(PaletteFade)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return append(paletteMessages, shiftMessage, modeMessage, applyMessage), nil
}

func CreateLightingPaletteBrightnessFlash(palette []Color, fadeOutSpeed uint16, targetBrightness byte, paletteShift byte,
	applyOnBeat bool) (MessageBlock, error) {
	paletteMessages, err := convertPalette(palette)
	if err != nil {
		return nil, err
	}

	speedMessage := BusMessage{LightingSpeed, convertShort(fadeOutSpeed)}
	gpMessage := BusMessage{LightingGeneralPurpose, []byte{targetBrightness}}
	shiftMessage := BusMessage{LightingColorShift, []byte{paletteShift}}
	modeMessage := BusMessage{LightingMode, []byte{byte(PaletteBrightnessFlash)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return append(paletteMessages, speedMessage, gpMessage, shiftMessage, modeMessage, applyMessage), nil
}

func CreateLightingPaletteHueFlash(palette []Color, fadeOutSpeed uint16, targetBrightness byte, paletteShift byte, applyOnBeat bool) (MessageBlock, error) {
	paletteMessages, err := convertPalette(palette)
	if err != nil {
		return nil, err
	}

	speedMessage := BusMessage{LightingSpeed, convertShort(fadeOutSpeed)}
	gpMessage := BusMessage{LightingGeneralPurpose, []byte{targetBrightness}}
	shiftMessage := BusMessage{LightingColorShift, []byte{paletteShift}}
	modeMessage := BusMessage{LightingMode, []byte{byte(PaletteHueFlash)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return append(paletteMessages, speedMessage, gpMessage, shiftMessage, modeMessage, applyMessage), nil
}

func CreateLightingStrobe(color Color, frequency uint16, rolloff byte, applyOnBeat bool) MessageBlock {

	colorMessage := BusMessage{LightingPaletteA + 0, convertColor(color)}
	speedMessage := BusMessage{LightingSpeed, convertShort(frequency)}
	gpMessage := BusMessage{LightingGeneralPurpose, []byte{rolloff}}
	modeMessage := BusMessage{LightingMode, []byte{byte(Strobe)}}
	applyMessage := BusMessage{LightingApply, convertBool(applyOnBeat)}
	return []BusMessage{colorMessage, speedMessage, gpMessage, modeMessage, applyMessage}
}

func convertShort(short uint16) []byte {
	return []byte{byte(short >> 8), byte(short & 0xff)}
}

func convertBool(input bool) []byte {
	if input {
		return []byte{1}
	} else {
		return []byte{0}
	}

}

func convertPalette(palette []Color) ([]BusMessage, error) {
	paletteLen := len(palette)
	if paletteLen > 8 {
		return nil, errors.New("palette length cannot exceed 8")
	}

	colorMessages := make([]BusMessage, paletteLen+1)
	colorMessages[0] = BusMessage{LightingPaletteSize, []byte{byte(paletteLen)}}
	for i := 0; i < paletteLen; i++ {
		colorMessages[i+1] = BusMessage{LightingPaletteA + MemoryField(i), convertColor(palette[i])}
	}
	return colorMessages, nil
}

func convertColor(color Color) []byte {
	return []byte{
		color.Red,
		color.Green,
		color.Blue,
		color.White,
		color.Amber,
		color.UltraViolet,
	}
}
