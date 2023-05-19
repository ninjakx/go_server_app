package model

// struct
// func ()initValidator() *validator.Validate {
// 	v := validator.New()
// 	return v
// }

// func(*validator.Validate) RegisterIPValidation(){
// 	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
// 		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
// 	}, func(ut ut.Translator, fe validator.FieldError) string {
// 		t, _ := ut.T("required", fe.Field())
// 		return t
// 	})
// }
// _ = v.RegisterValidation("ip_address", func(fl validator.FieldLevel) bool {
// 	return net.ParseIP(fl.Field().String()) == nil
// })

// if err := v.Struct(server); err != nil {
// 	log.Printf("[server][CreateServer][struct] error:%+v\n", err)
// 	respondError(c, http.StatusBadRequest, err.Error())
// 	return
// }
