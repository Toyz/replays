package replays

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"strings"
	"time"
)

/*
 public teResourceGUID MapGuid;
                public teResourceGUID HeroGuid;
                public teResourceGUID SkinGuid;
                public long Timestamp;
                public ulong UserId;
                public ReplayType Type;
                public int QualityPct;
*/

type Replay struct {
	ReplayName    string
	Map           Map
	Hero          Hero
	Skin          Skin
	Timestamp     time.Time
	UserId        uint64
	ReplayType    ReplayType
	ReplayQuality ReplayQuality
}

func Parse(inFile []byte) (Replay, error) {
	replay := Replay{}
	name, payload := processAtoms(inFile)
	if len(payload) == 0 {
		return replay, errors.New("unknown payload")
	}

	data, err := base64.StdEncoding.DecodeString(payload[1])
	if err != nil {
		return replay, err
	}

	buf := bytes.NewBuffer(data)

	var mapGUID ResourceGuid
	if err := binary.Read(buf, binary.LittleEndian, &mapGUID); err != nil {
		return replay, err
	}
	mapGUID = (mapGUID &^ 0xFFFFFFFF00000000) | 0x0790000000000000
	replay.Map = newMap(mapGUID)

	var heroGUID ResourceGuid
	if err := binary.Read(buf, binary.LittleEndian, &heroGUID); err != nil {
		return replay, err
	}
	replay.Hero = newHero(heroGUID)

	var skinGUID ResourceGuid
	if err := binary.Read(buf, binary.LittleEndian, &skinGUID); err != nil {
		return replay, err
	}
	replay.Skin = newSkin(skinGUID)

	var TimeStamp int64
	if err := binary.Read(buf, binary.LittleEndian, &TimeStamp); err != nil {
		return replay, err
	}
	t := time.Unix(TimeStamp, 0)
	replay.Timestamp = t

	if err := binary.Read(buf, binary.LittleEndian, &replay.UserId); err != nil {
		return replay, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &replay.ReplayType); err != nil {
		return replay, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &replay.ReplayQuality); err != nil {
		return replay, err
	}

	replay.ReplayName = name
	return replay, nil
}

func processAtoms(buf []byte) (filename string, payload []string) {
	var size = 0

	for size < len(buf) {
		atom := NewAtom(buf[size:])
		size += int(atom.Size)

		if atom.Name == "moov" || atom.Name == "udta" {
			return processAtoms(atom.Buffer)
		}

		if atom.Name == "meta" {
			filename = processMetaName(atom.Buffer)
			continue
		}

		if atom.Name == "Xtra" {
			payload = processXtra(atom.Buffer)
			continue
		}
	}

	return
}

func processMetaName(buf []byte) string {
	size := 0
	for size < len(buf) {
		a := NewAtom(buf[size:])
		size += int(math.Max(float64(4), float64(a.Size)))

		if a.Name == "ilst" || a.Name == string([]byte{169, 110, 97, 109}) {
			return processMetaName(a.Buffer)
		}

		if a.Name == "data" {
			return string(a.Buffer[8:])
		}
	}

	return ""
}

func processXtra(in []byte) []string {
	if len(in) < 0x1F {
		log.Print("replay file is messed up")
		return nil
	}

	buf := bytes.NewBuffer(in)

	var blockSize int32
	if err := binary.Read(buf, binary.BigEndian, &blockSize); err != nil {
		log.Printf("failed to decode size: %f", err)
		return nil
	}

	var blockNameLength int32
	if err := binary.Read(buf, binary.BigEndian, &blockNameLength); err != nil {
		log.Printf("failed to decode length: %f", err)
		return nil
	}
	if blockNameLength == 0 {
		return nil
	}

	nameBlock := make([]byte, blockNameLength)
	if err := binary.Read(buf, binary.BigEndian, nameBlock); err != nil {
		log.Panicf("failed to decode length: %f", err)
	}
	name := string(nameBlock)
	if name != "WM/EncodingSettings" {
		log.Print("replay is corrupted")
		return nil
	}

	var settingsCount int32
	if err := binary.Read(buf, binary.BigEndian, &settingsCount); err != nil {
		log.Printf("failed to decode settingsCount: %f", err)
		return nil
	}

	for i := 0; i < int(settingsCount); i++ {
		var encodedSettingLength int32
		if err := binary.Read(buf, binary.BigEndian, &encodedSettingLength); err != nil {
			log.Printf("failed to decode encodedSettingsLength: %f", err)
			continue
		}
		if encodedSettingLength == 0 {
			continue
		}

		var unknownType int16
		if err := binary.Read(buf, binary.BigEndian, &unknownType); err != nil {
			log.Printf("failed to decode unknownType: %f", err)
			continue
		}

		b64Block := make([]byte, encodedSettingLength-6)
		if err := binary.Read(buf, binary.BigEndian, b64Block); err != nil {
			log.Printf("failed to decode b64Block: %f", err)
			continue
		}

		data := string(bytes.Replace(b64Block, []byte{0}, []byte{' '}, -1))
		return strings.Split(strings.Replace(data, " ", "", -1), ":")
	}

	return nil
}
