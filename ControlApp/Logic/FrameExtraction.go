package Logic

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExtractFrames converts the specified source file into displayable frames.
func ExtractFrames(animationId uint32, sourcePath string) error {
	// Input and output file names
	rootPath := fmt.Sprintf("blob/animations/%d/", animationId)
	_ = os.RemoveAll(rootPath)
	err := os.MkdirAll(rootPath, 0o775)
	if err != nil {
		return err
	}

	outputFile := rootPath + "/%04d.png"
	cmd := exec.Command("./ffmpeg",
		"-i", sourcePath,
		"-vf", "fps=25,scale=160x128:force_original_aspect_ratio=increase,crop=160:128",
		"-q:v", "1",
		outputFile)

	// Run the command
	return cmd.Run()
}

func GetAnimationFrames(animationId uint32) ([]string, error) {
	rootPath := fmt.Sprintf("blob/animations/%d", animationId)

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
