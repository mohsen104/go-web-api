package service_errors

type ServiceError struct {
	EndUserMessages   string `json:"endUserMessage"`
	TechnicalMessages string `json:"technicalMessage"`
	Err               error
}

func (s *ServiceError) Error() string {
	return s.EndUserMessages
}
