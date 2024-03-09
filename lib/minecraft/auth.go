package minecraft

import (
	GMMAuth "github.com/maxsupermanhd/go-mc-ms-auth"
	"log"
)

func LoginMicrosoftAccount() (GMMAuth.BotAuth, error) {
	msAuth, err := GMMAuth.GetMCcredentials("./auth.json", "88650e7e-efee-4857-b9a9-cf580a00ef43")
	if err != nil {
		log.Fatal("登录微软账号失败，将使用离线模式: ", err)
		return GMMAuth.BotAuth{}, err
	}
	log.Print("Authenticated as ", msAuth.Name, " (", msAuth.UUID, ")")
	return msAuth, nil
}
