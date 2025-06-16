package Api

import (
	"ControlApp/Display"
	"ControlApp/Infrastructure"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type fetchResult struct {
	ConnectedDisplays []int
}

func (fixture Fixture) HandleDisplayConnectedApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	indices := make([]int, 0)
	for _, index := range fixture.Hardware.GetConnectedDisplays() {
		indices = append(indices, int(index))
	}

	//Encode data
	if err := json.NewEncoder(w).Encode(fetchResult{indices}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (fixture Fixture) HandleDisplayImportAnimationApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Limit file size to 100 MB. This line saves you from those accidental 100 MB uploads!
	err := r.ParseMultipartForm(10 << 24)
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
	}

	//Get animation ID
	var animationId uint32
	animationIdStr := r.FormValue("id")
	if animationIdStr != "" {
		tempId, err := strconv.ParseInt(animationIdStr, 10, 32)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing animation ID", http.StatusBadRequest)
			return
		}
		animationId = uint32(tempId)
	}

	info, err := os.Stat(fmt.Sprintf("blob/animations/%d", animationId))
	if err == nil && info.IsDir() {
		http.Error(w, "Animation with ID already exists.", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, _, err := r.FormFile("animationFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// Now let's save it locally
	dst, err := createTempFile()
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = dst.Close()
		_ = os.Remove(dst.Name())
	}()

	// Copy the uploaded file to the destination file
	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}

	//Convert animation
	err = Infrastructure.ExtractFrames(animationId, dst.Name())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting animation. error %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (fixture Fixture) HandleDisplayUploadAnimationApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//Get animation ID
	var animationId uint32
	animationIdStr := r.FormValue("id")
	if animationIdStr != "" {
		tempId, err := strconv.ParseInt(animationIdStr, 10, 32)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing animation ID.", http.StatusBadRequest)
			return
		}
		animationId = uint32(tempId)
	}

	//Get display byte
	var displayNr byte
	displayNrStr := r.FormValue("display")
	if displayNrStr != "" {
		tempId, err := strconv.ParseInt(displayNrStr, 10, 8)
		if err != nil || animationId < 0 {
			http.Error(w, "Error parsing display number.", http.StatusBadRequest)
			return
		}
		displayNr = byte(tempId)
	}

	dirPath := fmt.Sprintf("blob/animations/%d", animationId)
	if !exists(dirPath) {
		http.Error(w, "Animation does not exist.", http.StatusBadRequest)
	}

	//Read frames
	frames, err := Infrastructure.GetAnimationFrames(animationId)
	if err != nil {
		http.Error(w, "Error fetching animation.", http.StatusInternalServerError)
	}

	err = fixture.Hardware.DisplayServers.UploadAnimation(Display.AnimationId(animationId), frames, Display.ServerDisplay(displayNr))

	//Encode data
	if err != nil {
		http.Error(w, "Error uploading animation.", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (fixture Fixture) HandleDisplayPlayAnimationApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//Get animation ID
	var animationId uint32
	animationIdStr := r.FormValue("id")
	if animationIdStr != "" {
		tempId, err := strconv.ParseInt(animationIdStr, 10, 32)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing animation ID.", http.StatusBadRequest)
			return
		}
		animationId = uint32(tempId)
	}

	//Get display byte
	var displayNr byte
	displayNrStr := r.FormValue("display")
	if displayNrStr != "" {
		tempId, err := strconv.ParseInt(displayNrStr, 10, 8)
		if err != nil || animationId < 0 {
			http.Error(w, "Error parsing display number.", http.StatusBadRequest)
			return
		}
		displayNr = byte(tempId)
	}

	dirPath := fmt.Sprintf("blob/animations/%d", animationId)
	if !exists(dirPath) {
		http.Error(w, "Animation does not exist.", http.StatusBadRequest)
	}

	fixture.Hardware.SendAnimationInstruction(Display.AnimationId(animationId), []Display.ServerDisplay{Display.ServerDisplay(displayNr)})

	w.WriteHeader(http.StatusOK)
}

func (fixture Fixture) HandleDisplayShowTextApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//Get animation ID
	displayText := r.FormValue("text")

	//Get display byte
	var displayNr byte
	displayNrStr := r.FormValue("display")
	if displayNrStr != "" {
		tempId, err := strconv.ParseInt(displayNrStr, 10, 8)
		if err != nil {
			http.Error(w, "Error parsing display number.", http.StatusBadRequest)
			return
		}
		displayNr = byte(tempId)
	}

	fixture.Hardware.SendTextInstruction(displayText, []Display.ServerDisplay{Display.ServerDisplay(displayNr)})

	w.WriteHeader(http.StatusOK)
}

func (fixture Fixture) HandleDisplaySetBrightnessApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	//Get brightness value
	var valueNr float64
	valueNrStr := r.FormValue("value")
	if valueNrStr != "" {
		var err error
		valueNr, err = strconv.ParseFloat(valueNrStr, 64)
		if err != nil {
			http.Error(w, "Error parsing display number.", http.StatusBadRequest)
			return
		}
	}

	//Get brightness value
	var decrementNr uint16
	decrementNrStr := r.FormValue("decrement")
	if decrementNrStr != "" {
		nr, err := strconv.ParseInt(decrementNrStr, 10, 16)
		if err != nil {
			http.Error(w, "Error parsing display number.", http.StatusBadRequest)
			return
		}
		decrementNr = uint16(nr)
	}

	fixture.Hardware.SendBrightnessChange(&valueNr, decrementNr)

	w.WriteHeader(http.StatusOK)
}

func createTempFile() (*os.File, error) {
	// Create an uploads directory if it doesn’t exist
	if _, err := os.Stat("blob/temp"); os.IsNotExist(err) {
		err := os.MkdirAll("blob/temp", 0o775)
		if err != nil {
			return nil, err
		}
	}

	// Build the file path and create it
	dst, err := os.CreateTemp("blob/temp", "animation_*")
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
