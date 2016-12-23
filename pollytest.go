package main

import (
	"fmt"
	"log"
	"os"

	"github.com/niyeradori/polly-go"
)

func main() {
	client := polly.New("AWS_ACCESS_KEY", "AWS_SECRET_KEY")

	r, err := client.DescribeVoices()
	fmt.Println("returned back from ListVocies")
	if err != nil {
		fmt.Println("ListVocies returned error")
		log.Fatal(err)
	}

	log.Printf("Output of DescribeVoices: %v\n", r)

	mp3Filename := "from_polly.mp3"
	mp3file, err := os.Create(mp3Filename)
	if err != nil {
		fmt.Println("Error while creating", mp3Filename, "-", err)
		return
	}
	fmt.Printf("mp3 file path is %v\n", mp3Filename)

	defer mp3file.Close()

	options := polly.NewSpeechOptions("Hello this is Polly")
	r1, err1 := client.SynthesizeSpeech(options)

	if err1 != nil {
		fmt.Println(" Error on Return from CreateSpeech Function")
		log.Fatal(err1)
	}
	log.Printf("Audio Length = %v\n", len(r1.Audio))
	log.Printf("Audio Type = %v\n", r1.ContentType)
	log.Printf("No. of CharactersProcessed = %v\n", r1.RequestID)
	n2, err2 := mp3file.Write(r1.Audio)
	if err2 != nil {
		fmt.Println(" Error on MP3 Write")
		log.Fatal(err2)
	}
	fmt.Printf("wrote %d bytes\n", n2)
}
