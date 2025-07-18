package Api

import (
	"ControlApp/Display"
	"ControlApp/Lightshow"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

//Endpoints:
//GetAnimations
//UploadPalette
//DeletePalette

type animationsResponse struct {
	Animations []animationHeader `json:"animations"`
}

type animationHeader struct {
	Id            uint32 `json:"id"`
	Name          string `json:"name"`
	ThumbnailPath string `json:"thumbnail"`
	Mood          string `json:"mood"`
	IsNsfw        bool   `json:"nsfw"`
}

type animationUploaded struct {
	Id uint32 `json:"id"`
}

func (fixture Fixture) HandleAnimationsGetAllApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	var palettes []animationHeader

	for _, animation := range fixture.Data.Visuals.GetAnimations().GetAll() {
		palettes = append(palettes, animationHeader{
			uint32(animation.Id),
			animation.Name,
			fmt.Sprintf("/static/thumbs/%d.png", animation.Id),
			getMoodStr(animation.Mood),
			animation.IsNsfw,
		})
	}

	//Encode data
	if err := json.NewEncoder(w).Encode(animationsResponse{palettes}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getMoodStr(mood Lightshow.LightingMood) string {
	switch mood {
	case Lightshow.Moody:
		return "Moody"
	case Lightshow.Happy:
		return "Happy"
	case Lightshow.Regular:
		return "Regular"
	case Lightshow.Party:
		return "Party"
	default:
		return "Unknown"
	}
}

func (fixture Fixture) HandleSingleAnimationApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fixture.handleAnimationImportApi(w, r)
		break
	case "DELETE":
		fixture.handleAnimationDeleteApi(w, r)
		break
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

func (fixture Fixture) handleAnimationImportApi(w http.ResponseWriter, r *http.Request) {
	// Limit file size to 100 MB. This line saves you from those accidental 100 MB uploads!
	err := r.ParseMultipartForm(10 << 24)
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
	}

	var moodNr uint8
	moodNrStr := r.FormValue("mood")
	if moodNrStr != "" {
		tempId, err := strconv.ParseInt(moodNrStr, 10, 8)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing mood.", http.StatusBadRequest)
			return
		}
		moodNr = uint8(tempId)
	}

	isSplitStr := r.FormValue("split")
	isSplit := isSplitStr == "on"

	isNsfwStr := r.FormValue("nsfw")
	isNsfw := isNsfwStr == "on"

	nameStr := r.FormValue("name")

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
	id, err := fixture.Data.Visuals.ImportAnimation(dst.Name(), nameStr, Lightshow.LightingMood(moodNr), isSplit, isNsfw)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error importing animation. error %s", err), http.StatusInternalServerError)
	}

	returnData := animationUploaded{uint32(id)}

	//Encode data
	if err := json.NewEncoder(w).Encode(returnData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (fixture Fixture) handleAnimationDeleteApi(w http.ResponseWriter, r *http.Request) {
	var id Display.AnimationId
	idStr := r.FormValue("id")
	if idStr != "" {
		tempId, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil || tempId < 0 {
			http.Error(w, "Error parsing ID.", http.StatusBadRequest)
			return
		}
		id = Display.AnimationId(tempId)
	} else {
		http.Error(w, "ID not specified.", http.StatusBadRequest)
		return
	}

	exists, _ := fixture.Data.Visuals.GetAnimations().GetById(id)

	if !exists {
		http.Error(w, "Animation does not exist.", http.StatusBadRequest)
		return
	}

	fixture.Data.Visuals.GetAnimations().RemoveAnimation(id)
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
