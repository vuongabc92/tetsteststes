package template

import (
	"github.com/vuongabc92/octocv/config"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/session"
)

func MessageBagGet(msgBag *session.MessageBag, field string) string {
	if msgBag.Has(field) {
		return msgBag.Get(field)
	}

	return ""
}

func MessageBagHas(msgBag *session.MessageBag, field string) bool {
	return msgBag.Has(field)
}

func FormDataText(formData *session.FormData, field string) string {
	data := formData.Get(field)

	if len(data) > 0 {
		return data[0]
	}

	return ""
}

func Trans(key interface{}, agrs ...interface{}) string {
	return helpers.Trans(key, agrs)
}

func FrontendAsset(assetPath string) string {
	return "assets/frontend/" + assetPath + "?v=" + config.AssetVersion
}
