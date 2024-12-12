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

type FaceMonitorResult struct {
	Token        string   `json:"token"`
	IsFinished   bool     `json:"is_finished"`
	Success      bool     `json:"success"`
	ErrorMessage string   `json:"error_message"`
	Action       string   `json:"action"`
	FaceCodes    []string `json:"face_codes"`
}

func (s *JMService) SendFaceMonitorCallback(result FaceMonitorResult) error {
	var resp = map[string]interface{}{}
	if _, err := s.authClient.Post(FaceMonitorURL, &result, &resp); err != nil {
		return err
	}
	return nil
}

type JoinFaceMonitorRequest struct {
	FaceMonitorToken string `json:"face_monitor_token"`
	SessionId        string `json:"session_id"`
}

func (s *JMService) JoinFaceMonitor(result JoinFaceMonitorRequest) error {
	var resp = map[string]interface{}{}
	if _, err := s.authClient.Post(FaceMonitorContextUrl, &result, &resp); err != nil {
		return err
	}
	return nil
}
