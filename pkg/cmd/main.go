package main

import (
	"encoding/hex"
	"fmt"
	"frame_reductor/model"
	"net"
	"os"
	"time"
)

func main() {

	// Przykładowe dane ramki konfiguracyjnej (muszą być rzeczywiste)
	frameData := []byte{
		// ramka konfiguracyjna 1
		//0xaa, 0x52, 0x00, 0xb6, 0x00, 0x65, 0x67, 0x1e, 0xc1, 0xb8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x01, 0x06, 0x53, 0x48, 0x45, 0x4c, 0x42, 0x59, 0x00,
		//0x02, 0xb6, 0x8e, 0x96, 0x27, 0xc9, 0x58, 0x48, 0x35, 0xae, 0x84, 0xee, 0x58, 0x8d, 0x8a, 0x72,
		//0x5f, 0x00, 0x0f, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x05, 0x42, 0x75, 0x73, 0x20, 0x31, 0x05,
		//0x42, 0x75, 0x73, 0x20, 0x32, 0x07, 0x43, 0x6f, 0x72, 0x64, 0x6f, 0x76, 0x61, 0x04, 0x44, 0x65,
		//0x6c, 0x6c, 0x0c, 0x4c, 0x61, 0x67, 0x6f, 0x6f, 0x6e, 0x20, 0x43, 0x72, 0x65, 0x65, 0x6b, 0x00,
		//0x08, 0x01, 0x00, 0x3f, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x3f, 0x80, 0x00,
		//0x00, 0x00, 0x00, 0x00, 0x08, 0x0c, 0x00, 0x3f, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x09,
		//0x00, 0x3f, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x09, 0x00, 0x3f, 0x80, 0x00, 0x00, 0x00,
		//0x42, 0x16, 0x00, 0x00, 0xc2, 0xc5, 0x33, 0x33, 0x7f, 0x80, 0x00, 0x00, 0x4d, 0x00, 0x00, 0x00,
		//0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0x2d, 0x97, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// ramka konfiguracyjna 2
		//0xaa, 0x52, 0x00, 0xb6, 0x00, 0x65, 0x67, 0x26, 0x73, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//0xff, 0xff, 0xff, 0xff, 0x00, 0x01, 0x06, 0x53, 0x48, 0x45, 0x4c, 0x42, 0x59, 0x00, 0x02, 0xb6,
		//0x8e, 0x96, 0x27, 0xc9, 0x58, 0x48, 0x35, 0xae, 0x84, 0xee, 0x58, 0x8d, 0x8a, 0x72, 0x5f, 0x00,
		//0x0f, 0x00, 0x05, 0x00, 0x00, 0x00, 0x05, 0x42, 0x75, 0x73, 0x20, 0x31, 0x05, 0x42, 0x75, 0x73,
		//0x20, 0x32, 0x07, 0x43, 0x6f, 0x72, 0x64, 0x6f, 0x76, 0x61, 0x04, 0x44, 0x65, 0x6c, 0x6c, 0x0c,
		//0x4c, 0x61, 0x67, 0x6f, 0x6f, 0x6e, 0x20, 0x43, 0x72, 0x65, 0x65, 0x6b, 0x00, 0x08, 0x01, 0x00,
		//0x3f, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x01, 0x00, 0x3f, 0x80, 0x00, 0x00, 0x00, 0x00,
		//0x00, 0x08, 0x0c, 0x00, 0x3f, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x09, 0x00, 0x3f, 0x80,
		//0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x09, 0x00, 0x3f, 0x80, 0x00, 0x00, 0x00, 0x42, 0x16, 0x00,
		//0x00, 0xc2, 0xc5, 0x33, 0x33, 0x7f, 0x80, 0x00, 0x00, 0x4d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//0x00, 0x00, 0x19, 0x28, 0x83, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// ramka z danymi 1
		0xaa, 0x01, 0x00, 0x42, 0x00, 0x65, 0x67, 0x1e, 0xc1, 0x70, 0x00, 0xae, 0x14, 0x7a, 0x00, 0x00,
		0x48, 0x92, 0x4b, 0xe7, 0x3f, 0xd4, 0x72, 0xe7, 0x48, 0x91, 0xce, 0x86, 0x3f, 0xd4, 0x82, 0x28,
		0x43, 0x65, 0x87, 0x0a, 0xbf, 0xd2, 0xab, 0x61, 0x43, 0xfb, 0xd8, 0x2b, 0x3f, 0xd4, 0xf1, 0xc7,
		0x43, 0x54, 0xe7, 0xd8, 0x3f, 0x4e, 0xaa, 0xb8, 0x42, 0x6f, 0xda, 0x1d, 0xbf, 0x21, 0x47, 0xae,
		0x63, 0x91, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	header, err := model.DecodeC37Header(frameData[:18])
	if err != nil {
		fmt.Println("Błąd dekodowania nagłówka:", err)
		return
	}
	fmt.Printf("Header: %v", header)

	switch header.DataFrameType {
	case model.ConfigurationFrame3:
		// Dekodowanie ramki konfiguracyjnej
		cfgFrame3, err := model.DecodeConfigurationFrame3(frameData)
		if err != nil {
			fmt.Println("Błąd dekodowania ramki konfiguracyjnej 3:", err)
			return
		}
		fmt.Printf("Decoded configuration frame 3: %+v\n", cfgFrame3)

	case model.DataFrame:
		// Dekodowanie ramki z danymi
		dataFrame, err := model.DecodeDataFrame(frameData)
		if err != nil {
			fmt.Println("Błąd dekodowania ramki z danymi:", err)
			return
		}
		fmt.Printf("Decoded data frame: %+v\n", dataFrame)
	}

	// Wyświetlenie informacji o ramce konfiguracyjnej
	//fmt.Printf("Configuration Frame:\n")
	//fmt.Printf("ID Code: %d\n", cfg.IDCode)
	//fmt.Printf("Frame Size: %d\n", cfg.FrameSize)
	//fmt.Printf("Frame Type: %d\n", cfg.FrameType)
	//fmt.Printf("Num PMUs: %d\n", cfg.NumPMUs)
	//fmt.Printf("DataPhasor Names:\n")
	//for i, name := range cfg.PhasorNames {
	//	fmt.Printf("  DataPhasor %d: %s\n", i+1, name)
	//}

	// Adres lokalny na porcie 4716
	addr := net.UDPAddr{
		Port: 4716,
		IP:   net.ParseIP("0.0.0.0"),
	}

	// Otwieramy gniazdo UDP
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Błąd podczas otwierania gniazda:", err)
		return
	}
	defer conn.Close()

	// Otwieramy plik do zapisu ramek
	file, err := os.Create("udp_frames.txt")
	if err != nil {
		fmt.Println("Błąd podczas tworzenia pliku:", err)
		return
	}
	defer file.Close()

	// Ustawiamy czas zakończenia nasłuchu
	timeout := time.After(180 * time.Second)

	fmt.Println("Nasłuchuję ramek UDP przez 300 sekund...")

	// Bufor do odczytu danych
	//buf := make([]byte, 1024)

loop:
	for {
		select {
		case <-timeout:
			fmt.Println("Czas nasłuchu upłynął.")
			break loop
		default:
			// Przykładowa ramka UDP (66 bajtów),
			// Zwiększony rozmiar bufora, aby uniknąć błędów związanych z dużymi ramkami
			frame := make([]byte, 1024)
			// Odbieramy dane UDP
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := conn.ReadFromUDP(frame)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // kontynuuj nasłuch po timeout
				}
				fmt.Println("Błąd podczas odczytu ramki: ", err)
				fmt.Println("Wykryta długość ramki: ", n)
				break loop
			}

			// Konwersja ramki do formatu hex
			hexFrame := hex.EncodeToString(frame)

			// Zapisujemy ramkę do pliku
			_, err = file.WriteString(hexFrame + "\n")
			if err != nil {
				fmt.Println("Błąd podczas zapisu do pliku:", err)
				break loop
			}
			fmt.Println("Odebrana ramka hex:", hexFrame)

			header, err := model.DecodeC37Header(frame[:18])
			if err != nil {
				fmt.Println("Błąd dekodowania nagłówka:", err)
				return
			}
			fmt.Printf("Header: %v", header)

			/*
				// Obliczanie czasu UTC
				calculatedTime := model.CalculateTimeUTC(header)
				fmt.Printf("Nagłówek: %+v   czas UTC: %v\n", header, calculatedTime)

				// Example 48 bytes after header, filled with sample data for testing
				data := make([]byte, 48)

				// Pozostałe dane (zakładając, że reszta ramki to 48 bajtów danych)
				data = frame[18:n]
				//fmt.Printf("Pozostałe dane: % X\n", data) // Wyświetla resztę ramki w formie heksadecymalnej

				// Decode the data fields, assuming 2 phasors, 2 analogs, and 1 digital word
				dataFields, err := model.DecodeDataFrame(data, 2, 2, 1)
				if err != nil {
					fmt.Println("Error decoding data fields:", err)
					return
				}

				fmt.Printf("Decoded C37DataFrame: %+v\n", dataFields)

			*/
		}
	}

	fmt.Println("Nasłuch zakończony, ramki zapisane do pliku udp_frames.txt.")
}
