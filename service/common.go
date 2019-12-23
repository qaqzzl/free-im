package service

type CommonService struct {

}


func (s *CommonService) IsPhoneVerifyCode() (bool, error) {
	return true, nil
}