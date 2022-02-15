package main

import (
	"fmt"
	"github.com/Kagami/go-face"
	"log"
	"path/filepath"
)

func main() {
	fmt.Println("Facial Recognition System v0.01")

	rec, err := face.NewRecognizer("models")
	if err != nil {
		fmt.Println("Cannot INItialize recognizer")
		return
	}
	defer rec.Close()

	fmt.Println("Recognizer Initialized")

	avengersImage := filepath.Join("jay-zhou.jpeg")

	faces, err := rec.RecognizeFile(avengersImage)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	fmt.Println("Number of Faces in Image: ", len(faces))

	var samples []face.Descriptor
	var avengers []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so Goes to its own category.
		avengers = append(avengers, int32(i))
	}
	// Name the categories, i.e. people on the image.
	labels := []string{
		"周杰伦",
	}
	// Pass samples to the recognizer.
	rec.SetSamples(samples, avengers)

	// Now let's try to classify some not yet known image.
	testTonyStark := filepath.Join("jay.jpeg")
	tonyStark, err := rec.RecognizeSingleFile(testTonyStark)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if tonyStark == nil {
		log.Fatalf("Not a single face on the image")
	}
	avengerID := rec.ClassifyThreshold(tonyStark.Descriptor, 0.3)
	if avengerID < 0 {
		fmt.Println(avengerID)
		log.Fatalf("Can't classify")
	}

	fmt.Println(avengerID)

	fmt.Println(labels[avengerID])

}
