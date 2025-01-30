// TubeGrab: Minimal YouTube Downloader Utility in Go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/kkdai/youtube/v2"
)

func main() {
	// Check if the user provided a YouTube URL
	if len(os.Args) < 2 {
		fmt.Println("Usage: tg <YouTube URL>")
		return
	}

	// Extract the YouTube URL from command-line arguments
	videoURL := os.Args[1]

	// Create a new YouTube client
	client := youtube.Client{}

	// Get video information
	video, err := client.GetVideo(videoURL)
	if err != nil {
		fmt.Printf("Error fetching video info: %v\n", err)
		return
	}

	// Print video details
	fmt.Printf("Downloading: %s\n", video.Title)
	fmt.Printf("Duration: %s\n", video.Duration)

	// Choose the best quality format
	format := video.Formats.WithAudioChannels().FindByQuality("medium")
	if format == nil {
		fmt.Println("No suitable format found.")
		return
	}

	// Create a stream for the video
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		fmt.Printf("Error getting video stream: %v\n", err)
		return
	}
	defer stream.Close()

	// Create the output file
	outputFile := filepath.Join(".", sanitizeFileName(video.Title)+".mp4")
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	// Download the video
	fmt.Printf("Downloading to: %s\n", outputFile)
	_, err = file.ReadFrom(stream)
	if err != nil {
		fmt.Printf("Error downloading video: %v\n", err)
		return
	}

	fmt.Println("Download completed successfully!")
}

// sanitizeFileName removes invalid characters from the file name
func sanitizeFileName(name string) string {
	invalidChars := `<>:"/\|?*`
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, string(char), "_")
	}
	return name
}