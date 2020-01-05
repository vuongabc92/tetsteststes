package template

import (
	"flag"
	"github.com/vuongabc92/octocv/config"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/session"
	"strings"
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
	pathSplit := strings.Split(assetPath, ".")
	if pathSplit[len(pathSplit)-1] == "css" || pathSplit[len(pathSplit)-1] == "js" {
		return "/assets/frontend/" + assetPath + "?v=" + config.AssetVersion
	}

	return "/assets/frontend/" + assetPath
}

func FrontendFullAsset(assetPath string) string {
	flag.Parse()
	return *config.BaseUrl + "/assets/frontend/" + assetPath + "?v=" + config.AssetVersion
}

func SupportEmailAddress() string {
	flag.Parse()
	return *config.SupportEmailAddress
}
