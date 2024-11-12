package service

type FaceRecognitionResult struct {
	Token        string `json:"token"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
	FaceCode     string `json:"face_code"`
}

func (s *JMService) SendFaceRecognitionCallback(result FaceRecognitionResult) error {
	var resp = map[string]interface{}{}
	if _, err := s.authClient.Post(FaceRecognitionURL, &result, &resp); err != nil {
		return err
	}
	return nil
}
