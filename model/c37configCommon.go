package model

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

// PhasorScaleFactor reprezentuje współczynnik konwersji dla kanałów fazorów z dodatkowymi flagami.
type PhasorScaleFactor struct {
	Flags           map[string]bool `json:"flags"`            // Flagi z mapowaniem bitowym
	PhasorType      string          `json:"phasor_type"`      // Typ fazora: Voltage lub Current
	PhasorComponent string          `json:"phasor_component"` // Komponent fazora (np. Phase A, Phase B)
	ScaleFactor     float32         `json:"scale_factor"`     // Współczynnik skali
	AngleOffset     float32         `json:"angle_offset"`     // Przesunięcie kąta
}

// AnalogScaleFactor reprezentuje współczynnik konwersji dla kanałów analogowych.
type AnalogScaleFactor struct {
	MagnitudeScale float32 `json:"magnitude_scale"` // Współczynnik skali wielkości w formacie IEEE 32-bit
	Offset         float32 `json:"offset"`          // Przesunięcie w formacie IEEE 32-bit
}

// DigitalMask reprezentuje maskę dla cyfrowych słów statusu.
type DigitalMask struct {
	Mask1 uint16 `json:"mask1"` // Pierwsza maska cyfrowa (16 bitów)
	Mask2 uint16 `json:"mask2"` // Druga maska cyfrowa (16 bitów)
}

// Definicje stałych dla bitu 0 w polu FORMAT
const (
	PhasorMagnitudeAndAngle = 0 // 0: magnitude i angle (polar)
	PhasorRealAndImaginary  = 1 // 1: real i imaginary (rectangular)
)

// FormatBits struktura reprezentująca bity pola FORMAT
type FormatBits struct {
	FREQ_DFREQ uint8 // Bit 3: Format częstotliwości DFREQ (0: 16-bit, 1: floating point)
	AnalogFmt  uint8 // Bit 2: Format analogowy (0: 16-bit, 1: floating point)
	PhasorFmt  uint8 // Bit 1: Format fazorów (0: 16-bit, 1: floating point)
	PhasorType uint8 // Bit 0: Typ fazora (0: magnitude i angle/polar, 1: real i imaginary/rectangular)
}

// Funkcja dekodująca bity pola FORMAT na strukturę FormatBits
func decodeFormatBits(format uint16) FormatBits {
	return FormatBits{
		FREQ_DFREQ: uint8((format >> 3) & 1), // Bit 3
		AnalogFmt:  uint8((format >> 2) & 1), // Bit 2
		PhasorFmt:  uint8((format >> 1) & 1), // Bit 1
		PhasorType: uint8(format & 1),        // Bit 0
	}
}

// TimeBaseBits struktura reprezentująca bity pola TIME_BASE
type TimeBaseBits struct {
	Reserved       uint32 // Bits 31-15: Zarezerwowane, zawsze 0
	TimeMultiplier uint32 // Bits 14-0: Mnożnik podstawy czasu
}

// DecodeTimeBase - funkcja dekodująca bity pola TIME_BASE na strukturę TimeBaseBits
func DecodeTimeBase(timeBase uint32) TimeBaseBits {
	return TimeBaseBits{
		Reserved:       (timeBase >> 15) & 0x1FFFF, // Bits 31-15
		TimeMultiplier: timeBase & 0x7FFF,          // Bits 14-0
	}
}

func DecodeCHNAMWithOffsetAndLength(reader *bytes.Reader, totalNames int) ([]string, error) {
	channelNames := make([]string, totalNames)

	// Przesunięcie o 1 bajt
	if _, err := reader.Seek(-1, io.SeekCurrent); err != nil {
		return nil, fmt.Errorf("error applying offset: %v", err)
	}

	for i := 0; i < totalNames; i++ {
		// Odczytaj długość nazwy (1 bajt)
		nameLen, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error reading length of channel name at index %d: %v", i, err)
		}

		// Odczytaj nazwę o długości określonej w nameLen
		name := make([]byte, nameLen)
		if _, err := reader.Read(name); err != nil {
			return nil, fmt.Errorf("error reading channel name at index %d: %v", i, err)
		}

		// Konwertuj bajty na string i dodaj do listy nazw
		channelNames[i] = string(name)
	}

	return channelNames, nil
}

// DecodeChannelNames - funkcja dekodująca nazwy kanałów dla Frame2
func DecodeChannelNames(reader *bytes.Reader, phnmr, annmr, dgnmr uint16) ([]string, error) {
	// Oblicz całkowitą liczbę kanałów
	totalChannels := int(phnmr) + int(annmr) + int(dgnmr)*16

	// Inicjalizuj listę nazw kanałów
	channelNames := make([]string, totalChannels)

	// Iteruj przez liczbę kanałów
	for i := 0; i < totalChannels; i++ {
		// Odczytaj 16 bajtów dla każdego kanału
		name := make([]byte, 16)
		if _, err := reader.Read(name); err != nil {
			return nil, fmt.Errorf("błąd odczytu nazwy kanału: %v", err)
		}

		// Konwertuj na string i usuń null bajty
		channelNames[i] = strings.TrimRight(string(name), "\x00")
	}

	return channelNames, nil
}

// ChannelType definiuje typ kanału jako napięcie lub prąd.
type ChannelType int

const (
	Voltage ChannelType = iota // 0 = napięcie
	Current                    // 1 = prąd
)

// PhasorUnit reprezentuje współczynnik konwersji dla kanałów fazorów.
type PhasorUnit struct {
	ChannelType      ChannelType // Typ kanału: napięcie lub prąd
	ConversionFactor float64     // Współczynnik konwersji (w 10^-5 V lub A na bit)
}

func DecodePhasorUnits(reader *bytes.Reader, phnmr uint16) ([]PhasorUnit, error) {
	phasorUnits := make([]PhasorUnit, phnmr)

	for i := uint16(0); i < phnmr; i++ {
		// Odczyt 4 bajtów dla kanału
		var rawData [4]byte
		if _, err := reader.Read(rawData[:]); err != nil {
			return nil, fmt.Errorf("błąd odczytu PHUNIT dla kanału %d: %v", i+1, err)
		}

		// Interpretacja pierwszego bajtu jako typ kanału
		var channelType ChannelType
		switch rawData[0] {
		case 0: // 0000 - Napięcie
			channelType = Voltage
		case 128: // 10000000 - Prąd
			channelType = Current
		default:
			return nil, fmt.Errorf("nieznany typ kanału %d w PHUNIT", rawData[0])
		}

		// Odczyt współczynnika konwersji (3 ostatnie bajty jako 24-bitowa liczba bez znaku)
		conversionFactor := float64(uint32(rawData[1])<<16|uint32(rawData[2])<<8|uint32(rawData[3])) * 1e-5

		phasorUnits[i] = PhasorUnit{
			ChannelType:      channelType,
			ConversionFactor: conversionFactor,
		}
	}

	return phasorUnits, nil
}

// AnalogType reprezentuje typ kanału analogowego
type AnalogType string

const (
	SinglePointOnWave AnalogType = "SinglePointOnWave" // 0
	RMS               AnalogType = "RMS"               // 1
	Peak              AnalogType = "Peak"              // 2
	Reserved          AnalogType = "Reserved"          // 5–64
	UserDefined       AnalogType = "UserDefined"       // 65–255
	Unknown           AnalogType = "Unknown"           // Nieznany typ
)

// AnalogUnit przechowuje dane o kanale analogowym
type AnalogUnit struct {
	ChannelType   AnalogType
	ScalingFactor float64
}

func DecodeAnalogUnits(reader *bytes.Reader, annmr uint16) ([]AnalogUnit, error) {
	analogUnits := make([]AnalogUnit, annmr)

	for i := uint16(0); i < annmr; i++ {
		// Odczyt 4 bajtów dla kanału
		var rawData [4]byte
		if _, err := reader.Read(rawData[:]); err != nil {
			return nil, fmt.Errorf("błąd odczytu ANUNIT dla kanału %d: %v", i+1, err)
		}

		// Interpretacja pierwszego bajtu jako typ kanału
		var channelType AnalogType
		switch {
		case rawData[0] == 0:
			channelType = SinglePointOnWave
		case rawData[0] == 1:
			channelType = RMS
		case rawData[0] == 2:
			channelType = Peak
		case rawData[0] >= 5 && rawData[0] <= 64:
			channelType = Reserved
		case rawData[0] >= 65: // <= 255
			channelType = UserDefined
		default:
			channelType = Unknown
		}

		// Odczyt współczynnika konwersji (3 ostatnie bajty jako 24-bitowa liczba ze znakiem)
		rawScalingFactor := int32(int8(rawData[1]))<<16 | int32(rawData[2])<<8 | int32(rawData[3])
		scalingFactor := float64(rawScalingFactor)

		analogUnits[i] = AnalogUnit{
			ChannelType:   channelType,
			ScalingFactor: scalingFactor,
		}
	}

	return analogUnits, nil
}

type DigitalUnit struct {
	NormalStatusMask uint16 // Maska normalnego stanu
	ValidInputsMask  uint16 // Maska aktualnie ważnych wejść
}

func DecodeDigitalUnits(reader *bytes.Reader, dgnmr uint16) ([]DigitalUnit, error) {
	digitalUnits := make([]DigitalUnit, dgnmr)

	for i := uint16(0); i < dgnmr; i++ {
		// Odczyt dwóch 16-bitowych wartości (4 bajty)
		var normalStatusMask, validInputsMask uint16

		if err := binary.Read(reader, binary.BigEndian, &normalStatusMask); err != nil {
			return nil, fmt.Errorf("błąd odczytu NormalStatusMask: %v", err)
		}

		if err := binary.Read(reader, binary.BigEndian, &validInputsMask); err != nil {
			return nil, fmt.Errorf("błąd odczytu ValidInputsMask: %v", err)
		}

		// Dodanie odczytanych wartości do wyniku
		digitalUnits[i] = DigitalUnit{
			NormalStatusMask: normalStatusMask,
			ValidInputsMask:  validInputsMask,
		}
	}

	return digitalUnits, nil
}

// DecodeFlags dekoduje flagi na podstawie wartości uint16, zwracając mapę opisującą ustawione flagi
func DecodeFlags(flags uint16) map[string]bool {
	return map[string]bool{
		"reserved":                  (flags & 0x0001) != 0, // Bit 0: Zarezerwowane (nieużywane)
		"upsampled_with_interpol":   (flags & 0x0002) != 0, // Bit 1: Próbkowanie w górę za pomocą interpolacji
		"upsampled_with_extrapol":   (flags & 0x0004) != 0, // Bit 2: Próbkowanie w górę za pomocą ekstrapolacji
		"downsampled_with_reselect": (flags & 0x0008) != 0, // Bit 3: Próbkowanie w dół z wyborem próbek
		"downsampled_with_fir":      (flags & 0x0010) != 0, // Bit 4: Próbkowanie w dół z filtrem FIR
		"downsampled_non_fir":       (flags & 0x0020) != 0, // Bit 5: Próbkowanie w dół bez użycia filtra FIR
		"filtered_without_sampling": (flags & 0x0040) != 0, // Bit 6: Filtracja bez zmiany próbkowania
		"magnitude_adjusted":        (flags & 0x0080) != 0, // Bit 7: Dopasowanie wielkości
		"phase_adjusted_rotation":   (flags & 0x0100) != 0, // Bit 8: Dopasowanie fazy przez rotację
		"pseudo_phasor":             (flags & 0x0400) != 0, // Bit 10: Pseudofazor
		"modification_applied":      (flags & 0x8000) != 0, // Bit 15: Zastosowano modyfikację
	}
}

// DecodePhasorScale dekoduje PhasorScale z danych binarnych
func DecodePhasorScale(reader *bytes.Reader, count int) ([]PhasorScaleFactor, error) {
	phasorScales := make([]PhasorScaleFactor, count)

	for i := 0; i < count; i++ {
		var flags uint16
		var phasorTypeAndComponent uint8
		var reserved uint8
		var scaleFactor uint32
		var angleOffset uint32

		// Odczyt pierwszego 4-bajtowego słowa (flags + typ + komponent)
		if err := binary.Read(reader, binary.BigEndian, &flags); err != nil {
			return nil, fmt.Errorf("Błąd odczytu BitMappedFlags dla PhasorScale: %v", err)
		}

		if err := binary.Read(reader, binary.BigEndian, &phasorTypeAndComponent); err != nil {
			return nil, fmt.Errorf("Błąd odczytu Typu i Komponentu dla PhasorScale: %v", err)
		}

		if err := binary.Read(reader, binary.BigEndian, &reserved); err != nil {
			return nil, fmt.Errorf("Błąd odczytu Reserved dla PhasorScale: %v", err)
		}

		// Rozkodowanie flags
		decodedFlags := DecodeFlags(flags)

		// Rozbicie phasorTypeAndComponent na typ i komponent fazora
		phasorType := "voltage"
		if (phasorTypeAndComponent>>3)&0x01 == 1 {
			phasorType = "current"
		}

		phasorComponent := map[uint8]string{
			0b000: "zero sequence",
			0b001: "positive sequence",
			0b010: "negative sequence",
			0b011: "reserved",
			0b100: "phase A",
			0b101: "phase B",
			0b110: "phase C",
			0b111: "reserved",
		}[phasorTypeAndComponent&0x07]

		// Odczyt drugiego 4-bajtowego słowa - Scale Factor
		if err := binary.Read(reader, binary.BigEndian, &scaleFactor); err != nil {
			return nil, fmt.Errorf("Błąd odczytu ScaleFactor dla PhasorScale: %v", err)
		}

		// Odczyt trzeciego 4-bajtowego słowa - Angle Offset
		if err := binary.Read(reader, binary.BigEndian, &angleOffset); err != nil {
			return nil, fmt.Errorf("Błąd odczytu AngleOffset dla PhasorScale: %v", err)
		}

		// Konwersja ScaleFactor i AngleOffset do float32
		scaleFactorFloat := math.Float32frombits(scaleFactor)
		angleOffsetFloat := math.Float32frombits(angleOffset)

		phasorScales[i] = PhasorScaleFactor{
			Flags:           decodedFlags,
			PhasorType:      phasorType,
			PhasorComponent: phasorComponent,
			ScaleFactor:     scaleFactorFloat,
			AngleOffset:     angleOffsetFloat,
		}
	}
	return phasorScales, nil
}

// DecodeAnalogScale dekoduje AnalogScale z danych binarnych.
func DecodeAnalogScale(reader *bytes.Reader, count int) ([]AnalogScaleFactor, error) {
	analogScales := make([]AnalogScaleFactor, count)
	for i := 0; i < count; i++ {
		var scale AnalogScaleFactor
		if err := binary.Read(reader, binary.BigEndian, &scale.MagnitudeScale); err != nil {
			return nil, fmt.Errorf("Błąd odczytu MagnitudeScale dla AnalogScale: %v", err)
		}
		if err := binary.Read(reader, binary.BigEndian, &scale.Offset); err != nil {
			return nil, fmt.Errorf("Błąd odczytu Offset dla AnalogScale: %v", err)
		}
		analogScales[i] = scale
	}
	return analogScales, nil
}

// DecodeDigitalMask dekoduje maski cyfrowe z pola DIGUNIT o długości 4 bajtów.
func DecodeDigitalMask(reader *bytes.Reader, numDigitals uint16) ([]DigitalMask, error) {
	digitalMasks := make([]DigitalMask, numDigitals)

	numDigitals = 0 //TODO temporary hardcoded

	for i := 0; i < int(numDigitals); i++ {
		var mask DigitalMask

		// Pierwsze 2 bajty - maska 1
		if err := binary.Read(reader, binary.BigEndian, &mask.Mask1); err != nil {
			return nil, fmt.Errorf("Błąd odczytu Mask1 dla DigitalMask: %v", err)
		}

		// Kolejne 2 bajty - maska 2
		if err := binary.Read(reader, binary.BigEndian, &mask.Mask2); err != nil {
			return nil, fmt.Errorf("Błąd odczytu Mask2 dla DigitalMask: %v", err)
		}

		digitalMasks[i] = mask
	}

	return digitalMasks, nil
}

// FNom reprezentuje nominalną częstotliwość linii
type FNom struct {
	Is50Hz bool // true, jeśli częstotliwość podstawowa wynosi 50 Hz
	Is60Hz bool // true, jeśli częstotliwość podstawowa wynosi 60 Hz
}

// DecodeFreqNominal dekoduje FNOM na podstawie specyfikacji
func DecodeFreqNominal(reader *bytes.Reader) (*FNom, error) {
	// Odczyt 2 bajtów jako uint16
	var rawFNom uint16
	if err := binary.Read(reader, binary.BigEndian, &rawFNom); err != nil {
		return nil, fmt.Errorf("błąd odczytu FNOM: %v", err)
	}

	// Dekodowanie bitów
	fNom := &FNom{
		Is50Hz: rawFNom&0x0001 == 1, // Bit 0 ustawiony na 1 oznacza 50 Hz
		Is60Hz: rawFNom&0x0001 == 0, // Bit 0 ustawiony na 0 oznacza 60 Hz
	}

	return fNom, nil
}
