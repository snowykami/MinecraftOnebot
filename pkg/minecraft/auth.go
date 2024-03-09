package minecraft

// Go MC 微软认证

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// BotAuth 是 Tnze 的 bot.Auth 的文字替代
// 终有一天他会发疯并对其进行某些操作
// 而我不想在他之后修复它
type BotAuth struct {
	Name string
	UUID string
	AsTk string
}

// MSauth 包含 Microsoft 认证凭据
type MSauth struct {
	AccessToken  string
	ExpiresAfter int64
	RefreshToken string
}

// AzureClientIDEnvVar 用于通过 os.Getenv 查找 Azure 客户端 ID，如果未传递 cid
const AzureClientIDEnvVar = "AzureClientID"

// CheckRefreshMS 检查 MSauth 是否过期，并在需要时刷新
func CheckRefreshMS(auth *MSauth, cid string) error {
	if auth.ExpiresAfter > time.Now().Unix() {
		return nil
	}
	if cid == "" {
		cid = os.Getenv(AzureClientIDEnvVar)
	}
	if auth.RefreshToken == "" {
		return errors.New("MS 访问令牌已过期，且未提供刷新令牌")
	}
	MSdata := url.Values{
		"client_id": {cid},
		// "client_secret": {os.Getenv("AzureSecret")},
		"refresh_token": {auth.RefreshToken},
		"grant_type":    {"refresh_token"},
		"redirect_uri":  {"https://login.microsoftonline.com/common/oauth2/nativeclient"},
	}
	MSresp, err := http.PostForm("https://login.live.com/oauth20_token.srf", MSdata)
	if err != nil {
		return err
	}
	var MSres map[string]interface{}
	json.NewDecoder(MSresp.Body).Decode(&MSres)
	MSresp.Body.Close()
	if MSresp.StatusCode != 200 {
		return fmt.Errorf("MS 刷新尝试回应非HTTP200！而是收到 %s 和以下 JSON：%#v", MSresp.Status, MSres)
	}
	MSaccessToken, ok := MSres["access_token"].(string)
	if !ok {
		return errors.New("在响应中未找到 access_token")
	}
	auth.AccessToken = MSaccessToken
	MSrefreshToken, ok := MSres["refresh_token"].(string)
	if !ok {
		return errors.New("在响应中未找到 refresh_token")
	}
	auth.RefreshToken = MSrefreshToken
	MSexpireSeconds, ok := MSres["expires_in"].(float64)
	if !ok {
		return errors.New("在响应中未找到 expires_in")
	}
	auth.ExpiresAfter = time.Now().Unix() + int64(MSexpireSeconds)
	return nil
}

// AuthMSdevice 尝试通过设备流授权用户。将阻塞线程，直到出现错误、超时或实际授权
func AuthMSdevice(cid string) (MSauth, error) {
	var auth MSauth
	if cid == "" {
		cid = os.Getenv(AzureClientIDEnvVar)
	}
	DeviceResp, err := http.PostForm("https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode", url.Values{
		"client_id": {cid},
		"scope":     {`XboxLive.signin offline_access`},
	})
	if err != nil {
		return auth, err
	}
	var DeviceRes map[string]interface{}
	json.NewDecoder(DeviceResp.Body).Decode(&DeviceRes)
	DeviceResp.Body.Close()
	if DeviceResp.StatusCode != 200 {
		return auth, fmt.Errorf("MS 设备请求回应非HTTP200！而是收到 %s 和以下 JSON：%#v", DeviceResp.Status, DeviceRes)
	}
	DeviceCode, ok := DeviceRes["device_code"].(string)
	if !ok {
		return auth, errors.New("在响应中未找到设备代码")
	}
	UserCode, ok := DeviceRes["user_code"].(string)
	if !ok {
		return auth, errors.New("在响应中未找到用户代码")
	}
	log.Print("用户代码：", UserCode)
	VerificationURI, ok := DeviceRes["verification_uri"].(string)
	if !ok {
		return auth, errors.New("在响应中未找到验证 URI")
	}
	log.Print("验证 URI：", VerificationURI)
	ExpiresIn, ok := DeviceRes["expires_in"].(float64)
	if !ok {
		return auth, errors.New("在响应中未找到 expires_in")
	}
	log.Print("过期时间：", ExpiresIn, " 秒")
	PoolInterval, ok := DeviceRes["interval"].(float64)
	if !ok {
		return auth, errors.New("在响应中未找到轮询间隔")
	}
	UserMessage, ok := DeviceRes["message"].(string)
	if !ok {
		return auth, errors.New("在响应中未找到用户消息")
	}
	log.Println(UserMessage)
	time.Sleep(4 * time.Second)

	for {
		time.Sleep(time.Duration(int(PoolInterval)+1) * time.Second)
		CodeResp, err := http.PostForm("https://login.microsoftonline.com/consumers/oauth2/v2.0/token", url.Values{
			"client_id":   {cid},
			"scope":       {"XboxLive.signin offline_access"},
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
			"device_code": {DeviceCode},
		})
		if err != nil {
			return auth, err
		}
		var CodeRes map[string]interface{}
		json.NewDecoder(CodeResp.Body).Decode(&CodeRes)
		CodeResp.Body.Close()
		if CodeResp.StatusCode == 400 {
			PoolError, ok := CodeRes["error"].(string)
			if !ok {
				return auth, fmt.Errorf("在轮询令牌时收到未知的 JSON：%#v", CodeRes)
			}
			if PoolError == "authorization_pending" {
				continue
			}
			if PoolError == "authorization_declined" {
				return auth, errors.New("用户拒绝了授权")
			}
			if PoolError == "expired_token" {
				return auth, errors.New("原来 " + strconv.Itoa(int(PoolInterval)) + " 秒不足以授权用户，请更快点")
			}
			if PoolError == "invalid_grant" {
				return auth, errors.New("在轮询令牌时收到 invalid_grant 错误：" + CodeRes["error_description"].(string))
			}
		} else if CodeResp.StatusCode == 200 {
			MSaccessToken, ok := CodeRes["access_token"].(string)
			if !ok {
				return auth, errors.New("在响应中未找到 access_token")
			}
			auth.AccessToken = MSaccessToken
			MSrefreshToken, ok := CodeRes["refresh_token"].(string)
			if !ok {
				return auth, errors.New("在响应中未找到 refresh_token")
			}
			auth.RefreshToken = MSrefreshToken
			MSexpireSeconds, ok := CodeRes["expires_in"].(float64)
			if !ok {
				return auth, errors.New("在响应中未找到 expires_in")
			}
			auth.ExpiresAfter = time.Now().Unix() + int64(MSexpireSeconds)
			return auth, nil
		} else {
			return auth, fmt.Errorf("MS 回应非HTTP200！而是收到 %s 和以下 JSON：%#v", CodeResp.Status, CodeRes)
		}
	}
}

// AuthMSCode 尝试通过用户代码（默认浏览器流）授权用户
func AuthMSCode(code string, cid string) (MSauth, error) {
	var auth MSauth
	if cid == "" {
		cid = os.Getenv(AzureClientIDEnvVar)
	}
	MSdata := url.Values{
		"client_id": {cid},
		// "client_secret": {os.Getenv("AzureSecret")},
		"code":         {code},
		"grant_type":   {"authorization_code"},
		"redirect_uri": {"https://login.microsoftonline.com/common/oauth2/nativeclient"},
	}
	MSresp, err := http.PostForm("https://login.live.com/oauth20_token.srf", MSdata)
	if err != nil {
		return auth, err
	}
	var MSres map[string]interface{}
	json.NewDecoder(MSresp.Body).Decode(&MSres)
	MSresp.Body.Close()
	if MSresp.StatusCode != 200 {
		return auth, fmt.Errorf("MS 回应非HTTP200！而是收到 %s 和以下 JSON：%#v", MSresp.Status, MSres)
	}
	MSaccessToken, ok := MSres["access_token"].(string)
	if !ok {
		return auth, errors.New("在响应中未找到 access_token")
	}
	auth.AccessToken = MSaccessToken
	MSrefreshToken, ok := MSres["refresh_token"].(string)
	if !ok {
		return auth, errors.New("在响应中未找到 refresh_token")
	}
	auth.RefreshToken = MSrefreshToken
	MSexpireSeconds, ok := MSres["expires_in"].(float64)
	if !ok {
		return auth, errors.New("在响应中未找到 expires_in")
	}
	auth.ExpiresAfter = time.Now().Unix() + int64(MSexpireSeconds)
	return auth, nil
}

// AuthXBL 从 Microsoft 令牌获取 XBox Live 令牌
func AuthXBL(MStoken string) (string, error) {
	XBLdataMap := map[string]interface{}{
		"Properties": map[string]interface{}{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  "d=" + MStoken,
		},
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType":    "JWT",
	}
	XBLdata, err := json.Marshal(XBLdataMap)
	if err != nil {
		return "", err
	}
	XBLreq, err := http.NewRequest(http.MethodPost, "https://user.auth.xboxlive.com/user/authenticate", bytes.NewBuffer(XBLdata))
	if err != nil {
		return "", err
	}
	XBLreq.Header.Set("Content-Type", "application/json")
	XBLreq.Header.Set("Accept", "application/json")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
	XBLresp, err := client.Do(XBLreq)
	if err != nil {
		return "", err
	}
	var XBLres map[string]interface{}
	json.NewDecoder(XBLresp.Body).Decode(&XBLres)
	XBLresp.Body.Close()
	if XBLresp.StatusCode != 200 {
		return "", fmt.Errorf("XBL 回应非HTTP200！而是收到 %s 和以下 JSON：%#v", XBLresp.Status, XBLres)
	}
	XBLtoken, ok := XBLres["Token"].(string)
	if !ok {
		return "", errors.New("在 XBL 响应中未找到令牌")
	}
	return XBLtoken, nil
}

// XSTSauth 包含 XSTS 令牌和 UHS
type XSTSauth struct {
	Token string
	UHS   string
}

// AuthXSTS 使用 XBL 获取 XSTS 令牌
func AuthXSTS(XBLtoken string) (XSTSauth, error) {
	var auth XSTSauth
	XSTSdataMap := map[string]interface{}{
		"Properties": map[string]interface{}{
			"SandboxId":  "RETAIL",
			"UserTokens": []string{XBLtoken},
		},
		"RelyingParty": "rp://api.minecraftservices.com/",
		"TokenType":    "JWT",
	}
	XSTSdata, err := json.Marshal(XSTSdataMap)
	if err != nil {
		return auth, err
	}
	XSTSreq, err := http.NewRequest(http.MethodPost, "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewBuffer(XSTSdata))
	if err != nil {
		return auth, err
	}
	XSTSreq.Header.Set("Content-Type", "application/json")
	XSTSreq.Header.Set("Accept", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	XSTSresp, err := client.Do(XSTSreq)
	if err != nil {
		return auth, err
	}
	var XSTSres map[string]interface{}
	json.NewDecoder(XSTSresp.Body).Decode(&XSTSres)
	XSTSresp.Body.Close()
	if XSTSresp.StatusCode != 200 {
		return auth, fmt.Errorf("XSTS 回应非HTTP200！而是收到 %s 和以下 JSON：%#v", XSTSresp.Status, XSTSres)
	}
	XSTStoken, ok := XSTSres["Token"].(string)
	if !ok {
		return auth, errors.New("在 XSTS 响应中未找到令牌")
	}
	auth.Token = XSTStoken
	XSTSdc, ok := XSTSres["DisplayClaims"].(map[string]interface{})
	if !ok {
		return auth, errors.New("在 XSTS 响应中未找到 DisplayClaims 对象")
	}
	XSTSxui, ok := XSTSdc["xui"].([]interface{})
	if !ok {
		return auth, errors.New("在 DisplayClaims 对象中未找到 xui 数组")
	}
	if len(XSTSxui) < 1 {
		return auth, errors.New("在 DisplayClaims 对象的 xui 数组中没有任何元素")
	}
	XSTSuhsObject, ok := XSTSxui[0].(map[string]interface{})
	if !ok {
		return auth, errors.New("在 xui 数组中无法获取 ush 对象")
	}
	XSTSuhs, ok := XSTSuhsObject["uhs"].(string)
	if !ok {
		return auth, errors.New("无法从 ush 对象中获取 uhs 字符串")
	}
	auth.UHS = XSTSuhs
	return auth, nil
}

// MCauth 代表 Minecraft 授权响应
type MCauth struct {
	Token        string
	ExpiresAfter int64
}

// AuthMC 从 XSTS 令牌获取 Minecraft 授权
func AuthMC(token XSTSauth) (MCauth, error) {
	var auth MCauth
	MCdataMap := map[string]interface{}{
		"identityToken": "XBL3.0 x=" + token.UHS + ";" + token.Token,
	}
	MCdata, err := json.Marshal(MCdataMap)
	if err != nil {
		return auth, err
	}
	MCreq, err := http.NewRequest(http.MethodPost, "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewBuffer(MCdata))
	if err != nil {
		return auth, err
	}
	MCreq.Header.Set("Content-Type", "application/json")
	MCreq.Header.Set("Accept", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	MCresp, err := client.Do(MCreq)
	if err != nil {
		return auth, err
	}
	var MCres map[string]interface{}
	json.NewDecoder(MCresp.Body).Decode(&MCres)
	MCresp.Body.Close()
	if MCresp.StatusCode != 200 {
		return auth, fmt.Errorf("MC 回应非HTTP200！而是收到 %s 和以下 JSON：%#v", MCresp.Status, MCres)
	}
	MCtoken, ok := MCres["access_token"].(string)
	if !ok {
		return auth, errors.New("在 MC 响应中未找到 access_token")
	}
	auth.Token = MCtoken
	MCexpire, ok := MCres["expires_in"].(float64)
	if !ok {
		return auth, errors.New("在 MC 响应中未找到 expires_in")
	}
	auth.ExpiresAfter = time.Now().Unix() + int64(MCexpire)
	return auth, nil
}

// GetMCprofile 从令牌获取 BotAuth
func GetMCprofile(token string) (BotAuth, error) {
	var profile BotAuth
	PRreq, err := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
	if err != nil {
		return profile, err
	}
	PRreq.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	PRresp, err := client.Do(PRreq)
	if err != nil {
		return profile, err
	}
	var PRres map[string]interface{}
	json.NewDecoder(PRresp.Body).Decode(&PRres)
	PRresp.Body.Close()
	if PRresp.StatusCode != 200 {
		return profile, fmt.Errorf("MC（配置文件）回应非HTTP200！而是收到 %s 和以下 JSON：%#v", PRresp.Status, PRres)
	}
	PRuuid, ok := PRres["id"].(string)
	if !ok {
		return profile, errors.New("在配置文件响应中未找到 uuid")
	}
	profile.UUID = PRuuid
	PRname, ok := PRres["name"].(string)
	if !ok {
		return profile, errors.New("在配置文件响应中未找到用户名")
	}
	profile.Name = PRname
	return profile, nil
}

// DefaultCacheFilename 用于使用设备代码流从 0 到 Minecraft BotAuth 的缓存加载和保存 Microsoft 授权，因为它提供的令牌的持续时间从一天到一周不等
const DefaultCacheFilename = "./auth.cache"

// GetMCcredentials 从缓存使用设备代码流获取 Minecraft 授权
func GetMCcredentials(CacheFilename, cid string) (BotAuth, error) {
	var resauth BotAuth
	var MSa MSauth
	if CacheFilename == "" {
		CacheFilename = DefaultCacheFilename
	}
	if _, err := os.Stat(CacheFilename); os.IsNotExist(err) {
		var err error
		MSa, err = AuthMSdevice(cid)
		if err != nil {
			return resauth, err
		}
		tocache, err := json.Marshal(MSa)
		if err != nil {
			return resauth, err
		}
		err = os.WriteFile(CacheFilename, tocache, 0600)
		if err != nil {
			return resauth, err
		}
		log.Println("获取了授权令牌，正在尝试对 XBL 进行身份验证...")
	} else {
		cachecontent, err := os.ReadFile(CacheFilename)
		if err != nil {
			return resauth, err
		}
		err = json.Unmarshal(cachecontent, &MSa)
		if err != nil {
			return resauth, err
		}
		MSaOld := MSa
		err = CheckRefreshMS(&MSa, cid)
		if err != nil {
			return resauth, err
		}
		if MSaOld.AccessToken != MSa.AccessToken {
			tocache, err := json.Marshal(MSa)
			if err != nil {
				return resauth, err
			}
			err = os.WriteFile(CacheFilename, tocache, 0600)
			if err != nil {
				return resauth, err
			}
		}
		log.Println("获取了缓存的授权令牌，正在尝试对 XBL 进行身份验证...")
	}

	XBLa, err := AuthXBL(MSa.AccessToken)
	if err != nil {
		return resauth, err
	}
	log.Println("在 XBL 上获得了授权，正在尝试获取 XSTS 令牌...")

	XSTSa, err := AuthXSTS(XBLa)
	if err != nil {
		return resauth, err
	}
	log.Println("获取了 XSTS 令牌，正在尝试获取 MC 令牌...")

	MCa, err := AuthMC(XSTSa)
	if err != nil {
		return resauth, err
	}
	log.Println("获取了 MC 令牌，不检查你是否拥有游戏，因为这太复杂了，直接去获取 MC 配置文件...")

	resauth, err = GetMCprofile(MCa.Token)
	if err != nil {
		return resauth, err
	}
	resauth.AsTk = MCa.Token
	return resauth, nil
}
