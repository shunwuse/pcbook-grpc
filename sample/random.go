package sample

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/warnshun/pcbook/pb"
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

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomCPUName(brand string) string {
	switch brand {
	case "Intel":
		return randomStringFromSet("i3", "i5", "i7")
	case "AMD":
		return randomStringFromSet("Ryzen 3", "Ryzen 5", "Ryzen 7")
	}
	return ""
}

func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD", "Intel")
}

func randomGPUName(brand string) string {
	switch brand {
	case "Nvidia":
		return randomStringFromSet("RTX 4060", "RTX 4070", "RTX 4080", "RTX 4090")
	case "AMD":
		return randomStringFromSet("RX 7600", "RX 7700", "RX 7800", "RX 7900")
	case "Intel":
		return randomStringFromSet("Arc A750", "Arc A770")
	}
	return ""
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "HP", "Lenovo", "Dell", "Asus")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("MacBook Air", "MacBook Pro")
	case "HP":
		return randomStringFromSet("Envy", "Omen", "Pavilion")
	case "Lenovo":
		return randomStringFromSet("IdeaPad", "Legion")
	case "Dell":
		return randomStringFromSet("Vostro", "XPS")
	case "Asus":
		return randomStringFromSet("TUF", "ROG")
	}
	return ""
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	resolution := &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
	return resolution
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomStringFromSet(set ...string) string {
	n := len(set)
	if n == 0 {
		return ""
	}
	return set[rand.Intn(n)]
}

func randomBoolean() bool {
	return rand.Intn(2) == 1
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomId() string {
	return uuid.New().String()
}
