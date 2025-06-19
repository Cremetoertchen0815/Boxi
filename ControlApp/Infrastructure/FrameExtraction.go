package Infrastructure

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExtractFrames converts the specified source file into displayable frames.
func ExtractFrames(animationId uint32, sourcePath string) error {
	// Input and output file names
	rootPath := fmt.Sprintf("animations/%d/", animationId)
	_ = os.RemoveAll(rootPath)
	err := os.MkdirAll(rootPath, 0o775)
	if err != nil {
		return err
	}

	outputFile := rootPath + "/%04d.png"
	cmd := exec.Command(
		"nice", "-n", "15",
		"./ffmpeg",
		"-i", sourcePath,
		"-vf", "fps=25,scale=160x128:force_original_aspect_ratio=increase,crop=160:128",
		"-q:v", "1",
		outputFile)

	// Run the command
	return cmd.Run()
}

// ExtractDoubleFrames converts the specified source file into two sets of displayable frames
// for left and right screens showing two halves of the same animation.
func ExtractDoubleFrames(leftAnimationId uint32, rightAnimationId uint32, sourcePath string) error {
	// Root paths
	leftPath := fmt.Sprintf("animations/%d/", leftAnimationId)
	rightPath := fmt.Sprintf("animations/%d/", rightAnimationId)

	_ = os.RemoveAll(leftPath)
	_ = os.RemoveAll(rightPath)

	// Create directories
	if err := os.MkdirAll(leftPath, 0o775); err != nil {
		return err
	}
	if err := os.MkdirAll(rightPath, 0o775); err != nil {
		return err
	}

	// Command to generate left-side frames
	cmdLeft := exec.Command(
		"nice", "-n", "15",
		"./ffmpeg",
		"-i", sourcePath,
		"-vf", "fps=25,scale=320x128:force_original_aspect_ratio=increase,crop=160:128:0:0",
		"-q:v", "1",
		leftPath+"%04d.png",
	)

	// Command to generate right-side frames
	cmdRight := exec.Command(
		"nice", "-n", "15",
		"./ffmpeg",
		"-i", sourcePath,
		"-vf", "fps=25,scale=320x128:force_original_aspect_ratio=increase,crop=160:128:160:0",
		"-q:v", "1",
		rightPath+"%04d.png",
	)

	// Execute both commands
	if err := cmdLeft.Run(); err != nil {
		return fmt.Errorf("left side extraction failed: %w", err)
	}
	if err := cmdRight.Run(); err != nil {
		return fmt.Errorf("right side extraction failed: %w", err)
	}

	return nil
}

func GetAnimationFrames(animationId uint32) ([]string, error) {
	rootPath := fmt.Sprintf("animations/%d", animationId)

	files, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		absPath := filepath.Join(rootPath, entry.Name())
		result = append(result, absPath)
	}

	return result, nil
}
