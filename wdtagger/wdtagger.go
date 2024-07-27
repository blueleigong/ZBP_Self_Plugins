package wdtagger

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type TaggerInterrogateRequest struct {
	Image     string  `json:"image"`
	Model     string  `json:"model"`
	Threshold float64 `json:"threshold"`
}

type TaggerInterrogateResponse struct {
	Caption map[string]float64 `json:"caption"`
}

func readImageFromURL(url string) (string, error) {
	client := &http.Client{
		Timeout: 120 * time.Second, // 设置超时
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image, status code: %d", resp.StatusCode)
	}

	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageBytes), nil
}

func getTagsFromImageURL(imageURL, model string, threshold float64, apiURL string) (string, error) {
	imageBase64, err := readImageFromURL(imageURL)
	if err != nil {
		return "", err
	}

	requestPayload := TaggerInterrogateRequest{
		Image:     imageBase64,
		Model:     model,
		Threshold: threshold,
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时
	}
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var responsePayload TaggerInterrogateResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return "", err
	}

	tags := make([]string, 0, len(responsePayload.Caption))
	for tag := range responsePayload.Caption {
		tags = append(tags, tag)
	}

	return strings.Join(tags, ", "), nil
}

func filterTags(tags string) string {
	// 将下划线替换为空格
	tags = strings.ReplaceAll(tags, "_", " ")

	// 列出需要去除的词
	unwantedTags := []string{"sensitive", "general", "questionable", "explicit"}

	// 过滤标签
	tagList := strings.Split(tags, ", ")
	filteredTags := make([]string, 0, len(tagList))
	for _, tag := range tagList {
		if !contains(unwantedTags, tag) {
			filteredTags = append(filteredTags, tag)
		}
	}

	return strings.Join(filteredTags, ", ")
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func init() {
	model := "wd14-vit-v2-git"
	threshold := 0.35
	apiURL := "http://127.0.0.1:7860/tagger/v1/interrogate" //本地SDwebui的api地址

	engine := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "prompt反查",
		Help:              "- prompt反查[图片]",
		PrivateDataFolder: "prompt",
	})

	// 上传一张图进行评价
	engine.OnKeywordGroup([]string{"prompt反查"}, zero.OnlyGroup, zero.MustProvidePicture).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			for _, url := range ctx.State["image_url"].([]string) {
				tags, err := getTagsFromImageURL(url, model, threshold, apiURL)
				if err != nil {
					log.Printf("Error getting tags: %v", err)
					ctx.SendChain(
						message.At(ctx.Event.UserID),
						message.Text("获取标签失败，请稍后再试。"),
					)
					return
				}
				tags = filterTags(tags)
				fmt.Println("Tags:", tags)
				ctx.SendChain(
					message.At(ctx.Event.UserID),
					message.Text("\nprompt反查结果:\n", tags),
				)
			}
		})
}
