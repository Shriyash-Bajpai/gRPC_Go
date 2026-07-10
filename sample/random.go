package sample

import (
	"math/rand"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/google/uuid"
)

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}

}

func randomCPUName(brand string) string {
	switch brand {
	case "Intel":
		return randomStringFromSet("i7", "i5", "i3", "Pentium")
	default:
		return randomStringFromSet("Ryzen 9", "Ryzen 7", "Ryzen 5")
	}
}

func randomGPUName(brand string) string {
	switch brand {
	case "NVIDIA":
		return randomStringFromSet("GTX 1600", "GTX 1650", "RTX 4600")
	default:
		return randomStringFromSet("Rx 590", "Rx 580", "Rx 5700-XT")
	}
}

func RandomLaptopBrand() string {
	brand := randomStringFromSet("Hewlett-Packard", "Dell", "Lenovo", "Apple")
	return brand
}

func RandomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("M3 Air", "M2 Air", "M4 Pro", "M3 Pro", "M2 Pro")
	case "Hewlett-Packard":
		return randomStringFromSet("Pavillion", "Victus", "Omen")
	case "Dell":
		return randomStringFromSet("G-15", "G-17")
	default:
		return randomStringFromSet("ThinkPad", "WorkPad", "GamingPad")
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}
func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	ran := rand.Intn(n-1) + 1
	return a[ran]
}

func randomInt(min int, max int) int {
	return (min + rand.Intn(max-min+1))
}

func RandomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func RandomScreenResolution() *pb.Screen_Resolution {

	height := randomInt(1080, 4320)
	width := height * 16 / 9
	resolution := &pb.Screen_Resolution{
		Height: uint32(height),
		Width:  uint32(width),
	}
	return resolution
}

func RandomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func RandomID() string {
	return uuid.New().String()
}
