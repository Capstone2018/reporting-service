package handlers

const (
	headerContentType                = "Content-Type"
	headerAccessControlAllowOrigin   = "Access-Control-Allow-Origin"
	headerAccessControlAllowHeaders  = "Access-Control-Allow-Headers"
	headerAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	headerAccessControlAllowMethods  = "Access-Control-Allow-Methods"
)

const (
	charsetUTF8              = "; charset=utf-8"
	contentTypeJSON          = "application/json"
	contentTypeJSONUTF8      = contentTypeJSON + charsetUTF8
	contentTypeTextPlain     = "text/plain"
	contentTypeTextPlainUTF8 = contentTypeTextPlain + charsetUTF8
	contentTypeTextHTML      = "text/html"
	contentTypeTextHTMLUTF8  = contentTypeTextHTML + charsetUTF8
)
