package whatslink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
)

// ApiResponse 定义API响应的结构体
type ApiResponse struct {
	Type        string `json:"type"`
	FileType    string `json:"file_type"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Count       int    `json:"count"`
	Screenshots []struct {
		Time      int    `json:"time"`
		Screenshot string `json:"screenshot"`
	} `json:"screenshots"`
}

// GetApiResponse 调用API并返回响应内容
func GetApiResponse(targetURL string) (*ApiResponse, error) {
    baseURL := "https://whatslink.info/api/v1/link"

    // 对目标 URL 进行编码
    escapedURL := url.QueryEscape(targetURL)

    // 创建请求URL
    reqURL := fmt.Sprintf("%s?url=%s", baseURL, escapedURL)
    fmt.Printf("URL:%s\n", reqURL)

    // 发送GET请求
    resp, err := http.Get(reqURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // 检查HTTP响应状态码
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
    }

    // 读取响应内容
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // 解析JSON响应
    var apiResponse ApiResponse
    err = json.Unmarshal(body, &apiResponse)
    if err != nil {
        return nil, err
    }

    return &apiResponse, nil
}

// DownloadImage 下载图像数据
func DownloadImage(imageURL string) ([]byte, error) {
    resp, err := http.Get(imageURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
    }

    return ioutil.ReadAll(resp.Body)
}

// BytesToMB 将字节数转换为MB
func BytesToMB(bytes int64) float64 {
    return float64(bytes) / 1024 / 1024
}

func init() {
	en := control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "WhatsLink API 查询",
		Help:             "-WhatsLink 查询\n链接查询 DDL/Torrent/Ed2k 链接",
	}).ApplySingle(ctxext.DefaultSingle)

	en.OnRegex(`^链接查询\s+((?:ed2k:\/\/|magnet:\?xt=urn:btih:[^\s&]+(?:&[^\s&]+)*|thunder:\/\/|ddl:\/\/[^\s]+))`).Limit(ctxext.LimitByGroup).SetBlock(true).
	Handle(func(ctx *zero.Ctx) {
		// 提取正则匹配的链接，并将 &amp; 转换为普通的 &
		targetURL := strings.ReplaceAll(ctx.State["regex_matched"].([]string)[1], "&amp;", "&")

		// 获取API响应
		response, err := GetApiResponse(targetURL)
		if err != nil {
			ctx.SendChain(
				message.At(ctx.Event.UserID),
				message.Text(fmt.Sprintf("Error fetching data: %v", err)),
			)
			return
		}

		// 检查返回数据是否为空
		if response.Type == "" && response.FileType == "" && response.Name == "" {
			ctx.SendChain(
				message.At(ctx.Event.UserID),
				message.Text("API返回的数据为空，请检查链接格式是否正确。"),
			)
			return
		}

		// 输出响应内容
		fmt.Printf("Type: %s\n", response.Type)
		fmt.Printf("FileType: %s\n", response.FileType)
		fmt.Printf("Name: %s\n", response.Name)
		fmt.Printf("Size: %d\n", response.Size)
		fmt.Printf("Count: %d\n", response.Count)

		// 发送响应内容
		ctx.SendChain(
			message.At(ctx.Event.UserID),
			message.Text(fmt.Sprintf("Type: %s\n", response.Type)),
			message.Text(fmt.Sprintf("FileType: %s\n", response.FileType)),
			message.Text(fmt.Sprintf("Name: %s\n", response.Name)),
			message.Text(fmt.Sprintf("Size: %.2f MB\n", BytesToMB(response.Size))),
			message.Text(fmt.Sprintf("Count: %d\n", response.Count)),
		)

		for _, screenshot := range response.Screenshots {
			fmt.Printf("Screenshot at %d: %s\n", screenshot.Time, screenshot.Screenshot)

			// 下载图片数据
			imgData, err := DownloadImage(screenshot.Screenshot)
			if err != nil {
				fmt.Println("图像下载失败:", err)
				continue
			}

			// 发送图片给用户
			ctx.SendChain(
				message.At(ctx.Event.UserID),
				message.ImageBytes(imgData),
			)
		}
	})
}
