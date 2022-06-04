package api

type getWechatAuthUrlDto struct {
	Redirect string `binding:"required"`
}

type loginDto struct {
	Code string `binding:"required"`
}
